package handlers

import (
	"fmt"
	"net/http"
	"text/template"
)

func ServePredictor(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("template/predictor.html"))
	tmpl.Execute(w, nil)
	fmt.Println("Predictor page served.")
}
