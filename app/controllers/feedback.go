package controllers

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/Jacobious52/addsfeedback/app/models"
)

func Feedback(w http.ResponseWriter, r *http.Request) {
	log.Println("/feedback", r.Method)

	if r.Method != "POST" {
		io.WriteString(w, "bad request")
		return
	}

	r.ParseForm()

	tmpl, err := template.ParseFiles("app/views/feedback.html")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var buff bytes.Buffer

	buff.WriteString("Design:\n")
	for _, v := range models.Feedback.Design {
		if r.Form.Get(v.ID()) == "on" {
			buff.WriteString(v.Desc)
			buff.WriteString("\n\n")
		}
	}

	err = tmpl.Execute(w, buff.String())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
