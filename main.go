package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// determine path
	path, good := os.LookupEnv("HW")
	if !good {
		path = "~/hw.txt"
		fmt.Printf("WARNING: $HW was not defined, defaulting to %s\n", path)
	}

	// attempt to open file
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("ERROR: couldn't open %s for reading", path))
	}

	// for each line
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	somethingDue := false
	for scanner.Scan() {
		// keep track of line number
		lineNumber++

		// scan
		tokens := strings.Split(scanner.Text(), "|")

		// don't freak out about blank lines
		if len(strings.TrimSpace(scanner.Text())) == 0 {
			continue
		}

		// validate
		if len(tokens) != 3 {
			fmt.Printf("WARNING: %d tokens on line %d (%s), expected 3\n", len(tokens), lineNumber, scanner.Text())
			continue
		}

		// clean
		for i, token := range tokens {
			tokens[i] = strings.TrimSpace(token)
		}

		// parse
		dueDate, err := time.Parse("01/02", tokens[1])

		if err != nil {
			fmt.Printf("WARNING: couldn't parse %s as a date\n", tokens[1])
			continue
		}

		lookAhead, err := strconv.ParseInt(tokens[2], 10, 64)
		if err != nil {
			fmt.Printf("WARNING: couldn't parse %s as an integer\n", tokens[2])
			continue
		}

		if lookAhead < 0 {
			fmt.Printf("WARNING: look-ahead is negative, notification will occur AFTER due date")
		}

		// calculate result and print
		daysUntil := dueDate.YearDay() - time.Now().YearDay() - 1 // this totally breaks if years don't match but semesters never span years
		if daysUntil > int(lookAhead) {
			continue
		}

		somethingDue = true
		if daysUntil > 0 {
			fmt.Printf("%s due in %d %s (%s)\n", tokens[0], daysUntil, days(daysUntil), dueDate.Weekday().String())
		} else if daysUntil == 0 {
			fmt.Printf("%s due TODAY\n", tokens[0])
		} else {
			fmt.Printf("%s is %d days LATE!!\n", tokens[0], daysUntil*-1)
		}

	}

	if !somethingDue {
		fmt.Println("nothing due!")
	}
}

func days(n int) string {
	if n == 1 {
		return "day"
	} else {
		return "days"
	}
}
