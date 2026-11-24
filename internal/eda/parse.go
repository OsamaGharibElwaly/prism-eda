package eda

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidCSV = errors.New("invalid CSV file")
	ErrMaxRows    = errors.New("dataset exceeds the maximum row count")
)

const maxRows = 200_000

type Table struct {
	Name        string
	Headers     []string
	Rows        [][]string
	SizeBytes   int64
	Delimiter   rune
	InvalidRows int
}

// Parse reads a delimited file, normalises ragged rows and de-duplicates headers.
func Parse(r io.Reader, name string, maxBytes int64) (*Table, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCSV, err)
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("%w: file is larger than %d MB", ErrInvalidCSV, maxBytes/(1024*1024))
	}

	text := strings.TrimPrefix(string(data), "\ufeff")
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("%w: file is empty", ErrInvalidCSV)
	}

	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = detectDelimiter(text)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCSV, err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("%w: expected a header and at least one data row", ErrInvalidCSV)
	}
	if len(records)-1 > maxRows {
		return nil, ErrMaxRows
	}

	headers := normaliseHeaders(records[0])
	width := len(headers)
	table := &Table{
		Name:      cleanName(name),
		Headers:   headers,
		Rows:      make([][]string, 0, len(records)-1),
		SizeBytes: int64(len(data)),
		Delimiter: reader.Comma,
	}

	for _, rec := range records[1:] {
		if len(rec) == 0 {
			continue
		}
		if len(rec) != width {
			table.InvalidRows++
		}
		row := make([]string, width)
		for i := 0; i < width; i++ {
			if i < len(rec) {
				row[i] = strings.TrimSpace(rec[i])
			}
		}
		table.Rows = append(table.Rows, row)
	}
	if len(table.Rows) == 0 {
		return nil, fmt.Errorf("%w: no data rows found", ErrInvalidCSV)
	}
	return table, nil
}

func detectDelimiter(text string) rune {
	line := strings.SplitN(text, "\n", 2)[0]
	best, bestCount := ',', 0
	for _, d := range []rune{',', ';', '\t', '|'} {
		if n := strings.Count(line, string(d)); n > bestCount {
			best, bestCount = d, n
		}
	}
	return best
}

func normaliseHeaders(raw []string) []string {
	seen := make(map[string]int, len(raw))
	out := make([]string, len(raw))
	for i, h := range raw {
		h = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(h), "\ufeff"))
		if h == "" {
			h = fmt.Sprintf("column_%d", i+1)
		}
		if n := seen[h]; n > 0 {
			out[i] = fmt.Sprintf("%s_%d", h, n+1)
		} else {
			out[i] = h
		}
		seen[h]++
	}
	return out
}

func cleanName(name string) string {
	name = strings.TrimSpace(filepathBase(name))
	if name == "" {
		return "dataset.csv"
	}
	return name
}

func filepathBase(name string) string {
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		return name[i+1:]
	}
	return name
}

// ParseNumber accepts plain decimals and values with thousands separators.
// Blank cells and placeholders return ok=false and are treated as missing.
func ParseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "NA" || s == "N/A" || s == "null" || s == "NULL" || s == "?" {
		return 0, false
	}
	if strings.Contains(s, ",") && strings.Contains(s, ".") {
		s = strings.ReplaceAll(s, ",", "")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v != v || v > 1e300 || v < -1e300 {
		return 0, false
	}
	return v, true
}
