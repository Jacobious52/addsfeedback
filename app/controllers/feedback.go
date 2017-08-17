package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/Jacobious52/addsfeedback/app/models"
)

type FinalFeedback struct {
	Mark    float64
	Message string
}

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

	designMarks := 2.0
	styleMarks := 2.0
	functionalityMarks := 2.0

	// feedback buffer
	var allFeedbackBuff bytes.Buffer
	allFeedbackBuff.WriteString("Feedback:\n")

	for _, name := range models.OrderedKeys() {
		section, _ := models.Feedback[name]
		checked := false
		var feedbackBuff bytes.Buffer
		feedbackBuff.WriteString(fmt.Sprint(name, ":\n"))
		for _, v := range section {
			if r.Form.Get(v.ID()) == "on" {
				feedbackBuff.WriteString(v.Desc)
				feedbackBuff.WriteString(fmt.Sprint(" (", v.Penalty, ")."))
				feedbackBuff.WriteString("\n\n")

				checked = true

				// calculate for final marks
				if name == "Design" {
					designMarks += v.Penalty
				}
				if name == "Style" {
					styleMarks += v.Penalty
				}
				if name == "Functionality" {
					functionalityMarks += v.Penalty
				}
			}
		}
		// dont write if they got full marks for this section
		if checked {
			allFeedbackBuff.Write(feedbackBuff.Bytes())
		}
	}
	// cap the marks
	designMarks = math.Max(0, designMarks)
	styleMarks = math.Max(0, styleMarks)
	functionalityMarks = math.Max(0, functionalityMarks)

	// Create marks scheme
	var marksBuffer bytes.Buffer
	marksBuffer.WriteString("Design: ")
	marksBuffer.WriteString(fmt.Sprint(designMarks))
	marksBuffer.WriteString("/2\n")

	marksBuffer.WriteString("Style/Commenting: ")
	marksBuffer.WriteString(fmt.Sprint(styleMarks))
	marksBuffer.WriteString("/2\n")

	marksBuffer.WriteString("Functionality: ")
	marksBuffer.WriteString(fmt.Sprint(functionalityMarks))
	marksBuffer.WriteString("/2\n\n")

	// goodjob if full marks
	if (styleMarks + designMarks + functionalityMarks) == 6 {
		allFeedbackBuff.WriteString("Good work!\n")
	}

	// write all feedback
	marksBuffer.Write(allFeedbackBuff.Bytes())

	results := FinalFeedback{
		Mark:    styleMarks + designMarks + functionalityMarks,
		Message: marksBuffer.String(),
	}

	err = tmpl.Execute(w, results)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
