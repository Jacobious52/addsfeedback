package controllers

import (
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/Jacobious52/addsfeedback/app/models"
)

type statPack struct {
	Data  map[int]week
	Order []int
}

type week struct {
	Freq       models.Freq
	Order      []string
	Max        int
	Percentage func(w week, amount int) int
	Color      func(w week, amount int) string
}

func freqOrder(w week) []string {
	type pair struct {
		Key   string
		Value int
	}

	var order []pair
	for k, v := range w.Freq {
		order = append(order, pair{k, v})
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i].Value > order[j].Value
	})

	var keys []string
	for _, k := range order {
		keys = append(keys, k.Key)
	}
	return keys
}

func weekOrder(s statPack) []int {
	var keys []int
	for k := range s.Data {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	return keys
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

	var order []int
	pack := statPack{make(map[int]week), order}
	models.Stats.Lock.RLock()
	for w, f := range models.Stats.Data {
		max := 0
		for _, i := range f {
			if i > max {
				max = i
			}
		}
		var fOrder []string
		wk := week{f, fOrder, max, percentageFunc, colorFunc}
		wk.Order = freqOrder(wk)
		pack.Data[w] = wk
	}
	pack.Order = weekOrder(pack)
	models.Stats.Lock.RUnlock()

	err = tmpl.Execute(w, pack)
	if err != nil {
		http.Error(w, "Something bad happened. Sorry :(", 500)
		log.Println(err.Error())
		return
	}
}