package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Answers struct {
	Name                                     string
	Answered, Unanswered, Correct, Incorrect int
}

var readch chan string

//Default Total number of Quiz
var TotalQues int = 10

func (an *Answers) Read(filename string) {
	readch = make(chan string, 10)
	filecur, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file", err)
	}

	recor := csv.NewReader(filecur)
	defer filecur.Close()
	for {
		record, err := recor.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error in reading file", err)
		}

		readch <- record[0]
		readch <- record[1]
	}
}

var ans string

func (an *Answers) AskQuestion() {

	if len(readch) <= 0 {
		an.endquiz()
	}
	ques := <-readch
	anscsv := <-readch

	fmt.Println(ques)
	fmt.Println("Enter the Answer")
	fmt.Scanln(&ans)
	if strings.TrimSpace(anscsv) == ans {
		an.Answered += 1
		an.Correct += 1
	} else {
		an.Answered += 1
		an.Incorrect += 1
	}

}

func (an *Answers) endquiz() {
	fmt.Println("Quiz ended")
	an.Unanswered = TotalQues - (an.Answered)
	fmt.Printf(" Your Final Score is %+v", an)
	os.Exit(0)
}

func main() {
	fmt.Println("Quiz Program")

	fmt.Println("Enter your name")
	var name string
	fmt.Scanln(&name)
	var an Answers
	an.Name = name
	go an.Read("src/main/ques.csv")
	//Default timer as 60 seconds
	timer := time.NewTimer(60 * time.Second)
	defer timer.Stop()
	go func() {
		<-timer.C
		fmt.Println("Time Up....!!! ")
		an.endquiz()
	}()
	for {

		var a string

		fmt.Println("Press ENTER for next question, else press anyother key to EXIT")
		fmt.Scanln(&a)
		if a == "" {
			if len(readch) <= 0 {
				an.endquiz()
			}

			//fmt.Println("Next Question")
		} else {
			an.endquiz()
		}

		an.AskQuestion()
	}
}
