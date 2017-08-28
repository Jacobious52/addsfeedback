package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Jacobious52/addsfeedback/app/models"
)

type FinalFeedback struct {
	Mark    float64
	Message string
}

func Feedback(w http.ResponseWriter, r *http.Request) {
	log.Println("/feedback", r.Method)

	if r.Method != "POST" {
		io.WriteString(w, "only accepts POST request. got "+r.Method)
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

	for _, name := range models.OrderedKeys() {
		section, _ := models.Feedback[name]
		checked := false
		var feedbackBuff bytes.Buffer
		feedbackBuff.WriteString(fmt.Sprint("\n# ", name, ":\n"))
		for _, v := range section {
			if r.Form.Get(v.ID()) == "on" {

				v.Penalty = -math.Abs(v.Penalty)

				feedbackBuff.WriteString(v.Desc)
				if v.Penalty == 0 {
					feedbackBuff.WriteString(".")
				} else {
					feedbackBuff.WriteString(fmt.Sprint(" (", v.Penalty, ")."))
				}
				feedbackBuff.WriteString("\n")

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

		extraText := r.Form.Get(fmt.Sprint("extra-comments-", name))
		if extraText != "" {
			extraPenalty := r.Form.Get(fmt.Sprint("extra-penalty-", name))

			penalty, err := strconv.ParseFloat(extraPenalty, 64)
			if err != nil {
				log.Println("Bad extra penalty", err.Error())
				penalty = 0
			}
			penalty = -math.Abs(penalty)

			feedbackBuff.WriteString(extraText)
			if penalty == 0 {
				feedbackBuff.WriteString(".")
			} else {
				feedbackBuff.WriteString(fmt.Sprint(" (", penalty, ")."))
			}
			feedbackBuff.WriteString("\n")
			checked = true

			// calculate for final marks
			if name == "Design" {
				designMarks += penalty
			}
			if name == "Style" {
				styleMarks += penalty
			}
			if name == "Functionality" {
				functionalityMarks += penalty
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
	marksBuffer.WriteString("/2\n")

	// goodjob if full marks
	if (styleMarks + designMarks + functionalityMarks) == 6 {
		goodFeedback := []string{
			"Good work!",
			"Good job",
			"Nice job!",
			"Nice work",
			"6/6",
			"Good code.",
			"Top stuff",
			"Nice",
		}

		goodBuffer := fmt.Sprint("\n\n", goodFeedback[rand.Intn(len(goodFeedback))])
		allFeedbackBuff.WriteString(goodBuffer)
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
