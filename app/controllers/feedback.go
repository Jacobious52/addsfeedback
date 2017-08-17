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

	// Mark design
	var designFeedbackBuff bytes.Buffer
	designMarks := 2.0
	designFeedbackBuff.WriteString("Design:\n")
	for _, v := range models.Feedback.Design {
		if r.Form.Get(v.ID()) == "on" {
			designFeedbackBuff.WriteString(v.Desc)
			designFeedbackBuff.WriteString(fmt.Sprint(" (", v.Penalty, ")."))
			designFeedbackBuff.WriteString("\n\n")
			designMarks += v.Penalty
		}
	}
	designMarks = math.Max(0, designMarks)

	// Mark style
	var styleFeedbackBuff bytes.Buffer
	styleMarks := 2.0
	styleFeedbackBuff.WriteString("\nStyle/Commenting:\n")
	for _, v := range models.Feedback.Style {
		if r.Form.Get(v.ID()) == "on" {
			styleFeedbackBuff.WriteString(v.Desc)
			styleFeedbackBuff.WriteString(fmt.Sprint(" (", v.Penalty, ")."))
			styleFeedbackBuff.WriteString("\n\n")
			styleMarks += v.Penalty
		}
	}
	styleMarks = math.Max(0, styleMarks)

	// Create marks
	var marksBuffer bytes.Buffer
	marksBuffer.WriteString("Design: ")
	marksBuffer.WriteString(fmt.Sprint(designMarks))
	marksBuffer.WriteString("/2\n")

	marksBuffer.WriteString("Style/Commenting: ")
	marksBuffer.WriteString(fmt.Sprint(styleMarks))
	marksBuffer.WriteString("/2\n")

	marksBuffer.WriteString("Functionality: 2/2")
	marksBuffer.WriteString("\n\n")

	// Write feedback
	var feedbackBuff bytes.Buffer
	// feedbackBuff.WriteString("Feedback:\n")

	// write design feedback
	if designMarks < 2 {
		feedbackBuff.Write(designFeedbackBuff.Bytes())
	}

	// write style feedback
	if styleMarks < 2 {
		feedbackBuff.Write(styleFeedbackBuff.Bytes())
	}

	// write full marks
	if (styleMarks + designMarks) == 4 {
		feedbackBuff.WriteString("Good work!\n")
	}

	// write all feedback
	marksBuffer.Write(feedbackBuff.Bytes())

	err = tmpl.Execute(w, marksBuffer.String())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
