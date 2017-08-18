package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Jacobious52/addsfeedback/app/models"
)

type OrderedFeedback struct {
	Order    []string
	Feedback models.Rubric
}

func Build(w http.ResponseWriter, r *http.Request) {

	log.Println("/build", r.Method)

	tmpl, err := template.ParseFiles("app/views/build.html")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	pack := OrderedFeedback{Order: models.OrderedKeys(), Feedback: models.Feedback}

	err = tmpl.Execute(w, pack)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
