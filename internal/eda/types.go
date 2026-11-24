// Package eda turns a raw CSV table into a structured profiling report.
// It has no knowledge of HTTP, storage, or the frontend.
package eda

import "time"

type ColumnKind string

const (
	KindNumeric     ColumnKind = "numeric"
	KindCategorical ColumnKind = "categorical"
	KindBoolean     ColumnKind = "boolean"
)

type Bin struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Count int     `json:"count"`
}

type NumericStats struct {
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	Mean      float64 `json:"mean"`
	Median    float64 `json:"median"`
	StdDev    float64 `json:"std_dev"`
	Sum       float64 `json:"sum"`
	P25       float64 `json:"p25"`
	P75       float64 `json:"p75"`
	Outliers  int     `json:"outliers"`
	Histogram []Bin   `json:"histogram"`
}

type ValueCount struct {
	Value   string  `json:"value"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type CategoricalStats struct {
	TopValues []ValueCount `json:"top_values"`
	Entropy   float64      `json:"entropy"`
}

type ColumnProfile struct {
	Name        string            `json:"name"`
	Index       int               `json:"index"`
	Kind        ColumnKind        `json:"kind"`
	Count       int               `json:"count"`
	Missing     int               `json:"missing"`
	MissingPct  float64           `json:"missing_pct"`
	Distinct    int               `json:"distinct"`
	Numeric     *NumericStats     `json:"numeric,omitempty"`
	Categorical *CategoricalStats `json:"categorical,omitempty"`
}

type Correlation struct {
	X           string  `json:"x"`
	Y           string  `json:"y"`
	Coefficient float64 `json:"coefficient"`
	Strength    string  `json:"strength"`
	Direction   string  `json:"direction"`
}

type CorrelationMatrix struct {
	Labels []string    `json:"labels"`
	Values [][]float64 `json:"values"`
}

type Insight struct {
	Kind     string `json:"kind"`
	Severity string `json:"severity"` // warning | opportunity | info
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Column   string `json:"column,omitempty"`
}

type Quality struct {
	Score          float64 `json:"score"`
	Completeness   float64 `json:"completeness"`
	DuplicateRows  int     `json:"duplicate_rows"`
	InvalidRows    int     `json:"invalid_rows"`
	IssueCount     int     `json:"issue_count"`
}

type Profile struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Rows           int              `json:"rows"`
	ColumnCount    int              `json:"column_count"`
	CreatedAt      time.Time        `json:"created_at"`
	AnalysisMs     int64            `json:"analysis_ms"`
	NumericColumns []string         `json:"numeric_columns"`
	Columns        []ColumnProfile  `json:"columns"`
	Correlations   []Correlation    `json:"correlations"`
	Matrix         CorrelationMatrix `json:"matrix"`
	Insights       []Insight        `json:"insights"`
	Quality        Quality          `json:"quality"`
}

type Preview struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
	Total   int        `json:"total"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
