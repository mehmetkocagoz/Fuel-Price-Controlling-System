package handlers

import (
	"fmt"
	"mehmetkocagz/model"
	"net/http"
	"text/template"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page served.")
	var tmpl = template.Must(template.ParseFiles("template/index.html"))
	dataAll := model.GrabTemplateData()
	tmpl.Execute(w, dataAll)
	fmt.Println("Home page served.")
}
