package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Jacobious52/addsfeedback/app/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("app/views/index.html")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	err = tmpl.Execute(w, models.Feedback)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
