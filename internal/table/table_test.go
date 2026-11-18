package table

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tbl := New(3, 4)
	if len(tbl.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(tbl.Rows))
	}
	for i, r := range tbl.Rows {
		if len(r) != 4 {
			t.Fatalf("row %d: expected 4 cols, got %d", i, len(r))
		}
		for j, c := range r {
			if c != "" {
				t.Fatalf("cell (%d, %d) expected empty, got %q", i, j, c)
			}
		}
	}
}

func TestInsertRow(t *testing.T) {
	tests := []struct {
		name      string
		initial   [][]string
		targetIdx int
		pos       Position
		expected  [][]string
	}{
		{
			name:      "insert after row 0",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 0,
			pos:       PositionAfter,
			expected:  [][]string{{"a", "b"}, {"", ""}, {"c", "d"}},
		},
		{
			name:      "insert before row 0",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 0,
			pos:       PositionBefore,
			expected:  [][]string{{"", ""}, {"a", "b"}, {"c", "d"}},
		},
		{
			name:      "insert after row 1",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 1,
			pos:       PositionAfter,
			expected:  [][]string{{"a", "b"}, {"c", "d"}, {"", ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := FromRows(tt.initial)
			tbl.InsertRow(tt.targetIdx, tt.pos)
			if !reflect.DeepEqual(tbl.Rows, tt.expected) {
				t.Errorf("got %v, want %v", tbl.Rows, tt.expected)
			}
		})
	}
}

func TestDeleteRow(t *testing.T) {
	t.Run("delete middle row", func(t *testing.T) {
		tbl := FromRows([][]string{{"1"}, {"2"}, {"3"}})
		tbl.DeleteRow(1)
		expected := [][]string{{"1"}, {"3"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})

	t.Run("cannot delete when only 1 row remains", func(t *testing.T) {
		tbl := FromRows([][]string{{"only"}})
		tbl.DeleteRow(0)
		expected := [][]string{{"only"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})
}

func TestInsertCol(t *testing.T) {
	tests := []struct {
		name      string
		initial   [][]string
		targetIdx int
		pos       Position
		expected  [][]string
	}{
		{
			name:      "insert after col 0",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 0,
			pos:       PositionAfter,
			expected:  [][]string{{"a", "", "b"}, {"c", "", "d"}},
		},
		{
			name:      "insert before col 0",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 0,
			pos:       PositionBefore,
			expected:  [][]string{{"", "a", "b"}, {"", "c", "d"}},
		},
		{
			name:      "insert after col 1",
			initial:   [][]string{{"a", "b"}, {"c", "d"}},
			targetIdx: 1,
			pos:       PositionAfter,
			expected:  [][]string{{"a", "b", ""}, {"c", "d", ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := FromRows(tt.initial)
			tbl.InsertCol(tt.targetIdx, tt.pos)
			if !reflect.DeepEqual(tbl.Rows, tt.expected) {
				t.Errorf("got %v, want %v", tbl.Rows, tt.expected)
			}
		})
	}
}

func TestDeleteCol(t *testing.T) {
	t.Run("delete middle col", func(t *testing.T) {
		tbl := FromRows([][]string{{"1", "2", "3"}, {"4", "5", "6"}})
		tbl.DeleteCol(1)
		expected := [][]string{{"1", "3"}, {"4", "6"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})

	t.Run("cannot delete when only 1 col remains", func(t *testing.T) {
		tbl := FromRows([][]string{{"1"}, {"2"}})
		tbl.DeleteCol(0)
		expected := [][]string{{"1"}, {"2"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})
}

func TestAssignNumbers(t *testing.T) {
	t.Run("assign row", func(t *testing.T) {
		tbl := FromRows([][]string{{"a", "b", "c"}, {"d", "e", "f"}})
		tbl.AssignNumbers("row", 0)
		expected := [][]string{{"1", "2", "3"}, {"d", "e", "f"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})

	t.Run("assign col", func(t *testing.T) {
		tbl := FromRows([][]string{{"a", "b", "c"}, {"d", "e", "f"}})
		tbl.AssignNumbers("col", 1)
		expected := [][]string{{"a", "1", "c"}, {"d", "2", "f"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})
}

func TestCopy(t *testing.T) {
	t.Run("copy row", func(t *testing.T) {
		tbl := FromRows([][]string{{"a", "b"}, {"c", "d"}})
		tbl.Copy("row", 0)
		expected := [][]string{{"a", "b"}, {"a", "b"}, {"c", "d"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})

	t.Run("copy col", func(t *testing.T) {
		tbl := FromRows([][]string{{"a", "b"}, {"c", "d"}})
		tbl.Copy("col", 0)
		expected := [][]string{{"a", "a", "b"}, {"c", "c", "d"}}
		if !reflect.DeepEqual(tbl.Rows, expected) {
			t.Errorf("got %v, want %v", tbl.Rows, expected)
		}
	})
}

func TestApplyOperation(t *testing.T) {
	tbl := FromRows([][]string{{"a", "b"}, {"c", "d"}})
	err := tbl.ApplyOperation("copyRow", 0, PositionAfter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"a", "b"}, {"a", "b"}, {"c", "d"}}
	if !reflect.DeepEqual(tbl.Rows, expected) {
		t.Errorf("got %v, want %v", tbl.Rows, expected)
	}

	err = tbl.ApplyOperation("unknownAction", 0, PositionAfter)
	if err == nil {
		t.Fatalf("expected error for unknown action, got nil")
	}
}
