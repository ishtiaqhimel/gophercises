package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	fileName, timeLimit := readFlagValues()
	f, err := os.Open(fileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open %s file: %q\n", fileName, err.Error()))
	}

	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read %s file: %q\n", fileName, err.Error()))
	}

	totalQuestions := len(data)
	if totalQuestions == 0 {
		exit("No question available in the file!")
	}

	problems := parseData(data)
	correct := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
OutterLoop:
	for i, p := range problems {
		fmt.Printf("Question #%d: %s = ", i+1, p.question)
		answerChan := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s\n", &ans)
			answerChan <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break OutterLoop // We can also return from here
		case answer := <-answerChan:
			if strings.Compare(p.answer, strings.ToLower(strings.TrimSpace(answer))) == 0 {
				correct++
			}
		}
	}
	fmt.Printf("You score is %d out of %d.\n", correct, totalQuestions)
}

func readFlagValues() (string, int) {
	fileName := flag.String("filename", "problems.csv", "CSV File that conatins quiz questions")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz")
	flag.Parse()
	return *fileName, *timeLimit
}

func parseData(data [][]string) []problem {
	ret := make([]problem, len(data))
	for i, d := range data {
		ret[i] = problem{d[0], d[1]}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
