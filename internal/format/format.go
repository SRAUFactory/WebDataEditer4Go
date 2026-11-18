package format

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// LineEnding defines supported line terminator types.
type LineEnding string

const (
	LF   LineEnding = "LF"
	CRLF LineEnding = "CRLF"
	CR   LineEnding = "CR"
)

// Separator returns the string separator for the LineEnding.
func (l LineEnding) Separator() string {
	switch l {
	case CRLF:
		return "\r\n"
	case CR:
		return "\r"
	default:
		return "\n"
	}
}

// FileType defines supported file data formats.
type FileType string

const (
	TypeCSV      FileType = "CSV"
	TypeTSV      FileType = "TSV"
	TypeJSON     FileType = "JSON"
	TypeMarkdown FileType = "Markdown"
)

// normalizeLineEndings converts all standalone \r or \r\n to \n for standard Go parsing.
func normalizeLineEndings(input string) string {
	// First convert \r\n to \n
	s := strings.ReplaceAll(input, "\r\n", "\n")
	// Then convert remaining standalone \r to \n
	return strings.ReplaceAll(s, "\r", "\n")
}

// ParseData parses input content based on format.
func ParseData(content string, fileType FileType) ([][]string, error) {
	switch fileType {
	case TypeCSV:
		return ParseCSV(content)
	case TypeTSV:
		return ParseTSV(content)
	case TypeJSON:
		return ParseJSON(content)
	default:
		return nil, fmt.Errorf("unsupported parse format: %s", fileType)
	}
}

// ParseCSV parses CSV text using standard encoding/csv.
// It normalizes line endings so that CR/CRLF/LF are all handled properly.
func ParseCSV(content string) ([][]string, error) {
	if content == "" {
		return [][]string{}, nil
	}
	normalized := normalizeLineEndings(content)
	r := csv.NewReader(strings.NewReader(normalized))
	r.FieldsPerRecord = -1 // Allow variable number of fields per record
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV parse error: %w", err)
	}
	return records, nil
}

// ParseTSV parses TSV text.
func ParseTSV(content string) ([][]string, error) {
	if content == "" {
		return [][]string{}, nil
	}
	normalized := normalizeLineEndings(content)
	r := csv.NewReader(strings.NewReader(normalized))
	r.Comma = '\t'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("TSV parse error: %w", err)
	}
	return records, nil
}

// ParseJSON parses JSON input.
// Supports:
// 1. Array of arrays: [["a","b"], ["c","d"]]
// 2. Array of objects: [{"k1":"v1", "k2":"v2"}, ...]
// 3. Array of scalars: ["a", "b"] -> [["a"], ["b"]]
// 4. Scalar value: "a" -> [["a"]]
func ParseJSON(content string) ([][]string, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return [][]string{}, nil
	}

	var raw interface{}
	dec := json.NewDecoder(strings.NewReader(trimmed))
	// Prevent float numbers from converting unexpectedly if numbers are present
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return parseJSONValue(raw)
}

func stringifyVal(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case json.Number:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parseJSONValue(raw interface{}) ([][]string, error) {
	switch val := raw.(type) {
	case []interface{}:
		if len(val) == 0 {
			return [][]string{}, nil
		}
		// Check first element to determine array type
		first := val[0]
		switch first.(type) {
		case []interface{}:
			// Array of arrays
			result := make([][]string, len(val))
			for i, rowRaw := range val {
				if rowSlice, ok := rowRaw.([]interface{}); ok {
					row := make([]string, len(rowSlice))
					for j, cell := range rowSlice {
						row[j] = stringifyVal(cell)
					}
					result[i] = row
				} else {
					result[i] = []string{stringifyVal(rowRaw)}
				}
			}
			return result, nil

		case map[string]interface{}:
			// Array of objects
			// Collect keys in order of appearance
			keys := make([]string, 0)
			keySet := make(map[string]bool)
			for _, itemRaw := range val {
				if obj, ok := itemRaw.(map[string]interface{}); ok {
					for k := range obj {
						if !keySet[k] {
							keySet[k] = true
							keys = append(keys, k)
						}
					}
				}
			}
			result := make([][]string, len(val))
			for i, itemRaw := range val {
				row := make([]string, len(keys))
				if obj, ok := itemRaw.(map[string]interface{}); ok {
					for j, k := range keys {
						if cellVal, exists := obj[k]; exists {
							row[j] = stringifyVal(cellVal)
						}
					}
				}
				result[i] = row
			}
			return result, nil

		default:
			// Array of scalars
			result := make([][]string, len(val))
			for i, item := range val {
				result[i] = []string{stringifyVal(item)}
			}
			return result, nil
		}

	default:
		// Scalar value
		return [][]string{{stringifyVal(val)}}, nil
	}
}

// SerializeData serializes rows to text format.
func SerializeData(rows [][]string, fileType FileType, lineEnding LineEnding) (string, error) {
	switch fileType {
	case TypeCSV:
		return SerializeCSV(rows, lineEnding)
	case TypeTSV:
		return SerializeTSV(rows, lineEnding)
	case TypeJSON:
		return SerializeJSON(rows)
	case TypeMarkdown:
		return SerializeMarkdown(rows, lineEnding)
	default:
		return "", fmt.Errorf("unsupported serialize format: %s", fileType)
	}
}

// SerializeCSV serializes rows to CSV using encoding/csv and custom line terminator.
func SerializeCSV(rows [][]string, lineEnding LineEnding) (string, error) {
	if len(rows) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	// encoding/csv writes \n as line endings. If CRLF or CR is requested, convert.
	s := buf.String()
	// Strip the trailing \n if present for precise joining
	hasTrailingNewline := strings.HasSuffix(s, "\n")
	trimmed := strings.TrimSuffix(s, "\n")
	if lineEnding == CRLF {
		// Replace standard \n with \r\n
		trimmed = strings.ReplaceAll(trimmed, "\n", "\r\n")
		if hasTrailingNewline {
			trimmed += "\r\n"
		}
	} else if lineEnding == CR {
		trimmed = strings.ReplaceAll(trimmed, "\n", "\r")
		if hasTrailingNewline {
			trimmed += "\r"
		}
	}
	return trimmed, nil
}

// SerializeTSV serializes rows to TSV using custom line terminator.
func SerializeTSV(rows [][]string, lineEnding LineEnding) (string, error) {
	if len(rows) == 0 {
		return "", nil
	}
	sep := lineEnding.Separator()
	lines := make([]string, len(rows))
	for i, row := range rows {
		lines[i] = strings.Join(row, "\t")
	}
	return strings.Join(lines, sep), nil
}

// SerializeJSON serializes rows to JSON string (2D array format).
func SerializeJSON(rows [][]string) (string, error) {
	if rows == nil {
		return "[]", nil
	}
	b, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// escapeMarkdownCell escapes pipes and replaces newlines with space.
func escapeMarkdownCell(cell string) string {
	s := strings.ReplaceAll(cell, "|", `\|`)
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "\r", " ")
}

// SerializeMarkdown serializes rows into a GitHub Flavored Markdown table.
// Row 0 is the header, followed by a delimiter row, followed by data rows.
func SerializeMarkdown(rows [][]string, lineEnding LineEnding) (string, error) {
	if len(rows) == 0 {
		return "", nil
	}
	sep := lineEnding.Separator()
	colCount := len(rows[0])
	if colCount == 0 {
		return "", nil
	}

	formatRow := func(row []string) string {
		escaped := make([]string, colCount)
		for j := 0; j < colCount; j++ {
			if j < len(row) {
				escaped[j] = escapeMarkdownCell(row[j])
			} else {
				escaped[j] = ""
			}
		}
		return "| " + strings.Join(escaped, " | ") + " |"
	}

	headerRow := formatRow(rows[0])
	delims := make([]string, colCount)
	for j := range delims {
		delims[j] = "---"
	}
	delimRow := "| " + strings.Join(delims, " | ") + " |"

	lines := []string{headerRow, delimRow}
	for i := 1; i < len(rows); i++ {
		lines = append(lines, formatRow(rows[i]))
	}

	return strings.Join(lines, sep), nil
}

// Helper to check reader errors
var _ = io.EOF
