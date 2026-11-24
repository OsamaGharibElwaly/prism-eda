package eda

import (
	"fmt"
	"math"
)

// Preview returns the first limit rows for the data table widget.
func (t *Table) Preview(limit int) *Preview {
	if limit <= 0 || limit > len(t.Rows) {
		limit = len(t.Rows)
	}
	rows := make([][]string, 0, limit)
	rows = append(rows, t.Rows[:limit]...)
	return &Preview{Headers: t.Headers, Rows: rows, Total: len(t.Rows)}
}

// Scatter returns paired numeric values for two columns, skipping blanks.
func (t *Table) Scatter(x, y string, limit int) ([]Point, error) {
	xi, xok := t.columnIndex(x)
	yi, yok := t.columnIndex(y)
	if !xok {
		return nil, fmt.Errorf("unknown column: %s", x)
	}
	if !yok {
		return nil, fmt.Errorf("unknown column: %s", y)
	}
	if xi == yi {
		return nil, fmt.Errorf("x and y must be different columns")
	}
	if limit <= 0 || limit > len(t.Rows) {
		limit = len(t.Rows)
	}
	points := make([]Point, 0, limit)
	for _, row := range t.Rows {
		xv, okX := ParseNumber(row[xi])
		yv, okY := ParseNumber(row[yi])
		if okX && okY && !math.IsInf(xv, 0) && !math.IsInf(yv, 0) {
			points = append(points, Point{X: xv, Y: yv})
		}
		if len(points) == limit {
			break
		}
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("no overlapping numeric values for %s and %s", x, y)
	}
	return points, nil
}

func (t *Table) columnIndex(name string) (int, bool) {
	for i, h := range t.Headers {
		if h == name {
			return i, true
		}
	}
	return 0, false
}
