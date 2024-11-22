package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	createtimeblock "timeblock/create-time-block"
	readtimeblock "timeblock/read-time-block"

	"github.com/fatih/color"
)

func main() {

	if len(os.Args) >= 3 {
		if os.Args[1] == "v" {
			if os.Args[2] == "tmr" {
				color.Red("Viewing tomorrow timeblock is not available now.")
			}
		} else {
			createtimeblock.CreateTimeBlock()
		}
	}

	if len(os.Args) >= 2 {
		if os.Args[1] == "v" {
			readtimeblock.ReadTimeBlock()
		} else if os.Args[1] == "sync" {
			gitPull()
		} else {
			color.Red("Invalid Input")
		}
	} else {
		createtimeblock.CreateTimeBlock()
	}

}

func gitPull() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	homeDir := usr.HomeDir

	commands := [][]string{
		{"git", "pull"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = filepath.Join(homeDir, "time-block")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running %s: %v\n", strings.Join(cmdArgs, " "), err)
			return
		}
	}

}

func getCommandInput() (string, error) {

	fmt.Println("Welcome to Time Block Scheduler")
	fmt.Println("-- n for new time block")
	fmt.Println("-- v to view today time block")
	fmt.Println(" ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input == "n" || input == "v" {
		return input, nil
	}

	return "", errors.New("Invalid Input")
}

func getCommand() (string, error) {

	var command string
	var err error

	for {
		command, err = getCommandInput()
		if err != nil {
			fmt.Println(err)
			fmt.Println(" ")
		} else {
			break
		}
	}

	return command, nil
}
