package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Jacobious52/addsfeedback/app/models"
)

type statPack struct {
	Data map[int]week
}

type week struct {
	Freq       models.Freq
	Max        int
	Percentage func(w week, amount int) int
	Color      func(w week, amount int) string
}

func percentageFunc(w week, amount int) int {
	return int(float64(amount) / float64(w.Max) * 100.0)
}

func colorFunc(w week, amount int) string {
	p := percentageFunc(w, amount)
	if p <= 25 {
		return "bg-success"
	}
	if p <= 50 {
		return "bg-warning"
	}
	return "bg-danger"
}

// Statistics displays statistics neatly
func Statistics(w http.ResponseWriter, r *http.Request) {

	log.Println("/stats", r.Method)

	tmpl, err := template.ParseFiles("app/views/stats.html")
	if err != nil {
		http.Error(w, "Something bad happened. Sorry :(", 500)
		log.Println(err.Error())
		return
	}

	pack := statPack{make(map[int]week)}
	models.Stats.Lock.RLock()
	for w, f := range models.Stats.Data {
		max := 0
		for _, i := range f {
			if i > max {
				max = i
			}
		}
		wk := week{f, max, percentageFunc, colorFunc}
		pack.Data[w] = wk
	}
	models.Stats.Lock.RUnlock()

	err = tmpl.Execute(w, pack)
	if err != nil {
		http.Error(w, "Something bad happened. Sorry :(", 500)
		log.Println(err.Error())
		return
	}
}