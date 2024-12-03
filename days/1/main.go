package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
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

type lists struct {
	left, right []int
}

func getLists(input string) (*lists, error) {

	leftSide := make([]int, 0)
	rightSide := make([]int, 0)

	for _, line := range strings.Split(input, "\n") {
		values := strings.Split(line, "   ")

		left, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}

		leftSide = append(leftSide, left)

		right, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		rightSide = append(rightSide, right)
	}

	if len(leftSide) != len(rightSide) {
		return nil, errors.New("invalid input")
	}

	return &lists{
		left:  leftSide,
		right: rightSide,
	}, nil
}

func part1(input string) (int, error) {
	lists, err := getLists(input)

	if err != nil {
		return -1, err
	}

	leftSide := lists.left
	rightSide := lists.right

	slices.Sort(leftSide)
	slices.Sort(rightSide)

	dist := 0

	for idx, left := range leftSide {
		right := rightSide[idx]

		if left > right {
			dist += left - right
		} else {
			dist += right - left
		}
	}

	return dist, nil
}

func part2(input string) (int, error) {
	lists, err := getLists(input)
	if err != nil {
		return -1, err
	}

	simScore := 0

	scoreMap := make(map[int]int)

	for _, right := range lists.right {
		val, found := scoreMap[right]
		if found {
			scoreMap[right] = val + 1
		} else {
			scoreMap[right] = 1
		}
	}

	for _, left := range lists.left {
		count, found := scoreMap[left]
		if !found {
			continue
		}

		simScore += left * count
	}

	return simScore, nil
}
