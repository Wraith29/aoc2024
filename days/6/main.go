package main

import (
	"errors"
	"fmt"
	"math"
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

type Direction struct{ x, y int }

var (
	up    = Direction{0, -1}
	down  = Direction{0, 1}
	left  = Direction{-1, 0}
	right = Direction{1, 0}
)

func nextDirection(d Direction) Direction {
	if d == left {
		return up
	} else if d == up {
		return right
	} else if d == right {
		return down
	} else {
		return left
	}
}

type Guard struct {
	x, y, steps  int
	facing       Direction
	obstructions []*Cell
}

type Cell struct {
	x, y                                    int
	visited, obstructed, crossedX, crossedY bool
}

type Grid struct {
	cells      [][]*Cell
	lineLength int
	guard      Guard
}

func newGrid(input string) Grid {
	lineLength := strings.Index(input, "\n")

	guard := Guard{
		x:            0,
		y:            0,
		steps:        0,
		facing:       up,
		obstructions: make([]*Cell, 0),
	}
	cells := make([][]*Cell, 0)

	for y, line := range strings.Split(input, "\n") {
		row := make([]*Cell, 0)

		for x, tile := range strings.Split(line, "") {
			if tile == "^" {
				guard.x = x
				guard.y = y
			}

			row = append(row, &Cell{
				x:          x,
				y:          y,
				visited:    false,
				obstructed: tile == "#",
				crossedX:   false,
				crossedY:   false,
			})
		}

		cells = append(cells, row)
	}

	return Grid{
		cells,
		lineLength,
		guard,
	}
}

func (g *Grid) visit(x, y int) {
	cell := g.cells[y][x]

	if g.guard.facing == up || g.guard.facing == down {
		cell.crossedY = true
	} else {
		cell.crossedX = true
	}

	cell.visited = true
}

func (g *Grid) getNextCell() *Cell {
	x := g.guard.x + g.guard.facing.x
	y := g.guard.y + g.guard.facing.y

	if x < 0 || x >= g.lineLength || y < 0 || y >= g.lineLength {
		return nil
	}

	return g.cells[y][x]
}

func (g *Grid) moveGuard() (bool, error) {
	g.visit(g.guard.x, g.guard.y)

	nextCell := g.getNextCell()
	// Just exit, nextCell is nil so we are out of the grid
	if nextCell == nil {
		return false, nil
	}

	for nextCell.obstructed {
		g.guard.facing = nextDirection(g.guard.facing)
		g.visit(g.guard.x, g.guard.y)
		g.guard.obstructions = append(g.guard.obstructions, nextCell)

		nextCell = g.getNextCell()
		if nextCell == nil {
			return false, nil
		}
	}

	g.guard.x = nextCell.x
	g.guard.y = nextCell.y
	g.guard.steps++

	if g.guard.steps > int(math.Pow(float64(g.lineLength), 2)) {
		return false, errors.New("stuck in a loop")
	}

	return true, nil
}

func (g *Grid) draw() {
	for _, row := range g.cells {
		for _, cell := range row {
			if cell.crossedX && cell.crossedY {
				print("+")
			} else if cell.crossedX {
				print("-")
			} else if cell.crossedY {
				print("|")
			} else if cell.obstructed {
				print("#")
			} else {
				print(".")
			}
		}
		println()
	}
}

func (g *Grid) getCell(x, y int) *Cell {
	if x < 0 || x >= g.lineLength || y < 0 || y >= g.lineLength {
		return nil
	}

	return g.cells[y][x]
}

func part1(input string) (int, error) {
	result := 0
	grid := newGrid(input)

	exit, err := grid.moveGuard()
	for exit && err == nil {
		exit, err = grid.moveGuard()
	}

	for _, row := range grid.cells {
		for _, cell := range row {
			if cell.visited {
				result++
			}
		}
	}

	return result, nil
}

// 1732 -> too high
func part2(input string) (int, error) {
	result := 0

	lineLength := strings.Index(input, "\n")

	for y := 0; y < lineLength; y++ {
		for x := 0; x < lineLength; x++ {
			grid := newGrid(input)

			if y == grid.guard.y && x == grid.guard.x {
				continue
			}

			cell := grid.getCell(x, y)

			if !cell.obstructed {
				cell.obstructed = true

				exit, err := grid.moveGuard()
				for exit && err == nil {
					exit, err = grid.moveGuard()
				}

				if err != nil {
					result++
				}
			}

		}
	}

	return result, nil
}
