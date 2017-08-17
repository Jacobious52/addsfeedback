package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jacobious52/addsfeedback/app/controllers"
	"github.com/Jacobious52/addsfeedback/app/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	models.LoadDatabase("feedback.json")

	// routes
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/feedback", controllers.Feedback)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
