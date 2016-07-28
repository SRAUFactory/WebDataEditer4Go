package main

import (
	"fmt"
	"html/template"
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
	DataView       string
}

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"select": getSelectTemplateHtml,
	}
	v := &WCDView{
		Title:          "CSV/TSV形式編集ツール（Web版）",
		SelectFileType: &WCDHtmlSelect{Options: map[string]string{"csv": "CSV", "tsv": "TSV"}, Selected: "csv", Name: "fileType"},
		SelectLfCode:   &WCDHtmlSelect{Options: map[string]string{"crlf": "CR+LF", "lf": "LF", "cr": "CR"}, Selected: "lf", Name: "lfCode"},
		DataView:       "",
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
