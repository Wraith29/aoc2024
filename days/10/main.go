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

type Cell struct {
	x, y, height int
	reached      bool
}

type Grid struct {
	gridSize int
	cells    [][]*Cell
}

func newGrid(input string) (*Grid, error) {
	gridSize := strings.Index(input, "\n")
	cells := make([][]*Cell, 0)

	for y, line := range strings.Split(input, "\n") {
		row := make([]*Cell, 0)

		for x, col := range strings.Split(line, "") {
			height, err := strconv.Atoi(col)
			if err != nil {
				return nil, err
			}

			row = append(row, &Cell{x, y, height, false})
		}

		cells = append(cells, row)
	}

	return &Grid{gridSize, cells}, nil
}

func (g *Grid) reset() {
	for _, row := range g.cells {
		for _, cell := range row {
			cell.reached = false
		}
	}
}

func (g *Grid) getCellNeighbours(c *Cell) []*Cell {
	neighbours := make([]*Cell, 0)

	coords := []struct{ x, y int }{
		{c.x, c.y - 1},
		{c.x - 1, c.y}, {c.x + 1, c.y},
		{c.x, c.y + 1},
	}

	for _, coord := range coords {
		if coord.x < 0 || coord.x >= g.gridSize || coord.y < 0 || coord.y >= g.gridSize {
			continue
		}

		neighbours = append(neighbours, g.cells[coord.y][coord.x])
	}

	return neighbours
}

func (g *Grid) checkTrail(c *Cell) int {
	validRoutes := 0

	for _, neighbour := range g.getCellNeighbours(c) {
		if neighbour.height == c.height+1 && neighbour.height == 9 && !neighbour.reached {
			validRoutes++
			neighbour.reached = true
			continue
		} else if neighbour.height == c.height+1 {
			validRoutes += g.checkTrail(neighbour)
		}
	}

	return validRoutes
}

func (g *Grid) checkDistinctTrails(c *Cell) int {
	validRoutes := 0

	for _, neighbour := range g.getCellNeighbours(c) {
		if neighbour.height == c.height+1 && neighbour.height == 9 {
			validRoutes++
			continue
		} else if neighbour.height == c.height+1 {
			validRoutes += g.checkDistinctTrails(neighbour)
		}
	}

	return validRoutes
}

func part1(input string) (int, error) {
	result := 0

	grid, err := newGrid(input)
	if err != nil {
		return -1, err
	}

	for _, row := range grid.cells {
		for _, cell := range row {
			if cell.height != 0 {
				continue
			}
			grid.reset()
			result += grid.checkTrail(cell)
		}
	}

	return result, nil
}

func part2(input string) (int, error) {
	result := 0

	grid, err := newGrid(input)
	if err != nil {
		return -1, err
	}

	for _, row := range grid.cells {
		for _, cell := range row {
			if cell.height != 0 {
				continue
			}
			grid.reset()
			result += grid.checkDistinctTrails(cell)
		}
	}

	return result, nil
}
