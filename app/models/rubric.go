package models

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
)

// Rubric is all the feedback database
type Rubric struct {
	Design []FeedbackItem
	Style  []FeedbackItem
}

// FeedbackItem is a feedback for a section
type FeedbackItem struct {
	Name    string
	Desc    string
	Penalty float32
}

func (f FeedbackItem) ID() string {
	h := fnv.New32a()
	h.Write([]byte(f.Name))
	return fmt.Sprint(h.Sum32())
}

// Feedback all of it
var Feedback Rubric

// LoadDatabase loads from file into public var
func LoadDatabase(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	err = json.Unmarshal(data, &Feedback)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("Loaded Rubric")
}
