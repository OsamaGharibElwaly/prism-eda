package eda

import (
	"math"
	"sort"
	"strings"
	"time"
)

// BuildProfile profiles every column, then derives correlations, a data-quality
// score and human-readable insights. The input table is never mutated.
func BuildProfile(id string, t *Table) *Profile {
	started := time.Now()
	n := len(t.Rows)
	totalMissing := 0
	cols := make([]ColumnProfile, len(t.Headers))

	idx := make([]int, 0, len(t.Headers))
	columns := make([][]float64, 0, len(t.Headers))

	for i, header := range t.Headers {
		values := make([]string, n)
		for r, row := range t.Rows {
			values[r] = row[i]
		}

		cp := ColumnProfile{Name: header, Index: i, Count: n, Kind: inferKind(values)}
		if cp.Kind == KindNumeric {
			parsed, missing := parseValues(values)
			cp.Missing = missing
			cp.Numeric = numericStats(parsed)
			cp.Distinct = distinctNumbers(parsed)
			if len(idx) < 12 {
				idx = append(idx, i)
				columns = append(columns, parsed)
			}
		} else {
			counts := valueCounts(values)
			cp.Missing = counts[""]
			cp.Distinct = len(counts) - boolToInt(cp.Missing > 0)
			cp.Categorical = categoricalStats(counts, n)
		}
		if n > 0 {
			cp.MissingPct = round1(100 * float64(cp.Missing) / float64(n))
		}
		totalMissing += cp.Missing
		cols[i] = cp
	}

	matrix, pairs := correlate(t.Headers, idx, columns)
	dups := duplicateRows(t)
	insights := buildInsights(cols, pairs, dups, n)
	quality := scoreQuality(t, cols, totalMissing, dups, n)
	quality.IssueCount = countIssues(insights)

	return &Profile{
		ID:             id,
		Name:           t.Name,
		Rows:           n,
		ColumnCount:    len(t.Headers),
		CreatedAt:      time.Now().UTC(),
		AnalysisMs:     time.Since(started).Milliseconds(),
		NumericColumns: numericNames(cols),
		Columns:        cols,
		Correlations:   pairs,
		Matrix:         matrix,
		Insights:       insights,
		Quality:        quality,
	}
}

func inferKind(values []string) ColumnKind {
	sample := nonEmpty(values, 2000)
	if len(sample) == 0 {
		return KindCategorical
	}
	numbers, booleans := 0, 0
	for _, v := range sample {
		if _, ok := ParseNumber(v); ok {
			numbers++
		}
		if isBool(v) {
			booleans++
		}
	}
	ratio := func(n int) float64 { return float64(n) / float64(len(sample)) }
	switch {
	case ratio(numbers) >= 0.9:
		return KindNumeric
	case ratio(booleans) >= 0.9:
		return KindBoolean
	default:
		return KindCategorical
	}
}

func nonEmpty(values []string, limit int) []string {
	out := make([]string, 0, min(len(values), limit))
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			continue
		}
		out = append(out, v)
		if len(out) == limit {
			break
		}
	}
	return out
}

func isBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "yes", "no", "y", "n", "true", "false", "t", "f":
		return true
	default:
		return false
	}
}

func parseValues(values []string) ([]float64, int) {
	out := make([]float64, len(values))
	missing := 0
	for i, v := range values {
		if num, ok := ParseNumber(v); ok {
			out[i] = num
		} else {
			out[i] = math.NaN()
			missing++
		}
	}
	return out, missing
}

func numericStats(values []float64) *NumericStats {
	sorted := make([]float64, 0, len(values))
	sum := 0.0
	for _, v := range values {
		if !math.IsNaN(v) {
			sorted = append(sorted, v)
			sum += v
		}
	}
	if len(sorted) == 0 {
		return &NumericStats{}
	}
	sort.Float64s(sorted)

	n := float64(len(sorted))
	mean := sum / n
	p25, p75 := percentile(sorted, 0.25), percentile(sorted, 0.75)
	iqr := p75 - p25
	variance := 0.0
	outliers := 0
	for _, v := range sorted {
		d := v - mean
		variance += d * d
		if iqr > 0 && (v < p25-1.5*iqr || v > p75+1.5*iqr) {
			outliers++
		}
	}
	stdDev := 0.0
	if len(sorted) > 1 {
		stdDev = math.Sqrt(variance / (n - 1))
	}

	return &NumericStats{
		Min:       sorted[0],
		Max:       sorted[len(sorted)-1],
		Mean:      mean,
		Median:    percentile(sorted, 0.5),
		StdDev:    stdDev,
		Sum:       sum,
		P25:       p25,
		P75:       p75,
		Outliers:  outliers,
		Histogram: histogram(sorted, p25, p75),
	}
}

func percentile(sorted []float64, p float64) float64 {
	pos := p * float64(len(sorted)-1)
	lo, hi := int(math.Floor(pos)), int(math.Ceil(pos))
	if lo == hi {
		return sorted[lo]
	}
	return sorted[lo] + (sorted[hi]-sorted[lo])*(pos-float64(lo))
}

// Freedman-Diaconis width with a Sturges fallback, clamped to a readable range.
func histogram(sorted []float64, p25, p75 float64) []Bin {
	min, max := sorted[0], sorted[len(sorted)-1]
	k := int(math.Ceil(math.Log2(float64(len(sorted))))) + 1
	if iqr := p75 - p25; iqr > 0 {
		if width := 2 * iqr / math.Cbrt(float64(len(sorted))); width > 0 {
			k = int(math.Ceil((max - min) / width))
		}
	}
	k = clampInt(k, 5, 20)
	width := (max - min) / float64(k)
	if width <= 0 {
		width = 1
	}
	bins := make([]Bin, k)
	for i := range bins {
		bins[i].Start = min + float64(i)*width
		bins[i].End = min + float64(i+1)*width
	}
	bins[k-1].End = max
	for _, v := range sorted {
		i := clampInt(int((v-min)/width), 0, k-1)
		bins[i].Count++
	}
	return bins
}

func distinctNumbers(values []float64) int {
	seen := make(map[float64]struct{})
	for _, v := range values {
		if !math.IsNaN(v) {
			seen[v] = struct{}{}
		}
	}
	return len(seen)
}

func valueCounts(values []string) map[string]int {
	counts := make(map[string]int)
	for _, v := range values {
		counts[strings.TrimSpace(v)]++
	}
	return counts
}

func categoricalStats(counts map[string]int, n int) *CategoricalStats {
	top := make([]ValueCount, 0, len(counts))
	entropy := 0.0
	for value, count := range counts {
		if value == "" {
			continue
		}
		p := float64(count) / float64(n)
		entropy -= p * math.Log2(p)
		top = append(top, ValueCount{Value: value, Count: count, Percent: round1(100 * p)})
	}
	sort.Slice(top, func(i, j int) bool {
		if top[i].Count != top[j].Count {
			return top[i].Count > top[j].Count
		}
		return top[i].Value < top[j].Value
	})
	if len(top) > 10 {
		top = top[:10]
	}
	return &CategoricalStats{TopValues: top, Entropy: entropy}
}

func correlate(headers []string, idx []int, columns [][]float64) (CorrelationMatrix, []Correlation) {
	labels := make([]string, len(idx))
	for i, columnIndex := range idx {
		labels[i] = headers[columnIndex]
	}
	matrix := CorrelationMatrix{Labels: labels, Values: make([][]float64, len(columns))}
	for i := range matrix.Values {
		matrix.Values[i] = make([]float64, len(columns))
		matrix.Values[i][i] = 1
	}

	pairs := make([]Correlation, 0)
	for i := 0; i < len(columns); i++ {
		for j := i + 1; j < len(columns); j++ {
			r, ok := pearson(columns[i], columns[j])
			if !ok {
				continue
			}
			matrix.Values[i][j], matrix.Values[j][i] = r, r
			pairs = append(pairs, Correlation{
				X:           labels[i],
				Y:           labels[j],
				Coefficient: r,
				Strength:    strength(r),
				Direction:   direction(r),
			})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return math.Abs(pairs[i].Coefficient) > math.Abs(pairs[j].Coefficient)
	})
	if len(pairs) > 16 {
		pairs = pairs[:16]
	}
	return matrix, pairs
}

func pearson(x, y []float64) (float64, bool) {
	n := 0
	var sx, sy float64
	for i := range x {
		if math.IsNaN(x[i]) || math.IsNaN(y[i]) {
			continue
		}
		sx += x[i]
		sy += y[i]
		n++
	}
	if n < 3 {
		return 0, false
	}
	mx, my := sx/float64(n), sy/float64(n)
	var cov, vx, vy float64
	for i := range x {
		if math.IsNaN(x[i]) || math.IsNaN(y[i]) {
			continue
		}
		dx, dy := x[i]-mx, y[i]-my
		cov += dx * dy
		vx += dx * dx
		vy += dy * dy
	}
	if vx == 0 || vy == 0 {
		return 0, false
	}
	return clamp(cov/math.Sqrt(vx*vy), -1, 1), true
}

func strength(r float64) string {
	switch a := math.Abs(r); {
	case a >= 0.7:
		return "strong"
	case a >= 0.4:
		return "moderate"
	default:
		return "weak"
	}
}

func direction(r float64) string {
	if r < 0 {
		return "negative"
	}
	return "positive"
}

func duplicateRows(t *Table) int {
	seen := make(map[string]struct{}, len(t.Rows))
	duplicates := 0
	for _, row := range t.Rows {
		key := strings.Join(row, "\x00")
		if _, ok := seen[key]; ok {
			duplicates++
			continue
		}
		seen[key] = struct{}{}
	}
	return duplicates
}

func scoreQuality(t *Table, cols []ColumnProfile, missing, duplicates, n int) Quality {
	if n == 0 {
		return Quality{}
	}
	completeness := 1 - float64(missing)/(float64(n)*float64(len(cols)))
	constants := 0
	for _, c := range cols {
		if c.Distinct <= 1 {
			constants++
		}
	}
	penalty := 20*float64(duplicates)/float64(n) +
		10*float64(t.InvalidRows)/float64(n) +
		10*float64(constants)/float64(len(cols))
	return Quality{
		Score:         clamp(round1(100*completeness-penalty), 0, 100),
		Completeness:  round1(100 * completeness),
		DuplicateRows: duplicates,
		InvalidRows:   t.InvalidRows,
	}
}

func numericNames(cols []ColumnProfile) []string {
	names := make([]string, 0)
	for _, c := range cols {
		if c.Kind == KindNumeric {
			names = append(names, c.Name)
		}
	}
	return names
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
