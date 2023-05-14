package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const baseURL = "http://localhost:8080"

type Task struct {
	Description string
	DueDate     string
	Priority    string
	Completed   bool
}

func addTask() error {
	task := Task{
		Description: getDescription(),
		DueDate:     getDueDate(),
		Priority:    getPriority(),
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// Send a POST request to the API to add a task
	resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add task: %v", resp.Status)
	}

	return nil
}

func editTask(taskNumber int) error {
	task := Task{
		Description: getDescription(),
		DueDate:     getDueDate(),
		Priority:    getPriority(),
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// Send a PUT request to the API
	req, err := http.NewRequest("PUT", baseURL+"/tasks/"+strconv.Itoa(taskNumber-1), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to edit task: %v", resp.Status)
	}

	return nil
}

func deleteTask(taskNumber int) error {
	// Send a DELETE request to the API
	req, err := http.NewRequest("DELETE", baseURL+"/tasks/"+strconv.Itoa(taskNumber-1), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to delete task: %v", resp.Status)
	}

	return nil
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
	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return err
	}

	// Send a POST request to the API to save the tasks
	resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to save tasks: %v", resp.Status)
	}

	return nil
}

func loadTasks() ([]Task, error) {
	// Send a GET request to the API to load the tasks
	resp, err := http.Get(baseURL + "/tasks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to load tasks: %v", resp.Status)
	}

	var tasks []Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func getDescription() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the task description: ")
	scanner.Scan()
	return scanner.Text()
}

func getDueDate() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the task due date (format YYYY-MM-DD): ")
	scanner.Scan()
	return scanner.Text()
}

func getPriority() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the task priority (low, medium, high): ")
	scanner.Scan()
	return scanner.Text()
}

func main() {
	var tasks []Task
	var err error

	scanner := bufio.NewScanner(os.Stdin)

mainLoop:
	for {
		// Load the task list from the API
		tasks, err = loadTasks()
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			continue
		}

		// Display the task list
		printTasks(tasks)

		fmt.Print("Enter a command ([a]dd, [e]dit, [d]elete, [q]uit): ")
		scanner.Scan()
		command := scanner.Text()

		switch strings.ToLower(command) {
		case "a":
			err = addTask()
			if err != nil {
				fmt.Println("Error adding task:", err)
			}
		case "e":
			fmt.Print("Enter the task number to edit: ")
			scanner.Scan()
			taskNumber, _ := strconv.Atoi(scanner.Text())
			if taskNumber < 1 || taskNumber > len(tasks) {
				fmt.Println("Invalid task number, please try again")
				continue
			}
			err = editTask(taskNumber)
			if err != nil {
				fmt.Println("Error editing task:", err)
			}
		case "d":
			fmt.Print("Enter the task number to delete: ")
			scanner.Scan()
			taskNumber, _ := strconv.Atoi(scanner.Text())
			if taskNumber < 1 || taskNumber > len(tasks) {
				fmt.Println("Invalid task number, please try again")
				continue
			}
			err = deleteTask(taskNumber)
			if err != nil {
				fmt.Println("Error deleting task:", err)
			}
		case "q":
			break mainLoop
		default:
			fmt.Println("Invalid command, please try again")
		}
	}
}
