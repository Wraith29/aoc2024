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

type Step struct{ x, y int }

func (s *Step) invert() {
	s.x = -s.x
	s.y = -s.y
}

type Node struct {
	x, y      int
	frequency string
	antiNode  bool
}

func (n *Node) getStep(o *Node) Step {
	return Step{n.x - o.x, n.y - o.y}
}

type Grid struct {
	gridSize int
	grid     [][]*Node
	freqMap  map[string][]*Node
}

func newGrid(input string) Grid {
	gridSize := strings.Index(input, "\n")
	grid := make([][]*Node, 0)
	freqMap := make(map[string][]*Node)

	for y, line := range strings.Split(input, "\n") {
		row := make([]*Node, 0)

		for x, column := range strings.Split(line, "") {
			frequency := ""
			if column != "." {
				frequency = column
			}

			node := &Node{
				x, y, frequency, false,
			}

			row = append(row, node)

			if frequency == "" {
				continue
			}

			freqNodes, found := freqMap[frequency]
			if !found {
				freqNodes = make([]*Node, 0)
			}

			freqNodes = append(freqNodes, node)
			freqMap[frequency] = freqNodes
		}

		grid = append(grid, row)
	}

	return Grid{
		gridSize,
		grid,
		freqMap,
	}
}

func (g *Grid) antiNodeCount() int {
	result := 0

	for _, row := range g.grid {
		for _, node := range row {
			if node.antiNode {
				result++
			}
		}
	}

	return result
}

func (g *Grid) getFromStep(n *Node, s Step) *Node {
	x := n.x + s.x
	y := n.y + s.y

	if x < 0 || x >= g.gridSize || y < 0 || y >= g.gridSize {
		return nil
	}

	return g.grid[y][x]
}

func (g *Grid) draw() {
	for _, row := range g.grid {
		for _, node := range row {
			if node.antiNode {
				print("#")
			} else if node.frequency != "" {
				print(node.frequency)
			} else {
				print(".")
			}
		}

		println()
	}
}

func part1(input string) (int, error) {
	grid := newGrid(input)

	for _, nodes := range grid.freqMap {
		for i, outerNode := range nodes {
			for j, innerNode := range nodes {
				if i == j {
					continue
				}

				step := outerNode.getStep(innerNode)

				next := grid.getFromStep(outerNode, step)
				if next != nil {
					next.antiNode = true
				}
			}
		}
	}

	return grid.antiNodeCount(), nil
}

func part2(input string) (int, error) {
	grid := newGrid(input)

	for _, nodes := range grid.freqMap {
		for i, outerNode := range nodes {
			for j, innerNode := range nodes {
				if i == j {
					continue
				}

				step := outerNode.getStep(innerNode)

				outerNode.antiNode = true

				next := grid.getFromStep(outerNode, step)
				for next != nil {
					next.antiNode = true
					next = grid.getFromStep(next, step)
				}
			}
		}
	}

	return grid.antiNodeCount(), nil
}
