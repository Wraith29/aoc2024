package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func logErr(err error) {
	fmt.Printf("ERROR: %+v\n", err)

	os.Exit(1)
}

func useSample(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if args[0] == "sample" {
		return true
	}

	return false
}

func getInput(args []string) (string, error) {
	fileName := "input.txt"
	if useSample(args) {
		fileName = "sample.txt"
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(data), "\n"), nil
}

func main() {
	input, err := getInput(os.Args[1:])

	if err != nil {
		logErr(err)
	}

	p1, err := part1(input)

	if err != nil {
		logErr(err)
	}

	p2, err := part2(input)

	if err != nil {
		logErr(err)
	}

	fmt.Printf("Part 1: %d\nPart 2: %d\n", p1, p2)
}

func dist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

type report struct {
	levels []int
	isDesc bool
}

func createReport(line string) (*report, error) {
	levels := make([]int, 0)

	for _, level := range strings.Split(line, " ") {
		if len(level) == 0 {
			continue
		}

		val, err := strconv.Atoi(strings.Trim(level, " "))
		if err != nil {
			return nil, err
		}

		levels = append(levels, val)
	}

	left := 0
	right := 0
	idx := 0
	isDesc := true

	for left == right {
		left = levels[idx]
		right = levels[idx+1]

		idx += 1
	}

	if left < right {
		isDesc = false
	}

	return &report{
		levels,
		isDesc,
	}, nil
}

func (r *report) isSafe() bool {
	prevLevel := r.levels[0]

	for _, level := range r.levels[1:] {
		distance := dist(level, prevLevel)
		if distance > 3 || distance < 1 {
			return false
		}

		if (r.isDesc && level > prevLevel) || (!r.isDesc && level < prevLevel) {
			return false
		}

		prevLevel = level
	}

	return true
}

func part1(input string) (int, error) {
	safe := 0

	for _, line := range strings.Split(input, "\n") {
		rep, err := createReport(line)

		if err != nil {
			return -1, err
		}

		if rep.isSafe() {
			safe += 1
		}
	}

	return safe, nil
}

func part2(input string) (int, error) {
	safe := 0

	// for _, line := range strings.Split(input, "\n") {
	// rep, err := createReport(line)

	// if err != nil {
	// return -1, nil
	// }
	// }

	return safe, nil
}
