package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type WCDHtmlSelect struct {
	Options  map[string]string
	Selected string
	Name     string
}

type WCDView struct {
	Title          string
	SelectFileType *WCDHtmlSelect
	SelectLfCode   *WCDHtmlSelect
	CsvData        [][]string
	FileName       string
}

/**
 * html Helper
 */
func getSelectTemplateHtml(sel *WCDHtmlSelect) template.HTML {
	html := ""
	for value, name := range sel.Options {
		var selected string
		if value == sel.Selected {
			selected = " selected"
		}
		html = fmt.Sprintf("%s<option value='%s'%s>%s</option>", html, value, selected, name)
	}
	return template.HTML(fmt.Sprintf("<select name='%s'>%s</select>", sel.Name, html))
}

func getAddButtonTemplateHtml(row int, col int) template.HTML {
	return template.HTML(fmt.Sprintf("<input type='submit' name='add[%d][%d]' value='+'>", row, col))
}

func getInputTextTemplateHtml(row int, col int, value string) template.HTML {
	return template.HTML(fmt.Sprintf("<input type='text' name='data[%d][]' value='%s'>", row /*col,*/, value))
}

/**
 * テンプレートの事前パース
 */
var (
	funcMap = template.FuncMap{
		"select":    getSelectTemplateHtml,
		"addButton": getAddButtonTemplateHtml,
		"inputText": getInputTextTemplateHtml,
	}
	tmpl = template.Must(template.New("base").Funcs(funcMap).ParseFiles("template/index.html"))
)

/**
 * http.Request Helper
 */
func getFormValue(r *http.Request, key string, dValue string) string {
	value := r.FormValue(key)
	if value == "" {
		value = dValue
	}
	return value
}

/**
 * request Handler
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fileType := getFormValue(r, "fileType", "csv")
	lfCode := getFormValue(r, "lfCode", "lf")
	fileName := getFormValue(r, "fileName", "")

	v := &WCDView{
		Title:          "CSV/TSV形式編集ツール（Web版）",
		SelectFileType: &WCDHtmlSelect{Options: map[string]string{"csv": "CSV", "tsv": "TSV"}, Selected: fileType, Name: "fileType"},
		SelectLfCode:   &WCDHtmlSelect{Options: map[string]string{"crlf": "CR+LF", "lf": "LF", "cr": "CR"}, Selected: lfCode, Name: "lfCode"},
		FileName:       fileName,
	}

	file, handler, _ := r.FormFile("uploadFile")
	if file != nil {
		defer file.Close()
		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				// サーバーを終了させないよう修正
				log.Printf("CSV Read Error: %v", err)
				http.Error(w, "CSVの解析に失敗しました: "+err.Error(), http.StatusBadRequest)
				return
			}

			v.CsvData = append(v.CsvData, record)
		}
		v.FileName = handler.Filename
	} else if r.Method == "POST" {
		for k, v := range r.Form {
			log.Printf("key : %s, value : %s", k, v)
		}
	}

	err := tmpl.ExecuteTemplate(w, "base", v)
	if err != nil {
		log.Printf("Template Render Error: %v", err)
		http.Error(w, "レンダリングエラーが発生しました", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CSV Editor on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
