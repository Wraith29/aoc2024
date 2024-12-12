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

	fmt.Printf("Part 1: %d\n", p1)

	p2, err := part2(input)

	if err != nil {
		logErr(err)
	}

	fmt.Printf("Part 2: %d\n", p2)
}

func dist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

type report struct {
	idx    int
	levels []int
	isDesc bool
}

func createReport(idx int, line string) (*report, error) {
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

	isDesc := levels[0] > levels[len(levels)-1]

	return &report{
		idx,
		levels,
		isDesc,
	}, nil
}

func checkLevels(levels []int, isDesc bool) (bool, int) {
	prevLevel := levels[0]

	for idx, level := range levels[1:] {
		distance := dist(level, prevLevel)
		if distance > 3 || distance < 1 {
			return false, idx + 1
		}

		if (isDesc && level > prevLevel) || (!isDesc && level < prevLevel) {
			return false, idx + 1
		}

		prevLevel = level
	}

	return true, -1
}

func copyAndRemove(s []int, i int) []int {
	new := make([]int, len(s))

	copy(new, s)

	return append(new[:i], new[i+1:]...)
}

func (r report) isSafeWithTolerance() bool {
	safe, idx := checkLevels(r.levels, r.isDesc)
	if safe {
		return safe
	}

	rightSide := copyAndRemove(r.levels, idx)
	if s, _ := checkLevels(rightSide, r.isDesc); s {
		return true
	}

	leftSide := copyAndRemove(r.levels, idx-1)
	if s, _ := checkLevels(leftSide, r.isDesc); s {
		return true
	}

	return false
}

// Answer: 486
func part1(input string) (int, error) {
	safe := 0

	for idx, line := range strings.Split(input, "\n") {
		rep, err := createReport(idx, line)

		if err != nil {
			return -1, err
		}

		if isSafe, _ := checkLevels(rep.levels, rep.isDesc); isSafe {
			safe += 1
		}
	}

	return safe, nil
}

// Answer: 540
func part2(input string) (int, error) {
	safe := 0

	for idx, line := range strings.Split(input, "\n") {
		rep, err := createReport(idx, line)

		if err != nil {
			return -1, nil
		}

		if rep.isSafeWithTolerance() {
			safe += 1
		}
	}

	return safe, nil
}
