package main

import (
	"fmt"
	"os"
	"slices"
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
	id, size int
	isFree   bool
}

func getFiles(input string) ([]File, error) {
	fileId := 0
	files := make([]File, 0)

	for idx, char := range strings.Split(input, "") {
		intVal, err := strconv.Atoi(char)
		if err != nil {
			return nil, err
		}

		if idx%2 != 0 {
			files = append(files, File{
				id:     0,
				size:   intVal,
				isFree: true,
			})
			continue
		}

		files = append(files, File{
			id:     fileId,
			size:   intVal,
			isFree: false,
		})

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

func (f *File) toSpacedIntArray() []int {
	digits := strconv.Itoa(f.id)
	result := make([]int, len(digits)*f.size)

	for idx := 0; idx < f.size; idx++ {
		for digitIdx, digit := range strings.Split(digits, "") {
			dig, err := strconv.Atoi(digit)
			if err != nil {
				panic(err)
			}

			result[idx*len(digits)+digitIdx] = dig
		}
	}

	return result
}

func draw(files []File) {
	for _, file := range files {
		if file.isFree {
			print(strings.Repeat(".", file.size))
		} else {
			print(strings.Repeat(strconv.Itoa(file.id), file.size))
		}
	}
	println()
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

type DynamicBlock struct {
	memory      []int
	initialId   int
	full, empty bool
}

func getDynamicBlocks(files []File) []DynamicBlock {
	blocks := make([]DynamicBlock, 0)

	for _, file := range files {
		block := DynamicBlock{
			memory:    file.toIntArray(),
			full:      !file.isFree,
			empty:     file.isFree,
			initialId: file.id,
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func repeat(value, size int) []int {
	result := make([]int, size)

	for idx := 0; idx < size; idx++ {
		result[idx] = value
	}

	return result
}

// We know size <= len(block.memory)
func (b *DynamicBlock) hasContiguousSpace(size int) (int, bool) {
	space := repeat(-1, size)

	for idx := 0; idx+size <= len(b.memory); idx++ {
		if slices.Equal(b.memory[idx:idx+size], space) {
			return idx, true
		}
	}
	return -1, false
}

func (b *DynamicBlock) insert(file *File, startIdx int) {
	for idx := startIdx; idx < startIdx+file.size; idx++ {
		b.memory[idx] = file.id
	}
}

func (b *DynamicBlock) clear() {
	b.memory = repeat(-1, len(b.memory))
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
		if file.isFree {
			continue
		}

		fileMap[file.id] = &file
	}

	// Max file id is either the last file or the second to last file
	// Dependent on the last file being empty or not
	fileId := max(files[len(files)-1].id, files[len(files)-2].id)

	blocks := getDynamicBlocks(files)

	for fileId > 0 {
		for idx, block := range blocks {
			if block.full || fileId <= idx {
				continue
			}

			file, _ := fileMap[fileId]
			// memory := file.toSpacedIntArray()

			if file.size > len(block.memory) {
				continue
			}

			startIdx, fits := block.hasContiguousSpace(file.size)
			if !fits {
				continue
			}

			block.insert(file, startIdx)

			for idx := len(blocks) - 1; idx >= 0; idx-- {
				if blocks[idx].initialId == fileId {
					blocks[idx].clear()
					break
				}
			}

			break
		}

		fileId--
	}

	for _, b := range blocks {
		for _, m := range b.memory {
			if m < 0 {
				print(".")
			} else {
				print(m)
			}
		}
	}
	println()

	totalIdx := 0

	for _, block := range blocks {
		for _, value := range block.memory {
			if value > 0 {
				result += value * totalIdx
			}
			totalIdx++
		}
	}

	return result, nil
}
