package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type WCDHtmlSelect struct {
	Options  map[string]string
	Selected string
}

type WCDView struct {
	Title          string
	SelectFileType *WCDHtmlSelect
	SelectLfCode   *WCDHtmlSelect
	DataView       string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"select": func(sel *WCDHtmlSelect) template.HTML {
			html := ""
			for value, name := range sel.Options {
				var selected string
				if value == sel.Selected {
					selected = " selected"
				}

				html = fmt.Sprintf("%s<option value='%s'%s>%s</option>", html, value, selected, name)
			}
			return template.HTML(fmt.Sprintf("<select>%s</select>", html))
		},
	}

	selectFileType := new(WCDHtmlSelect)
	selectFileType.Options = map[string]string{"csv": "CSV", "tsv": "TSV"}
	selectFileType.Selected = "csv"

	selectLfCode := new(WCDHtmlSelect)
	selectLfCode.Options = map[string]string{"crlf": "CR+LF", "lf": "LF", "cr": "CR"}
	selectLfCode.Selected = "lf"

	v := &WCDView{Title: "CSV/TSV形式編集ツール（Web版）", SelectFileType: selectFileType, SelectLfCode: selectLfCode, DataView: ""}
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
