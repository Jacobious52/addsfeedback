package models

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/go-github/github"
)

// Stat contains all the statisics for every week
type Stat struct {
	Data map[int]Freq
	Lock sync.RWMutex `json:"-"`
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

	Stats.Lock.Lock()

	if _, ok := Stats.Data[w]; !ok {
		Stats.Data[w] = make(Freq)
	}

	for _, r := range result {
		Stats.Data[w][r]++
	}

	Stats.Lock.Unlock()
}

// SaveStats Save the stats to disk
func SaveStats(ctx context.Context, client *github.Client, gist *github.Gist) {
	Stats.Lock.RLock()
	data, err := json.MarshalIndent(&Stats, "", "    ")
	Stats.Lock.RUnlock()

	if err != nil {
		log.Println(err.Error())
		return
	}
	filename := github.GistFilename("stats.json")
	file := gist.Files[filename]
	str := string(data)
	file.Content = &str
	gist.Files[filename] = file

	_, _, err = client.Gists.Edit(ctx, "2b896ec6e671a2b8b16c0e05198dcc83", gist)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Saved Stats")
}

// LoadStats loads the stats from disk
func LoadStats(ctx context.Context, client *github.Client) *github.Gist {
	gist, _, err := client.Gists.Get(ctx, "2b896ec6e671a2b8b16c0e05198dcc83")
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	file := gist.Files[github.GistFilename("stats.json")]
	data := file.GetContent()

	err = json.Unmarshal([]byte(data), &Stats)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Println("Loaded Stats")
	return gist
}