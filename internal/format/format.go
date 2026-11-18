package format

// LineEnding defines supported line terminator types.
type LineEnding string

const (
	LF   LineEnding = "\n"
	CRLF LineEnding = "\r\n"
	CR   LineEnding = "\r"
)

// FileType defines supported file data formats.
type FileType string

const (
	TypeCSV      FileType = "CSV"
	TypeTSV      FileType = "TSV"
	TypeJSON     FileType = "JSON"
	TypeMarkdown FileType = "Markdown"
)
