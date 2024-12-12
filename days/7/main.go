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

type Equation struct {
	value    int
	children []int
}

type Operator struct {
	rep  string
	exec func(int, int) int
}

var (
	plus = Operator{
		rep: "+",
		exec: func(a, b int) int {
			return a + b
		},
	}
	mult = Operator{
		rep: "*",
		exec: func(a, b int) int {
			return a * b
		},
	}
	concat = Operator{
		rep: "||",
		exec: func(a, b int) int {
			sb := strings.Builder{}

			sb.WriteString(strconv.Itoa(a) + strconv.Itoa(b))

			val, err := strconv.Atoi(sb.String())

			if err != nil {
				panic(err)
			}

			return val
		},
	}
)

func generate(max int, operations []Operator) [][]Operator {
	var gen func([]Operator, int, int) [][]Operator

	gen = func(items []Operator, max, depth int) [][]Operator {
		if items == nil {
			items = make([]Operator, 0)
		}

		if depth == max {
			// Copy the slice, so that the items slice doesn't then get modified later
			final := make([]Operator, len(items))
			copy(final, items)

			return [][]Operator{final}
		}

		result := make([][]Operator, 0)
		for _, op := range operations {
			current := append(items, op)

			result = append(result, gen(current, max, depth+1)...)
		}

		return result
	}

	return gen(make([]Operator, 0), max, 0)
}

func (e *Equation) isPossible(operators []Operator) bool {
	// Amount of operations that need to happen
	opCount := len(e.children) - 1

	operationCombos := generate(opCount, operators)

	for _, opSet := range operationCombos {
		sum := e.children[0]

		for i := 0; i < len(opSet); i++ {
			sum = opSet[i].exec(sum, e.children[i+1])
		}

		if sum == e.value {
			return true
		}
	}

	return false
}

func getEquations(input string) ([]Equation, error) {
	results := make([]Equation, 0)

	for _, val := range strings.Split(input, "\n") {
		splitIdx := strings.Index(val, ":")
		equationValue, err := strconv.Atoi(val[:splitIdx])
		if err != nil {
			return nil, err
		}

		children := make([]int, 0)

		for _, child := range strings.Split(strings.Trim(val[splitIdx+1:], " "), " ") {
			childVal, err := strconv.Atoi(child)
			if err != nil {
				return nil, err
			}

			children = append(children, childVal)
		}

		results = append(results, Equation{
			value:    equationValue,
			children: children,
		})
	}

	return results, nil
}

// 2399558932629 - Too low
func part1(input string) (int, error) {
	result := 0

	equations, err := getEquations(input)
	if err != nil {
		return -1, err
	}

	operators := []Operator{plus, mult}

	for _, eq := range equations {
		if eq.isPossible(operators) {
			result += eq.value
		}
	}

	return result, nil
}

func part2(input string) (int, error) {
	result := 0

	equations, err := getEquations(input)
	if err != nil {
		return -1, err
	}

	operators := []Operator{plus, mult, concat}

	for _, eq := range equations {
		if eq.isPossible(operators) {
			result += eq.value
		}
	}

	return result, nil
}
