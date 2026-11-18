package format

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseCSV_Basic(t *testing.T) {
	input := "a,b,c\n1,2,3"
	got, err := ParseCSV(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestParseCSV_MultilineCells(t *testing.T) {
	// Cell containing internal newline
	input := "\"col1\",\"col\n2\",\"col3\"\n\"val\n1\",\"val2\",\"val3\""
	got, err := ParseCSV(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{
		{"col1", "col\n2", "col3"},
		{"val\n1", "val2", "val3"},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestParseTSV(t *testing.T) {
	input := "a\tb\tc\n1\t2\t3"
	got, err := ParseTSV(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestParseJSON(t *testing.T) {
	t.Run("Array of Arrays", func(t *testing.T) {
		input := `[["a", "b"], ["1", 2]]`
		got, err := ParseJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := [][]string{{"a", "b"}, {"1", "2"}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, want %v", got, expected)
		}
	})

	t.Run("Array of Objects", func(t *testing.T) {
		input := `[{"col1": "val1", "col2": "val2"}, {"col1": "val3", "col2": "val4"}]`
		got, err := ParseJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 || len(got[0]) != 2 {
			t.Fatalf("unexpected dims: %v", got)
		}
	})

	t.Run("Array of Scalars", func(t *testing.T) {
		input := `["a", "b", "c"]`
		got, err := ParseJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := [][]string{{"a"}, {"b"}, {"c"}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, want %v", got, expected)
		}
	})

	t.Run("Scalar", func(t *testing.T) {
		input := `"single"`
		got, err := ParseJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := [][]string{{"single"}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, want %v", got, expected)
		}
	})
}

func TestSerializeCSV(t *testing.T) {
	rows := [][]string{
		{"a", "b,c", "d\ne"},
		{"1", "2", "3"},
	}

	t.Run("LF", func(t *testing.T) {
		out, err := SerializeCSV(rows, LF)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "\"b,c\"") {
			t.Errorf("expected quoted comma, got: %q", out)
		}
		// verify roundtrip
		parsed, err := ParseCSV(out)
		if err != nil {
			t.Fatalf("parse error on roundtrip: %v", err)
		}
		if !reflect.DeepEqual(parsed, rows) {
			t.Errorf("roundtrip failed: got %v, want %v", parsed, rows)
		}
	})

	t.Run("CRLF", func(t *testing.T) {
		out, err := SerializeCSV(rows, CRLF)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "\r\n") {
			t.Errorf("expected CRLF line endings, got: %q", out)
		}
	})
}

func TestSerializeTSV(t *testing.T) {
	rows := [][]string{{"a", "b"}, {"c", "d"}}
	out, err := SerializeTSV(rows, LF)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "a\tb\nc\td"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestSerializeMarkdown(t *testing.T) {
	rows := [][]string{
		{"Header 1", "Header | 2"},
		{"Val 1", "Val\n2"},
	}
	out, err := SerializeMarkdown(rows, LF)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "| Header 1 | Header \\| 2 |\n| --- | --- |\n| Val 1 | Val 2 |"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestTestdataFixtures(t *testing.T) {
	dir := "../../testdata"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read testdata dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}
		content := string(contentBytes)

		t.Run(entry.Name(), func(t *testing.T) {
			if strings.HasSuffix(entry.Name(), ".csv") {
				records, err := ParseCSV(content)
				if err != nil {
					t.Fatalf("ParseCSV failed on fixture %s: %v", entry.Name(), err)
				}
				if len(records) == 0 {
					t.Fatalf("expected non-empty records for %s", entry.Name())
				}
			} else if strings.HasSuffix(entry.Name(), ".tsv") {
				records, err := ParseTSV(content)
				if err != nil {
					t.Fatalf("ParseTSV failed on fixture %s: %v", entry.Name(), err)
				}
				if len(records) == 0 {
					t.Fatalf("expected non-empty records for %s", entry.Name())
				}
			}
		})
	}
}

func FuzzCSVParser(f *testing.F) {
	// Seed corpus
	f.Add("a,b,c\n1,2,3")
	f.Add("\"hello\nworld\",test,123")
	f.Add("<html lang=\"jp\">,</html>,<script/>")
	f.Add("a\r\nb\r\nc")
	f.Add("a\rb\rc")

	f.Fuzz(func(t *testing.T, input string) {
		records, err := ParseCSV(input)
		if err != nil {
			// Parsing error on arbitrary input is acceptable
			return
		}
		// If parse succeeded, re-serializing should not crash
		out, err := SerializeCSV(records, LF)
		if err != nil {
			t.Fatalf("SerializeCSV failed after valid parse: %v", err)
		}
		// Roundtrip parse
		_, err = ParseCSV(out)
		if err != nil {
			t.Fatalf("ParseCSV failed on serialized output: %v", err)
		}
	})
}
