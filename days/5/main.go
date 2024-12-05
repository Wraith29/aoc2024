package main

import (
	"fmt"
	"maps"
	"math"
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

type Page struct {
	value   int
	before  []int
	printed bool
}

// Page Number to Must be Before
type Instructions map[int]*Page

func loadInstructions(input string) (Instructions, error) {
	ins := make(Instructions)

	for _, line := range strings.Split(input, "\n") {
		items := strings.Split(strings.Trim(line, "\n"), "|")

		pageValue, err := strconv.Atoi(items[0])
		if err != nil {
			return nil, err
		}

		before, err := strconv.Atoi(items[1])
		if err != nil {
			return nil, err
		}

		val, exists := ins[pageValue]

		if exists {
			val.before = append(val.before, before)

			continue
		}

		ins[pageValue] = &Page{
			value:   pageValue,
			before:  []int{before},
			printed: false,
		}
	}

	return ins, nil
}

func (i Instructions) reset() {
	for _, value := range i {
		value.printed = false
	}
}

func (i Instructions) copy() Instructions {
	return maps.Clone(i)
}

func (i Instructions) isPrinted(pageNum int) bool {
	page, exists := i[pageNum]
	if !exists {
		return false
	}

	return page.printed
}

func isLineCorrect(instructions Instructions, pages []int) bool {
	instructions.reset()

	for idx, pageNumber := range pages {
		page, exists := instructions[pageNumber]

		if !exists && idx == len(pages)-1 {
			return true
		} else if !exists {
			return false
		}

		for _, child := range page.before {
			if instructions.isPrinted(child) {
				return false
			}
		}

		page.printed = true
	}

	return true
}

func reorderLine(instructions Instructions, pages []int) []int {
	instructions.reset()

	for idx, pageNumber := range pages {
		page, exists := instructions[pageNumber]

		if !exists && idx == len(pages)-1 {
			return pages
		} else if !exists {
			pages[idx], pages[idx+1] = pages[idx+1], pages[idx]
			break
		}

		for _, child := range page.before {
			if instructions.isPrinted(child) {
				pages[idx-1], pages[idx] = pages[idx], pages[idx-1]
				break
			}
		}
		page.printed = true
	}

	if !isLineCorrect(instructions, pages) {
		return reorderLine(instructions, pages)
	}

	return pages
}

func part1(input string) (int, error) {
	result := 0

	splitIdx := strings.Index(input, "\n\n")

	instructions, err := loadInstructions(input[:splitIdx])
	if err != nil {
		return -1, err
	}

	manuals := input[splitIdx+2:]

	for _, line := range strings.Split(manuals, "\n") {
		pageNumbers := strings.Split(strings.Trim(line, "\n"), ",")
		pages := make([]int, 0)

		for _, pageNumber := range pageNumbers {
			value, err := strconv.Atoi(pageNumber)
			if err != nil {
				return -1, err
			}

			pages = append(pages, value)
		}

		if isLineCorrect(instructions, pages) {
			result += pages[int(math.Ceil(float64(len(pages)/2)))]
		}
	}

	return result, nil
}

func part2(input string) (int, error) {
	result := 0

	splitIdx := strings.Index(input, "\n\n")

	instructions, err := loadInstructions(input[:splitIdx])
	if err != nil {
		return -1, err
	}

	manuals := input[splitIdx+2:]

	for _, line := range strings.Split(manuals, "\n") {
		pageNumbers := strings.Split(strings.Trim(line, "\n"), ",")
		pages := make([]int, 0)

		for _, pageNumber := range pageNumbers {
			value, err := strconv.Atoi(pageNumber)
			if err != nil {
				return -1, err
			}

			pages = append(pages, value)
		}

		// If they're right straight away, ignore
		if isLineCorrect(instructions, pages) {
			continue
		}
		pages = reorderLine(instructions, pages)

		result += pages[int(math.Ceil(float64(len(pages)/2)))]
	}

	return result, nil
}
