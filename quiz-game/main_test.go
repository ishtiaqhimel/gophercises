package main

import (
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

func testReadCSVFile(t *testing.T) {
	str := "1+1,2\n3+4,7\n5+1,6\n"
	data, err := readCSVFile(strings.NewReader(str))
	assert.NilError(t, err)

	validData := [][]string{
		{"1+1", "2"},
		{"3+4", "7"},
		{"5+1", "6"},
	}
	assert.DeepEqual(t, data, validData)
}

func testGetTotalScore(t *testing.T) {
	timer := time.NewTimer(time.Duration(5) * time.Second)

	problems := []problem{
		{"3+3", "6"},
		{"2+4", "6"},
		{"7+9", "16"},
	}

	answerChan := make(chan string)
	answers := []string{"6", "6", "16"}
	go func() {
		for _, answer := range answers {
			answerChan <- answer
		}
	}()

	correct := 0
	for i, p := range problems {
		c := getScoreForSingleQuestion(i, p, timer.C, answerChan)
		if c == -1 {
			break
		}
		correct += c
	}

	assert.Equal(t, correct, len(answers))
}

func TestReadCSVFile(t *testing.T) {
	t.Run("test ReadCSVFile", testReadCSVFile)
}

func TestGetTotalScore(t *testing.T) {
	t.Run("test GetTotalScore", testGetTotalScore)
}
