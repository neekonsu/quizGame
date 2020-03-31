package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Question represents a single question
type Question struct {
	content string
	answer  int
}

// Quiz represents all information about a quiz; all questions, accuracy, and itme elapsed (seconds)
type Quiz struct {
	filepath                    string
	questions                   []Question
	correct, incorrect, timeCap int
	timeElapsed                 float32
}

// adds single question from row of csv
func (q *Quiz) addQuestion(row []string) {
	var output Question
	output.content = row[0]
	output.answer, _ = strconv.Atoi(row[1])

	q.questions = append(q.questions, output)
}

// simple error checker, reports line as well
func checkErr(line string, e error) {
	if e != nil {
		log.Fatal("Error on line: "+line, e)
	}
}

var quiz Quiz

// read csv and populate global quiz object
func init() {
	flag.StringVar(&quiz.filepath, "path", "./quiz.csv", "expects string filepath")
	flag.IntVar(&quiz.timeCap, "timeCap", 5, "sets cap for quiz time")
	flag.Parse()
	file, err := os.Open(quiz.filepath)
	checkErr("28", err)
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		checkErr("32", err)
		quiz.addQuestion(row)
	}
}

// run the quiz itself
func main() {
	var wg sync.WaitGroup
	timer := time.NewTimer(time.Duration(quiz.timeCap) * time.Second)
	wg.Add(1)
	go func(wgi *sync.WaitGroup) {
		for _, question := range quiz.questions {
			var input string
			fmt.Printf("What is %s: ", question.content)
			fmt.Scanln(&input)
			if response, _ := strconv.Atoi(input); response == question.answer {
				quiz.correct++
			} else {
				quiz.incorrect++
			}
		}
		wgi.Done()
	}(&wg)
	<-timer.C
	wg.Done()
	if (quiz.correct + quiz.incorrect) == 0 {
		fmt.Println("\nTime's up! You got 0%")
	} else {
		fmt.Println(fmt.Sprintf("\nTime's up! You got %.2f%%\n", 100*(float32(quiz.correct)/float32(quiz.correct+quiz.incorrect))))
	}
}
