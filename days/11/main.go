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

func insert(original []int, splitIndex, left, right int) []int {
	result := make([]int, 0)

	leftHalf := original[:splitIndex]

	rightHalf := original[splitIndex+1:]

	result = append(result, leftHalf...)
	result = append(result, left, right)
	result = append(result, rightHalf...)

	return result
}

type Transformer struct {
	fn    func(int) ([]int, error)
	cache map[int][]int
}

func newTransformer(fn func(int) ([]int, error)) Transformer {
	return Transformer{
		fn:    fn,
		cache: make(map[int][]int),
	}
}

func (t *Transformer) transform(stone int) ([]int, error) {
	if value, ok := t.cache[stone]; ok {
		return value, nil
	}

	result, err := t.fn(stone)
	if err != nil {
		return nil, err
	}
	t.cache[stone] = result

	return result, nil
}

type Stone struct {
	idx    int
	buffer []int
}

func getStones(input string) ([]*Stone, error) {
	result := make([]*Stone, 0)

	for idx, stone := range strings.Split(strings.Trim(input, " "), " ") {
		value, err := strconv.Atoi(stone)
		if err != nil {
			return nil, err
		}

		result = append(result, &Stone{
			idx:    idx,
			buffer: []int{value},
		})
	}

	return result, nil
}

func (s *Stone) applyTransformation(t *Transformer, count int) error {
	for n := 0; n < count; n++ {
		println(n, len(s.buffer))

		for idx := 0; idx < len(s.buffer); idx++ {
			transformed, err := t.transform(s.buffer[idx])
			if err != nil {
				return err
			}

			if len(transformed) == 1 {
				s.buffer[idx] = transformed[0]
			} else {
				s.buffer = insert(s.buffer, idx, transformed[0], transformed[1])
				idx++
			}
		}
	}

	return nil
}

func transformStone(stone int) ([]int, error) {
	if stone == 0 {
		return []int{1}, nil
	}

	digits := strconv.Itoa(stone)

	if len(digits)%2 == 0 {
		halfIdx := len(digits) / 2

		left, err := strconv.Atoi(digits[:halfIdx])
		if err != nil {
			return nil, err
		}

		right, err := strconv.Atoi(digits[halfIdx:])
		if err != nil {
			return nil, err
		}

		return []int{left, right}, nil
	}

	return []int{stone * 2024}, nil
}

func part1(input string) (int, error) {
	result := 0

	stones, err := getStones(input)
	if err != nil {
		return -1, err
	}

	transformer := newTransformer(transformStone)

	for _, stone := range stones {
		if err = stone.applyTransformation(&transformer, 25); err != nil {
			return -1, err
		}

		result += len(stone.buffer)
	}

	return result, nil
}

func part2(input string) (int, error) {
	result := 0

	stones, err := getStones(input)
	if err != nil {
		return -1, err
	}

	transformer := newTransformer(transformStone)

	for _, stone := range stones {
		if err = stone.applyTransformation(&transformer, 75); err != nil {
			return -1, err
		}

		result += len(stone.buffer)
	}

	return result, nil
}
