package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
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
	return template.HTML(fmt.Sprintf("<input type='text' name='data[%d][%d]' value='%s'>", row, col, value))
}

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
	funcMap := template.FuncMap{
		"select":    getSelectTemplateHtml,
		"addButton": getAddButtonTemplateHtml,
		"inputText": getInputTextTemplateHtml,
	}

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
				log.Fatal("Error : ", err)
			}

			v.CsvData = append(v.CsvData, record)
		}
		v.FileName = handler.Filename
	}

	t := template.Must(template.New(v.Title).Funcs(funcMap).ParseFiles("template/index.html"))
	err := t.ExecuteTemplate(w, "base", v)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
