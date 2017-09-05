package main

import (
	"crypto/subtle"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Jacobious52/addsfeedback/app/controllers"
	"github.com/Jacobious52/addsfeedback/app/models"
)

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	rand.Seed(time.Now().Unix())

	models.LoadDatabase("db/feedback.json")
	models.LoadStats("db/stats.json")

	// routes
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	dbfs := http.FileServer(http.Dir("db"))
	http.Handle("/db/", http.StripPrefix("/db/", dbfs))

	http.HandleFunc("/", BasicAuth(controllers.Build, "addsmarker", "c++11", "addsmarkersite"))
	http.HandleFunc("/feedback", BasicAuth(controllers.Feedback, "addsmarker", "c++11", "addsmarkersite"))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
