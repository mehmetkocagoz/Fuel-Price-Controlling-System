package handlers

import (
	"fmt"
	"mehmetkocagz/datafunctions"
	"mehmetkocagz/model"
	"net/http"
	"text/template"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("template/index.html"))
	dataAll := model.GrabTemplateData()
	tmpl.Execute(w, dataAll)
	fmt.Println("Home page served.")
}

func ServeHomeWithDate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Date requested")
	selectedDate := r.FormValue("selected-date")
	selectedDateTimestamp := datafunctions.ConvertTimestampFormatYMD(selectedDate)
	dataAllWDate := model.GrabTemplateDataWDate(selectedDateTimestamp)
	var tmpl = template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, dataAllWDate)
	fmt.Println("Home page served.")
}
