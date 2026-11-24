package eda

import (
	"math"
	"strings"
	"testing"
)

const fixture = "price,area,room,flag\n10,100,1,yes\n20,200,2,no\n,300,2,yes\n40,,1,no\n60,500,3,yes\n"

func TestParseAndProfile(t *testing.T) {
	table, err := Parse(strings.NewReader(fixture), "sample.csv", 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if table.Rows[0][0] != "10" || len(table.Headers) != 4 || len(table.Rows) != 5 {
		t.Fatalf("unexpected table: %+v", table)
	}

	p := BuildProfile("abc", table)
	if p.Rows != 5 || p.ColumnCount != 4 || p.ID != "abc" {
		t.Fatalf("unexpected profile meta: %+v", p)
	}
	if p.Columns[0].Kind != KindNumeric || p.Columns[3].Kind != KindBoolean {
		t.Fatalf("bad type inference: %+v", p.Columns)
	}
	if p.Columns[0].Missing != 1 || p.Columns[0].MissingPct != 20 {
		t.Fatalf("bad missing count: %+v", p.Columns[0])
	}
	if got := p.Columns[0].Numeric.Mean; math.Abs(got-32.5) > 0.001 {
		t.Fatalf("mean = %v, want 32.5", got)
	}
	if got := p.Columns[1].Numeric.Median; got != 250 {
		t.Fatalf("median = %v, want 250", got)
	}
	if len(p.Matrix.Labels) != 3 || p.Matrix.Values[0][1] < 0.99 {
		t.Fatalf("bad correlation matrix: %+v", p.Matrix)
	}
	if len(p.Insights) == 0 || p.Quality.Score < 0 || p.Quality.Score > 100 {
		t.Fatalf("bad insights/quality: %+v %+v", p.Insights, p.Quality)
	}
	if p.Quality.Completeness != 90 {
		t.Fatalf("completeness = %v, want 90", p.Quality.Completeness)
	}
}

func TestParseRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"", "name\n", strings.Repeat("x", 100)} {
		if _, err := Parse(strings.NewReader(input), "bad.csv", 10); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}

func TestScatterSkipsMissingPairs(t *testing.T) {
	table, err := Parse(strings.NewReader(fixture), "sample.csv", 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	points, err := table.Scatter("price", "area", 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3", len(points))
	}
	if points[0].X != 10 || points[2].Y != 500 {
		t.Fatalf("unexpected points: %+v", points)
	}
	if _, err := table.Scatter("nope", "area", 100); err == nil {
		t.Fatal("expected unknown column error")
	}
}

func TestRaggedRowsAndDuplicateHeaders(t *testing.T) {
	table, err := Parse(strings.NewReader("a,a,b\n1,2\n3,4,5,6\n"), "ragged.csv", 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if table.Headers[1] != "a_2" || table.InvalidRows != 2 {
		t.Fatalf("headers=%v invalid=%d", table.Headers, table.InvalidRows)
	}
	if len(table.Rows[0]) != 3 || table.Rows[0][2] != "" {
		t.Fatalf("row not padded: %v", table.Rows[0])
	}
}
