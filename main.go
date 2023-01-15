package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func problemPuller(fileName string) ([]problem, error) {
	// read all the problems from the file

	// 1. open the file
	if fObj, err := os.Open(fileName); err == nil {
		// 2. create new reader
		csvR := csv.NewReader(fObj)
		// 3. read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			// 4. call the parseProblem func
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error reading data in csv"+"format from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}
}

func parseProblem(lines [][]string) []problem {
	// go over the line and parse with problem struct
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{
			q: lines[i][0],
			a: lines[i][1],
		}
	}
	return r
}

func exit(msg string) {
	fmt.Print(msg)
	os.Exit(1)
}

func main() {
	// 1.input the name of the file
	// fName := flag("f", "quiz.csv", "path of csv file")
	fName := flag.String("f", "quiz.csv", "path of csv file")
	// 2.set the duration of timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	// 3.pull the problems from the file (call the puller func)
	problems, err := problemPuller(*fName)
	// 4.handle the error
	if err != nil {
		exit(fmt.Sprintf("something wrong : %s", err.Error()))
	}
	// 5. create a variable to count the correct answer
	correctAns := 0
	// 6. use the duration of timer, initialize the time
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)
	// 7. loop through the problem, print question, accept answer

problemLoop:

	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}
	// 8. calculate and print result
	fmt.Printf("Your result is %d out of %d \n", correctAns, len(problems))
	fmt.Printf("Press enter to exit\n")
	<-ansC
}
