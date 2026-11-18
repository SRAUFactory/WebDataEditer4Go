//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/SRAUFactory/WebDataEditer4Go/internal/format"
	"github.com/SRAUFactory/WebDataEditer4Go/internal/table"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Content string      `json:"content,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func jsonResponse(success bool, data interface{}, content string, errMsg string) string {
	res := Response{
		Success: success,
		Data:    data,
		Content: content,
		Error:   errMsg,
	}
	b, err := json.Marshal(res)
	if err != nil {
		return fmt.Sprintf(`{"success":false,"error":%q}`, err.Error())
	}
	return string(b)
}

// parseData(content: string, format: string, lineEnding: string) -> json string
func parseDataWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return jsonResponse(false, nil, "", "parseData requires at least 2 arguments: content, format")
	}
	content := args[0].String()
	fileType := format.FileType(args[1].String())

	rows, err := format.ParseData(content, fileType)
	if err != nil {
		return jsonResponse(false, nil, "", err.Error())
	}
	return jsonResponse(true, rows, "", "")
}

// serializeData(rowsJson: string, format: string, lineEnding: string) -> json string
func serializeDataWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return jsonResponse(false, nil, "", "serializeData requires 3 arguments: rowsJson, format, lineEnding")
	}
	rowsJson := args[0].String()
	fileType := format.FileType(args[1].String())
	lineEnding := format.LineEnding(args[2].String())

	var rows [][]string
	if err := json.Unmarshal([]byte(rowsJson), &rows); err != nil {
		return jsonResponse(false, nil, "", fmt.Sprintf("invalid rowsJson: %v", err))
	}

	content, err := format.SerializeData(rows, fileType, lineEnding)
	if err != nil {
		return jsonResponse(false, nil, "", err.Error())
	}
	return jsonResponse(true, nil, content, "")
}

type opParams struct {
	Index    int            `json:"index"`
	Position table.Position `json:"position"`
}

// applyTableOperation(rowsJson: string, action: string, paramsJson: string) -> json string
func applyTableOperationWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return jsonResponse(false, nil, "", "applyTableOperation requires 3 arguments: rowsJson, action, paramsJson")
	}
	rowsJson := args[0].String()
	action := args[1].String()
	paramsJson := args[2].String()

	var rows [][]string
	if err := json.Unmarshal([]byte(rowsJson), &rows); err != nil {
		return jsonResponse(false, nil, "", fmt.Sprintf("invalid rowsJson: %v", err))
	}

	var params opParams
	if paramsJson != "" {
		if err := json.Unmarshal([]byte(paramsJson), &params); err != nil {
			return jsonResponse(false, nil, "", fmt.Sprintf("invalid paramsJson: %v", err))
		}
	}

	tbl := table.FromRows(rows)
	if err := tbl.ApplyOperation(action, params.Index, params.Position); err != nil {
		return jsonResponse(false, nil, "", err.Error())
	}

	return jsonResponse(true, tbl.Rows, "", "")
}

func main() {
	c := make(chan struct{}, 0)

	api := js.Global().Get("Object").New()
	api.Set("parseData", js.FuncOf(parseDataWrapper))
	api.Set("serializeData", js.FuncOf(serializeDataWrapper))
	api.Set("applyTableOperation", js.FuncOf(applyTableOperationWrapper))

	js.Global().Set("WebDataEditor", api)

	// Dispatch ready event if window exists
	window := js.Global().Get("window")
	if !window.IsUndefined() && !window.IsNull() {
		eventConstructor := window.Get("CustomEvent")
		if !eventConstructor.IsUndefined() {
			evt := eventConstructor.New("WebDataEditorReady")
			window.Call("dispatchEvent", evt)
		}
	}

	<-c
}
