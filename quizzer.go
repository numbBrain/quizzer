package main

import (
	
	"os"
	"fmt"
	"flag"
	"time"
	"encoding/csv"
)

type problem struct{

	question string
	answer string
}

func parseLines(csvData [][]string) []problem{

	var problems []problem

	for _, line := range csvData{

		p := problem{question : line[0], answer : line[1], }
		problems = append(problems, p)
	}

	return problems
}

func main(){

	csvFilename := flag.String("csv", "questions.csv", "CSV files to read questions from")
	timeout := flag.Int("timeout", 5, "time limit to answer a question")
	flag.Parse()

	timeoutDuration := time.Duration(*timeout) * time.Second
	csvFile, err := os.Open(*csvFilename)

	if err != nil{
		fmt.Printf("Error opening file: %s", *csvFilename)
		os.Exit(1)
	}

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()

	if err != nil{
		fmt.Printf("Error parsing file: %s", *csvFilename)
		os.Exit(1)
	}
	
	problems := parseLines(csvData)
	score := 0
	ticker := time.NewTicker(timeoutDuration)

	for i, p := range problems{

		answerChan := make(chan string)
		quit := false

		go func(){

			var a string
			fmt.Scanf("%s", &a)
			answerChan <- a

		}()

		fmt.Printf("\nProblem %d of %d: %s ", i+1, len(problems), p.question)

		select{

			case <- ticker.C:
				break

			case answer := <- answerChan:
				if answer == p.answer{
					score++		
				
				}else if answer == "quit"{
	
					quit = true
					ticker.Stop()	
					break		
				}

				ticker.Reset(timeoutDuration) 
		}

		if quit{

			fmt.Printf("\nYou quit the test! I'm surprised you're the sperm that won.\n")
			break
		}

	}

	fmt.Printf("\nYou scored %d out of %d\n", score, len(problems))
}
