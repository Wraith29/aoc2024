package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	re "github.com/dlclark/regexp2"
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

var (
	mulPtn    = re.MustCompile(`mul\(([0-9]*),([0-9]*)\)`, 0)
	togglePtn = re.MustCompile(`(do)\(\)|(don't)\(\)`, 0)
)

func findAll(ptn *re.Regexp, input string) ([]*re.Match, error) {
	results := make([]*re.Match, 0)

	m, err := ptn.FindStringMatch(input)
	if err != nil {
		return nil, err
	}

	for m != nil {
		results = append(results, m)

		m, err = ptn.FindNextMatch(m)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func calculateMul(match *re.Match) (int, error) {
	left, err := strconv.Atoi(match.GroupByNumber(1).String())
	if err != nil {
		return -1, err
	}

	right, err := strconv.Atoi(match.GroupByNumber(2).String())
	if err != nil {
		return -1, err
	}

	return left * right, nil
}

func part1(input string) (int, error) {
	matches, err := findAll(mulPtn, input)

	if err != nil {
		return -1, err
	}

	result := 0

	for _, match := range matches {
		val, err := calculateMul(match)
		if err != nil {
			return -1, err
		}

		result += val
	}

	return result, nil
}

func part2(input string) (int, error) {
	toggles, err := findAll(togglePtn, input)

	if err != nil {
		return -1, err
	}

	result := 0
	enabled := true
	idx := 0

	for _, toggle := range toggles {
		if enabled {
			inp := input[idx : toggle.Index+toggle.Length]

			muls, err := findAll(mulPtn, inp)

			if err != nil {
				return -1, err
			}

			for _, mul := range muls {
				val, err := calculateMul(mul)

				if err != nil {
					return -1, nil
				}
				result += val
			}
		}

		idx = toggle.Index + toggle.Length

		enabled = toggle.String() == "do()"
	}

	if enabled {
		finalText := input[idx:]

		muls, err := findAll(mulPtn, finalText)

		if err != nil {
			return -1, err
		}

		for _, mul := range muls {
			val, err := calculateMul(mul)

			if err != nil {
				return -1, nil
			}
			result += val
		}
	}

	return result, nil
}
