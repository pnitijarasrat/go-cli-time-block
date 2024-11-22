package readtimeblock

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ReadTimeBlock() {
	readTodayFile()
	fmt.Println(" ")
}

func readTodayFile() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		fmt.Println(" ")
		return
	}

	homeDir := usr.HomeDir
	date := strings.Split(time.Now().String(), " ")[0]
	filePath := filepath.Join(homeDir, "time-block", date+".txt")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Time block doesn't exist, create new one.")
		fmt.Println(" ")
		return
	}

	scanner := bufio.NewScanner(file)
	fmt.Println(" ")
	for scanner.Scan() {
		handleLogging(scanner)
	}

}

func handleLogging(scanner *bufio.Scanner) error {

	line := scanner.Text()
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, line)
	if match {
		color.Red("Date: " + line)
		fmt.Println(" ")
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 4 {
		return errors.New("Something went wrong")
	}
	timeRange := parts[0] + " - " + parts[2]

	times := strings.Split(timeRange, " - ")
	if len(times) != 2 {
		return errors.New("Something went wrong")
	}
	start, err := strconv.Atoi(strings.Split(times[0], ".")[0])
	end, err := strconv.Atoi(strings.Split(times[1], ".")[0])

	if err != nil {
		return err
	}

	currentTime := time.Now().Hour()

	if start <= currentTime && end > currentTime {
		color.Cyan(line)
	} else {
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}
