package main

import (
	"context"
	"crypto/subtle"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jacobious52/addsfeedback/app/controllers"
	"github.com/Jacobious52/addsfeedback/app/models"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("No Gist Token specified. Add env var $TOKEN")
	}

	rand.Seed(time.Now().Unix())

	// setup gist api
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	gist := models.LoadStats(ctx, client)

	models.LoadDatabase("db/feedback.json")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			models.SaveStats(ctx, client, gist)
			os.Exit(0)
		}
	}()

	// routes
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	dbfs := http.FileServer(http.Dir("db"))
	http.Handle("/db/", http.StripPrefix("/db/", dbfs))

	http.HandleFunc("/", BasicAuth(controllers.Build, "addsmarker", "c++11", "addsmarkersite"))
	http.HandleFunc("/feedback", BasicAuth(controllers.Feedback, "addsmarker", "c++11", "addsmarkersite"))
	http.HandleFunc("/stats", BasicAuth(controllers.Statistics, "addsmarker", "c++11", "addsmarkersite"))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err.Error())
	}
	models.SaveStats(ctx, client, gist)
}