package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Task struct {
	Description string
	DueDate     string
	Priority    string
	Completed   bool
}

type TaskList struct {
	Tasks []Task
}

func addTask() Task {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter task description: ")
	scanner.Scan()
	description := scanner.Text()

	fmt.Print("Enter time frame(today, tomorrow, next week, year 3000?): ")
	scanner.Scan()
	dueDate := scanner.Text()

	fmt.Print("Enter priority (high/medium/low): ")
	scanner.Scan()
	priority := scanner.Text()

	return Task{
		Description: description,
		DueDate:     dueDate,
		Priority:    priority,
	}
}

func editTask(tasks []Task) []Task {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the task number to edit: ")
	scanner.Scan()
	taskNumber, _ := strconv.Atoi(scanner.Text())
	if taskNumber < 1 || taskNumber > len(tasks) {
		fmt.Println("Invalid task number, please try again")
		return tasks
	}
	taskIndex := taskNumber - 1

	tasks[taskIndex] = addTask()
	return tasks
}

func deleteTask(tasks []Task) []Task {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the task number to delete: ")
	scanner.Scan()
	taskNumber, _ := strconv.Atoi(scanner.Text())
	if taskNumber < 1 || taskNumber > len(tasks) {
		fmt.Println("Invalid task number, please try again")
		return tasks
	}

	return append(tasks[:taskNumber-1], tasks[taskNumber:]...)
}

func printTasks(tasks []Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Description", "Due Date", "Priority"})

	descWidth := 200
	dueWidth := 40
	priWidth := 40

	// Set the column widths
	table.SetColWidth(descWidth)
	table.SetColWidth(dueWidth)
	table.SetColWidth(priWidth)

	for i, task := range tasks {
		description := task.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}

		var priorityColor string
		switch task.Priority {
		case "low":
			priorityColor = "\x1b[32m" // green
		case "medium":
			priorityColor = "\x1b[33m" // yellow
		case "high":
			priorityColor = "\x1b[31m" // red
		}

		table.Append([]string{
			strconv.Itoa(i + 1),
			description,
			task.DueDate,
			priorityColor + task.Priority + "\x1b[0m", // reset color
		})
	}

	table.Render()
}

func saveTasks(tasks []Task) error {
	taskList := TaskList{Tasks: tasks}
	data, err := json.MarshalIndent(taskList, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("tasks.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadTasks() ([]Task, error) {
	data, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		return nil, err
	}

	var taskList TaskList
	err = json.Unmarshal(data, &taskList)
	if err != nil {
		return nil, err
	}

	return taskList.Tasks, nil
}

func main() {
	var tasks []Task

	// Load the task list from the JSON file
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
	}

	// Display the task list
	printTasks(tasks)

	scanner := bufio.NewScanner(os.Stdin)
mainLoop:
	for {
		fmt.Print("Enter a command ([a]dd, [e]dit, [d]elete, [q]uit): ")
		scanner.Scan()
		command := scanner.Text()

		switch strings.ToLower(command) {
		case "a":
			task := addTask()
			tasks = append(tasks, task)
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("Error saving tasks:", err)
			}
			printTasks(tasks)
		case "e":
			tasks = editTask(tasks)
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("Error saving tasks:", err)
			}
			printTasks(tasks)
		case "d":
			tasks = deleteTask(tasks)
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("Error saving tasks:", err)
			}
			printTasks(tasks)
		case "q":
			// Save the task list to the JSON file
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("Error saving tasks:", err)
			}
			printTasks(tasks)
			break mainLoop
		default:
			fmt.Println("Invalid command, please try again")
		}
	}
}
