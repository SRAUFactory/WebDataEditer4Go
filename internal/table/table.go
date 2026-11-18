package table

import (
	"fmt"
	"strconv"
)

// Table represents a two-dimensional grid of string cells.
type Table struct {
	Rows [][]string
}

// New creates a new Table with given rows and columns initialized to empty strings.
func New(rowCount, colCount int) *Table {
	if rowCount < 0 {
		rowCount = 0
	}
	if colCount < 0 {
		colCount = 0
	}
	rows := make([][]string, rowCount)
	for i := 0; i < rowCount; i++ {
		rows[i] = make([]string, colCount)
	}
	return &Table{Rows: rows}
}

// FromRows creates a new Table from an existing 2D slice, copying all cells.
func FromRows(rows [][]string) *Table {
	if rows == nil {
		return &Table{Rows: [][]string{}}
	}
	newRows := make([][]string, len(rows))
	for i, r := range rows {
		newRows[i] = make([]string, len(r))
		copy(newRows[i], r)
	}
	return &Table{Rows: newRows}
}

// Clone returns a deep copy of Table.
func (t *Table) Clone() *Table {
	return FromRows(t.Rows)
}

// Position specifies whether an insertion occurs before or after the index.
type Position string

const (
	PositionBefore Position = "before"
	PositionAfter  Position = "after"
)

// InsertRow inserts a new empty row before or after targetIdx.
// If table is empty, inserts a row with 1 empty column.
func (t *Table) InsertRow(targetIdx int, pos Position) {
	colCount := 1
	if len(t.Rows) > 0 {
		colCount = len(t.Rows[0])
	}
	newRow := make([]string, colCount)

	if len(t.Rows) == 0 {
		t.Rows = [][]string{newRow}
		return
	}

	insertIdx := targetIdx + 1
	if pos == PositionBefore {
		insertIdx = targetIdx
	}
	if insertIdx < 0 {
		insertIdx = 0
	}
	if insertIdx > len(t.Rows) {
		insertIdx = len(t.Rows)
	}

	t.Rows = append(t.Rows[:insertIdx], append([][]string{newRow}, t.Rows[insertIdx:]...)...)
}

// DeleteRow deletes the row at targetIdx.
// If the table has 1 or fewer rows, it does nothing (at least 1 row must remain).
func (t *Table) DeleteRow(targetIdx int) {
	if len(t.Rows) <= 1 {
		return
	}
	if targetIdx < 0 || targetIdx >= len(t.Rows) {
		return
	}
	t.Rows = append(t.Rows[:targetIdx], t.Rows[targetIdx+1:]...)
}

// InsertCol inserts a new empty column before or after targetIdx.
// If table is empty, inserts a 1x1 table.
func (t *Table) InsertCol(targetIdx int, pos Position) {
	if len(t.Rows) == 0 {
		t.Rows = [][]string{{""}}
		return
	}

	insertIdx := targetIdx + 1
	if pos == PositionBefore {
		insertIdx = targetIdx
	}
	if insertIdx < 0 {
		insertIdx = 0
	}
	numCols := len(t.Rows[0])
	if insertIdx > numCols {
		insertIdx = numCols
	}

	for i := range t.Rows {
		row := t.Rows[i]
		newRow := make([]string, 0, len(row)+1)
		newRow = append(newRow, row[:insertIdx]...)
		newRow = append(newRow, "")
		newRow = append(newRow, row[insertIdx:]...)
		t.Rows[i] = newRow
	}
}

// DeleteCol deletes the column at targetIdx.
// If columns are 1 or fewer, it does nothing (at least 1 column must remain).
func (t *Table) DeleteCol(targetIdx int) {
	if len(t.Rows) == 0 || len(t.Rows[0]) <= 1 {
		return
	}
	if targetIdx < 0 || targetIdx >= len(t.Rows[0]) {
		return
	}

	for i := range t.Rows {
		row := t.Rows[i]
		t.Rows[i] = append(row[:targetIdx], row[targetIdx+1:]...)
	}
}

// AssignNumbers fills sequential numbers starting from 1 across a row or column.
// mode can be "row" or "col".
func (t *Table) AssignNumbers(mode string, targetIdx int) {
	if len(t.Rows) == 0 {
		return
	}
	if mode == "row" {
		if targetIdx < 0 || targetIdx >= len(t.Rows) {
			return
		}
		for j := range t.Rows[targetIdx] {
			t.Rows[targetIdx][j] = strconv.Itoa(j + 1)
		}
	} else if mode == "col" {
		if targetIdx < 0 || targetIdx >= len(t.Rows[0]) {
			return
		}
		for i := range t.Rows {
			t.Rows[i][targetIdx] = strconv.Itoa(i + 1)
		}
	}
}

// Copy duplicates a row or column and inserts it right after targetIdx.
// mode can be "row" or "col".
func (t *Table) Copy(mode string, targetIdx int) {
	if len(t.Rows) == 0 {
		return
	}
	if mode == "row" {
		if targetIdx < 0 || targetIdx >= len(t.Rows) {
			return
		}
		copiedRow := make([]string, len(t.Rows[targetIdx]))
		copy(copiedRow, t.Rows[targetIdx])

		insertIdx := targetIdx + 1
		t.Rows = append(t.Rows[:insertIdx], append([][]string{copiedRow}, t.Rows[insertIdx:]...)...)
	} else if mode == "col" {
		if targetIdx < 0 || targetIdx >= len(t.Rows[0]) {
			return
		}
		insertIdx := targetIdx + 1
		for i := range t.Rows {
			val := t.Rows[i][targetIdx]
			row := t.Rows[i]
			newRow := make([]string, 0, len(row)+1)
			newRow = append(newRow, row[:insertIdx]...)
			newRow = append(newRow, val)
			newRow = append(newRow, row[insertIdx:]...)
			t.Rows[i] = newRow
		}
	}
}

// ApplyOperation dispatches table operations based on action name.
func (t *Table) ApplyOperation(action string, index int, pos Position) error {
	switch action {
	case "insertRow":
		t.InsertRow(index, pos)
	case "deleteRow":
		t.DeleteRow(index)
	case "insertCol":
		t.InsertCol(index, pos)
	case "deleteCol":
		t.DeleteCol(index)
	case "assignRow":
		t.AssignNumbers("row", index)
	case "assignCol":
		t.AssignNumbers("col", index)
	case "copyRow":
		t.Copy("row", index)
	case "copyCol":
		t.Copy("col", index)
	default:
		return fmt.Errorf("unknown table operation: %s", action)
	}
	return nil
}
