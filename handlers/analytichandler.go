package handlers

import (
	"fmt"
	"net/http"
	"text/template"
)

func ServeAnalysis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Analysis page served.")
	var tmpl = template.Must(template.ParseFiles("template/analytic.html"))
	tmpl.Execute(w, nil)
}
