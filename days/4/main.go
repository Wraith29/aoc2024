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

	fmt.Printf("Part 1: %d\n", p1)

	p2, err := part2(input)

	if err != nil {
		logErr(err)
	}

	fmt.Printf("Part 2: %d\n", p2)
}

type Cell struct {
	x, y int
	char string
}

type Grid struct {
	lineLength int
	cells      [][]Cell
}

func createGrid(input string) Grid {
	lineLength := strings.Index(input, "\n")
	cells := make([][]Cell, 0)

	for y, line := range strings.Split(input, "\n") {
		row := make([]Cell, 0)

		for x, char := range strings.Split(line, "") {
			row = append(row, Cell{x, y, char})
		}

		cells = append(cells, row)
	}

	return Grid{lineLength, cells}
}

func (c *Cell) getNeighbours(g Grid) []Cell {
	neighbours := make([]Cell, 0)

	indexes := []struct{ x, y int }{
		{c.x - 1, c.y - 1}, {c.x, c.y - 1}, {c.x + 1, c.y - 1},
		{c.x - 1, c.y}, {c.x + 1, c.y},
		{c.x - 1, c.y + 1}, {c.x, c.y + 1}, {c.x + 1, c.y + 1},
	}

	for _, index := range indexes {
		if index.x < 0 || index.x >= g.lineLength || index.y < 0 || index.y >= g.lineLength {
			continue
		}

		neighbours = append(neighbours, g.cells[index.y][index.x])
	}

	return neighbours
}

func (c Cell) step(o Cell) (int, int) {
	return c.x - o.x, c.y - o.y
}

type coord struct{ x, y int }

func (g Grid) getFromStep(c Cell, step coord, gap int) *Cell {
	x := c.x - step.x*gap
	y := c.y - step.y*gap

	if x < 0 || x >= g.lineLength || y < 0 || y >= g.lineLength {
		return nil
	}

	return &g.cells[y][x]
}

func part1(input string) (int, error) {
	result := 0

	grid := createGrid(input)

	for _, row := range grid.cells {
		for _, item := range row {
			if item.char != "X" {
				continue
			}

			neighbours := item.getNeighbours(grid)

			for _, nb := range neighbours {
				if nb.char != "M" {
					continue
				}

				x, y := item.step(nb)

				aCell := grid.getFromStep(nb, coord{x, y}, 1)

				if aCell == nil || aCell.char != "A" {
					continue
				}

				sCell := grid.getFromStep(nb, coord{x, y}, 2)

				if sCell != nil && sCell.char == "S" {
					result += 1
				}

			}
		}
	}

	return result, nil
}

var (
	topLeft  = coord{x: -1, y: -1}
	topRight = coord{x: 1, y: -1}
	btmLeft  = coord{x: -1, y: 1}
	btmRight = coord{x: 1, y: 1}
)

func part2(input string) (int, error) {
	result := 0
	grid := createGrid(input)

	for _, row := range grid.cells {
		for _, cell := range row {
			if cell.char != "A" {
				continue
			}

			tl := grid.getFromStep(cell, topLeft, 1)
			tr := grid.getFromStep(cell, topRight, 1)
			bl := grid.getFromStep(cell, btmLeft, 1)
			br := grid.getFromStep(cell, btmRight, 1)

			// These could probably all be merged into 1 giant IF, but I don't care
			if tl == nil || tr == nil || bl == nil || br == nil {
				continue
			}

			if tl.char == "X" || tl.char == "A" ||
				tr.char == "X" || tr.char == "A" ||
				bl.char == "X" || bl.char == "A" ||
				br.char == "X" || br.char == "A" {
				continue
			}

			if tl.char == br.char || tr.char == bl.char {
				continue
			}

			result += 1
		}
	}

	return result, nil
}
