package eda

import (
	"fmt"
	"math"
	"sort"
)

func buildInsights(cols []ColumnProfile, corrs []Correlation, duplicates, n int) []Insight {
	out := make([]Insight, 0)
	for _, c := range cols {
		if c.MissingPct >= 5 {
			out = append(out, Insight{
				Kind: "missing", Severity: "warning", Column: c.Name,
				Title:  fmt.Sprintf("%s is %v%% empty", c.Name, c.MissingPct),
				Detail: fmt.Sprintf("%d of %d rows have no value. Impute or drop this column before modelling.", c.Missing, n),
			})
		}
		if c.Distinct <= 1 {
			out = append(out, Insight{
				Kind: "constant", Severity: "warning", Column: c.Name,
				Title:  fmt.Sprintf("%s is constant", c.Name),
				Detail: "A single distinct value carries no signal; consider removing it.",
			})
		}
		if n > 20 && c.Kind == KindNumeric && c.Numeric != nil {
			s := c.Numeric
			if s.Outliers > 0 && float64(s.Outliers)/float64(n) > 0.02 {
				out = append(out, Insight{
					Kind: "outlier", Severity: "info", Column: c.Name,
					Title:  fmt.Sprintf("%s has %d outliers", c.Name, s.Outliers),
					Detail: fmt.Sprintf("Values outside 1.5 × IQR (%.1f – %.1f).", s.P25-1.5*(s.P75-s.P25), s.P75+1.5*(s.P75-s.P25)),
				})
			}
			if s.StdDev > 0 && math.Abs(s.Mean-s.Median) > 0.5*s.StdDev {
				out = append(out, Insight{
					Kind: "skew", Severity: "info", Column: c.Name,
					Title:  fmt.Sprintf("%s is skewed", c.Name),
					Detail: fmt.Sprintf("Mean %.2f vs median %.2f — prefer the median for reporting.", s.Mean, s.Median),
				})
			}
		}
		if c.Kind != KindNumeric && n > 20 && c.Distinct > 20 && float64(c.Distinct)/float64(n) > 0.5 {
			out = append(out, Insight{
				Kind: "cardinality", Severity: "info", Column: c.Name,
				Title:  fmt.Sprintf("%s has high cardinality", c.Name),
				Detail: fmt.Sprintf("%d distinct values across %d rows; group rare levels before encoding.", c.Distinct, n),
			})
		}
	}
	if duplicates > 0 {
		out = append(out, Insight{
			Kind: "duplicate", Severity: "warning",
			Title:  fmt.Sprintf("%d duplicate rows", duplicates),
			Detail: "Identical rows can bias statistics and inflate model confidence.",
		})
	}
	shown := 0
	for _, c := range corrs {
		if math.Abs(c.Coefficient) < 0.7 || shown == 4 {
			continue
		}
		out = append(out, Insight{
			Kind: "correlation", Severity: "opportunity",
			Title:  fmt.Sprintf("%s ↔ %s (r = %.2f)", c.X, c.Y, c.Coefficient),
			Detail: fmt.Sprintf("Strong %s relationship — a useful candidate for modelling or a scatter plot.", c.Direction),
		})
		shown++
	}

	severity := map[string]int{"warning": 0, "opportunity": 1, "info": 2}
	sort.SliceStable(out, func(i, j int) bool {
		return severity[out[i].Severity] < severity[out[j].Severity]
	})
	if len(out) > 12 {
		out = out[:12]
	}
	return out
}

func countIssues(insights []Insight) int {
	n := 0
	for _, i := range insights {
		if i.Severity == "warning" {
			n++
		}
	}
	return n
}
