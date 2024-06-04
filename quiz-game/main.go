package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
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
		exit(fmt.Sprintf("failed to open %s file: %q\n", fileName, err.Error()))
	}

	data, err := readCSVFile(f)
	if err != nil {
		exit(err.Error())
	}

	totalQuestions := len(data)
	if totalQuestions == 0 {
		exit("no question available in the file!")
	}

	problems := parseData(data)

	correct := getTotalScore(problems, timeLimit)

	fmt.Printf("You score is %d out of %d.\n", correct, totalQuestions)
}

func readFlagValues() (string, int) {
	fileName := flag.String("filename", "problems.csv", "CSV File that conatins quiz questions")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz")
	flag.Parse()
	return *fileName, *timeLimit
}

func readCSVFile(f io.Reader) ([][]string, error) {
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %q", err)
	}
	return data, nil
}

func parseData(data [][]string) []problem {
	ret := make([]problem, len(data))
	for i, d := range data {
		ret[i] = problem{d[0], d[1]}
	}
	return ret
}

func getTotalScore(problems []problem, timeLimit int) int {
	correct := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	answerChan := make(chan string)
	go getAnswer(answerChan)

	for i, p := range problems {
		c := getScoreForSingleQuestion(i, p, timer.C, answerChan)
		if c == -1 {
			fmt.Println()
			return correct
		}
		correct += c
	}
	return correct
}

func getScoreForSingleQuestion(i int, p problem, timerC <-chan time.Time, answerChan <-chan string) int {
	fmt.Printf("Question #%d: %s = ", i+1, p.question)

	select {
	case <-timerC:
		return -1
	case answer := <-answerChan:
		if strings.Compare(p.answer, strings.ToLower(strings.TrimSpace(answer))) == 0 {
			return 1
		}
		return 0
	}
}

func getAnswer(answerChan chan string) {
	for {
		var ans string
		fmt.Scanf("%s\n", &ans)
		answerChan <- ans
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
