package table

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
