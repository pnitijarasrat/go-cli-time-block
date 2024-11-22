package createtimeblock

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

type TimeBlock struct {
	Start string
	End   string
	Task  string
}

func CreateTimeBlock() {
	date, _ := GetDate()

	fmt.Println("This is time block for:", date)
	fmt.Println(" ")

	timeRange, _ := getTimeRange()

	fmt.Println("Your day will start at " + timeRange[0] + ".00 and end at " + timeRange[1] + ".00")
	fmt.Println(" ")

	schedule, _ := getSchedule(timeRange)

	createTimeBlockFile(date, schedule)

}

func createTimeBlockFile(date string, schedule map[string]string) {

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	homeDir := usr.HomeDir

	filePath := filepath.Join(homeDir, "time-block", date+".txt")

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println(err)
			fmt.Println(" ")
			return
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		fmt.Println(" ")
		return
	}
	defer file.Close()

	file.WriteString(date + "\n\n")

	timeBlocks := make([]TimeBlock, 0, len(schedule))

	for timeRange, task := range schedule {
		times := strings.Split(timeRange, " - ")
		start := formatTime(times[0])
		end := formatTime(times[1])
		timeBlocks = append(timeBlocks, TimeBlock{Start: start, End: end, Task: task})
	}

	sort.Slice(timeBlocks, func(i, j int) bool {
		startTimeI := parseTime(timeBlocks[i].Start)
		startTimeJ := parseTime(timeBlocks[j].Start)
		return startTimeI < startTimeJ
	})

	for _, block := range timeBlocks {
		_, err := file.WriteString(block.Start + " - " + block.End + "  " + block.Task + "\n")
		if err != nil {
			fmt.Println(err)
			fmt.Println(" ")
			return
		}
	}

	color.Cyan("Your time block has been created for " + date)
	fmt.Println(" ")

	HandleGitPush(homeDir, date)
}

func HandleGitPush(homeDir string, date string) {

	filePath := filepath.Join(homeDir, "time-block", date+".txt")

	commands := [][]string{
		{"git", "add", filePath},
		{"git", "commit", "-m", "Add time block file for " + date},
		{"git", "push"},
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

	fmt.Println(" ")
	color.Green("Files Synced")
}

func parseTime(timeStr string) int {
	parts := strings.Split(timeStr, ".")
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	return hours*60 + minutes
}

func formatTime(timeStr string) string {
	parts := strings.Split(timeStr, ".")
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	return fmt.Sprintf("%02d.%02d", hours, minutes)
}

func getDateInput(prompt string, reader *bufio.Reader) (string, error) {

	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "today" {
		input = strings.Split(time.Now().String(), " ")[0]
	} else if input == "tomorrow" {
		input = strings.Split(time.Now().AddDate(0, 0, 1).String(), " ")[0]
	}

	for _, char := range input {
		if unicode.IsLetter(char) {
			return "", errors.New("Date Input should not contain alphabets, try today")
		}
	}

	return input, nil
}

func GetDate() (string, error) {

	reader := bufio.NewReader(os.Stdin)
	var date string
	var err error

	for {
		date, err = getDateInput("What's the date: ", reader)
		if err != nil {
			color.Red(err.Error())
			fmt.Println(" ")
		} else {
			break
		}
	}
	return date, err
}

func getTimeRangeInput(propmt string, reader *bufio.Reader) ([]string, error) {
	fmt.Print(propmt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	arr := strings.Split(input, " ")

	if len(arr) < 2 {
		return nil, errors.New("Time Range must contain at least 2 arguments")
	}

	left, err := strconv.Atoi(arr[0])
	right, err := strconv.Atoi(arr[1])

	if err != nil {
		return nil, errors.New("Invalid Time Range Input, should be INTERGER")
	}

	if right < left {
		return nil, errors.New("Start should be less than Stop")
	}

	return arr, nil
}

func getTimeRange() ([]string, error) {

	reader := bufio.NewReader(os.Stdin)
	var timeRange []string
	var err error

	for {
		timeRange, err = getTimeRangeInput("What's the time range: ", reader)
		if err != nil {
			color.Red(err.Error())
			fmt.Println(" ")
		} else {
			break
		}
	}

	return timeRange, err
}

func getSchedule(arr []string) (map[string]string, error) {

	reader := bufio.NewReader(os.Stdin)
	start, err := strconv.Atoi(arr[0])
	stop, err := strconv.Atoi(arr[1])
	var task string
	var howLongInt int
	schedule := make(map[string]string)

	if err != nil {
		color.Red(err.Error())
		fmt.Println(" ")
	}

	for i := start; i < stop; i++ {
		startAt := strconv.Itoa(i)
		for {
			prompt := "What do you wanna do at " + startAt + ".00: "
			task, howLongInt, err = getScheduleInput(prompt, reader)
			if err != nil {
				color.Red(err.Error())
				fmt.Println(" ")
			} else {
				break
			}
		}
		i += howLongInt - 1
		endAt := strconv.Itoa(i + 1)
		schedule[startAt+".00 - "+endAt+".00"] = task
		fmt.Println(" ")
	}

	fmt.Println("All Booked")
	fmt.Println(" ")

	return schedule, nil
}

func getScheduleInput(propmt string, reader *bufio.Reader) (string, int, error) {

	fmt.Println(propmt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return "", 0, errors.New("Task might not be empty string")
	}

	fmt.Println("For how long (hours): ")
	howLong, _ := reader.ReadString('\n')
	howLong = strings.TrimSpace(howLong)

	for _, char := range howLong {
		if unicode.IsLetter(char) {
			return "", 0, errors.New("Time duration must be INTEGER")
		}
	}

	howLongInt, _ := strconv.Atoi(howLong)

	return input, howLongInt, nil
}
