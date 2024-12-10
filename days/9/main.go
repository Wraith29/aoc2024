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

type File struct {
	start, id, size int
	isFree          bool
}

func getFiles(input string) ([]File, error) {
	fileId := 0
	files := make([]File, 0)
	totalIndex := 0

	for idx, char := range strings.Split(input, "") {
		size, err := strconv.Atoi(char)
		if err != nil {
			return nil, err
		}

		if idx%2 != 0 {
			files = append(files, File{
				id:     0,
				start:  totalIndex,
				size:   size,
				isFree: true,
			})
			totalIndex += size
			continue
		}

		files = append(files, File{
			id:     fileId,
			start:  totalIndex,
			size:   size,
			isFree: false,
		})
		totalIndex += size

		fileId++
	}

	return files, nil
}

func (f *File) toIntArray() []int {
	result := make([]int, f.size)

	for idx := 0; idx < f.size; idx++ {
		if f.isFree {
			result[idx] = -1
		} else {
			result[idx] = f.id
		}
	}

	return result
}

const blockSize = 5

type Block struct {
	memory      [blockSize]int
	empty, full bool
}

func getBlocks(files []File) []Block {
	blocks := make([]Block, 0)

	memoryIndex := 0
	memorySlice := [blockSize]int{}

	for _, file := range files {
		values := file.toIntArray()

		for _, value := range values {
			if memoryIndex == blockSize {
				memoryIndex = 0

				blocks = append(blocks, Block{memory: memorySlice})

				memorySlice = [blockSize]int{}
			}

			memorySlice[memoryIndex] = value
			memoryIndex++
		}
	}

	for memoryIndex != blockSize {
		memorySlice[memoryIndex] = -1
		memoryIndex++
	}

	blocks = append(blocks, Block{memory: memorySlice})

	return blocks
}

func (b *Block) pop() int {
	for idx := blockSize - 1; idx >= 0; idx-- {
		if b.memory[idx] != -1 {
			result := b.memory[idx]
			b.memory[idx] = -1

			if idx == 0 {
				b.empty = true
			}

			return result
		}
	}

	return -1
}

// True means the insert has passed, false means it has failed, and is full
func (b *Block) insert(v int) bool {
	for idx := 0; idx < blockSize; idx++ {
		if b.memory[idx] == -1 {
			b.memory[idx] = v

			if idx == blockSize-1 {
				b.full = true
			}

			return true
		}
	}

	b.full = true

	return false
}

func part1(input string) (int, error) {
	result := 0

	files, err := getFiles(input)
	if err != nil {
		return -1, err
	}

	blocks := getBlocks(files)

	innerBlockIdx := 0
	outerBlockIdx := len(blocks) - 1

	for innerBlockIdx != outerBlockIdx {
		outerValue := blocks[outerBlockIdx].pop()

		if outerValue == -1 {
			outerBlockIdx--
			continue
		}

		for !blocks[innerBlockIdx].insert(outerValue) {
			innerBlockIdx++
		}
	}

	for blockIndex, block := range blocks {
		for index, value := range block.memory {
			if value < 0 {
				continue
			}

			totalIndex := blockIndex*blockSize + index

			result += totalIndex * value
		}
	}

	return result, nil
}

type Space struct {
	start, size int
}

type MemoryManager struct {
	memory []int
}

func newMemoryManager(files []File) MemoryManager {
	memory := make([]int, 0)
	totalIndex := 0

	for _, file := range files {
		memory = append(memory, file.toIntArray()...)
		totalIndex += file.size
	}

	return MemoryManager{
		memory: memory,
	}
}

func (m *MemoryManager) getNextSpace(start int) *Space {
	idx := start

	for idx < len(m.memory) {
		value := m.memory[idx]
		if value > 0 {
			idx++
			continue
		}

		spaceStartPos := idx
		size := 0

		for m.memory[spaceStartPos+size] < 0 && spaceStartPos+size < len(m.memory)-1 {
			size++
		}

		return &Space{
			start: spaceStartPos,
			size:  size,
		}
	}

	return nil
}

func (m *MemoryManager) insert(f *File, space *Space) {
	for idx := space.start; idx < space.start+f.size; idx++ {
		m.memory[idx] = f.id
	}

	for idx := f.start; idx < f.start+f.size; idx++ {
		m.memory[idx] = -1
	}
}

// 8515929533392 -> High
func part2(input string) (int, error) {
	result := 0

	files, err := getFiles(input)
	if err != nil {
		return -1, err
	}

	fileMap := make(map[int]*File)
	for _, file := range files {
		fileMap[file.id] = &file
	}

	mem := newMemoryManager(files)

	fileId := max(files[len(files)-1].id, files[len(files)-2].id)

	for fileId > 0 {
		file, _ := fileMap[fileId]
		if file.isFree {
			continue
		}

		pos := 0
		space := mem.getNextSpace(pos)

		for space != nil {
			if space.size < file.size || file.start < space.start {
				pos = space.size + space.start + 1
				space = mem.getNextSpace(pos)
				continue
			}

			mem.insert(file, space)

			break
		}

		fileId--
	}

	for idx, value := range mem.memory {
		if value >= 0 {
			result += value * idx
		}
	}

	return result, nil
}
