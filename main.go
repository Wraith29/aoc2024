package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("ERROR: Missing required positional arg: day")
		os.Exit(1)
	}

	day, err := strconv.Atoi(args[0])
	if err != nil {
		logErr(err)
	}

	err = os.Mkdir(fmt.Sprintf("days/%d", day), 0700)
	if err != nil {
		logErr(err)
	}

	if createMod(day) != nil {
		logErr(err)
	}

	os.Exit(0)
}

func logErr(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
	os.Exit(1)
}

func createMod(day int) error {
	cmd := exec.Command("go", "mod", "init", fmt.Sprintf("d%d", day))

	cmd.Dir = getWorkingDir(day)

	err := cmd.Run()

	if err != nil {
		return err
	}

	return createFiles(day)
}

func createFiles(day int) error {
	touchCmd := exec.Command("touch", "sample.txt", "input.txt")

	touchCmd.Dir = getWorkingDir(day)

	err := touchCmd.Run()

	if err != nil {
		return err
	}

	cpCmd := exec.Command("cp", "template/template.go", fmt.Sprintf("days/%d/main.go"))

	return cpCmd.Run()
}

func getWorkingDir(day int) string {
	return fmt.Sprintf("/home/iacna/dev/AOC2024/days/%d", day)
}
