package main

import (
	"fmt"
	"os"
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

func part1(input string) (int, error) {
	return 0, nil
}

func part2(input string) (int, error) {
	return 0, nil
}
