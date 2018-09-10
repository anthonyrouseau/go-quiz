package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file in format 'question,answer'")
	timeLimit := flag.Int("limit", 30, "time limit in seconds")
	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Couldn't read the file")
	}
	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerChannel := make(chan string)
	defer close(answerChannel)
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		go getAnswer(answerChannel)
		select {
		case <-timer.C:
			fmt.Printf("\nOut of time! You scored %d out of %d", correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d", correct, len(problems))
}

func getAnswer(answerChannel chan string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	answerChannel <- answer
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{question: line[0], answer: strings.TrimSpace(line[1])}
	}
	return problems
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
