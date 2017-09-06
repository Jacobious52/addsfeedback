package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// Stat contains all the statisics for every week
type Stat struct {
	Data map[int]Freq
	lock sync.RWMutex
}

// Freq is the frequency a name of a penatly occurs
type Freq map[string]int

// Result is all the results from a student build feedback
type Result []string

// NewResult creates a new result for a request
func NewResult() Result {
	return make(Result, 0)
}

// RecordResult appends to the results array
func (r Result) RecordResult(name string) Result {
	return append(r, name)
}

// SaveResult Saves the result concurrently
func (r Result) SaveResult() {
	go Stats.Add(r)
}

// Stats is the global stats database object
var Stats *Stat

// Add a result to the feedback
func (s *Stat) Add(result Result) {
	_, w := time.Now().ISOWeek()

	Stats.lock.Lock()

	if _, ok := Stats.Data[w]; !ok {
		Stats.Data[w] = make(Freq)
	}

	for _, r := range result {
		Stats.Data[w][r]++
	}

	Stats.lock.Unlock()

	SaveStats("db/stats.json")
}

// SaveStats Save the stats to disk
func SaveStats(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	Stats.lock.RLock()
	data, err := json.Marshal(&Stats)
	Stats.lock.RUnlock()

	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = file.Write(data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Saved Stats")
}

// LoadStats loads the stats from disk
func LoadStats(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(data, &Stats)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Loaded Stats")
}
