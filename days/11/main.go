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

var cache = make(map[int][]int)

func transformStone(stone int) ([]int, error) {
	if stone == 0 {
		return []int{1}, nil
	}

	digits := strconv.Itoa(stone)
	if len(digits)%2 == 0 {
		left, err := strconv.Atoi(digits[:len(digits)/2])
		if err != nil {
			return nil, err
		}

		right, err := strconv.Atoi(digits[len(digits)/2:])
		if err != nil {
			return nil, err
		}

		return []int{left, right}, nil
	}

	return []int{stone * 2024}, nil
}

func getTransformation(stone int) ([]int, error) {
	if result, ok := cache[stone]; ok {
		return result, nil
	}

	result, err := transformStone(stone)
	if err != nil {
		return nil, err
	}

	cache[stone] = result
	return result, nil
}

func getStones(input string) (map[int]int, error) {
	stones := make(map[int]int)

	for _, stoneValue := range strings.Split(strings.Trim(input, " "), " ") {
		value, err := strconv.Atoi(stoneValue)
		if err != nil {
			return nil, err
		}

		stones[value]++
	}

	return stones, nil
}

func part1(input string) (int, error) {
	result := 0

	stones, err := getStones(input)
	if err != nil {
		return -1, err
	}

	for count := 0; count < 25; count++ {
		newStones := make(map[int]int)

		for stone, value := range stones {
			transformedStones, err := getTransformation(stone)
			if err != nil {
				return -1, err
			}

			for _, newStone := range transformedStones {
				newStones[newStone] += value
			}
			stones = newStones
		}
	}

	for _, count := range stones {
		result += count
	}

	return result, nil
}

func part2(input string) (int, error) {
	result := 0

	stones, err := getStones(input)
	if err != nil {
		return -1, err
	}

	for count := 0; count < 75; count++ {
		newStones := make(map[int]int)

		for stone, value := range stones {
			transformedStones, err := getTransformation(stone)
			if err != nil {
				return -1, err
			}

			for _, newStone := range transformedStones {
				newStones[newStone] += value
			}
			stones = newStones
		}
	}

	for _, count := range stones {
		result += count
	}

	return result, nil
}
