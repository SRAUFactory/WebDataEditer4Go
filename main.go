package main

import (
	"html/template"
	"net/http"
)

type View struct {
	Title          string
	SelectFileType string
	SelectLfCode   string
	DataView       string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	v := &View{Title: "CSV/TSV形式編集ツール（Web版）", SelectFileType: "", SelectLfCode: "", DataView: ""}
	t := template.Must(template.ParseFiles("templete/index.html"))
	err := t.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
