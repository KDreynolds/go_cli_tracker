package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/olekukonko/tablewriter"
)

type Task struct {
    Description string
    DueDate     time.Time
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

    fmt.Print("Enter due date (YYYY-MM-DD): ")
    scanner.Scan()
    dueDate, _ := time.Parse("2006-01-02", scanner.Text())

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

        if task.Completed {
            description = "\x1b[2m" + description + "\x1b[0m" // dimmed text
        }

        table.Append([]string{
            strconv.Itoa(i + 1),
            description,
            task.DueDate.Format("2006-01-02"),
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
            fmt.Print("Enter a command (add, edit, delete, quit): ")
            scanner.Scan()
            command := scanner.Text()
        
            switch strings.ToLower(command) {
            case "add":
                task := addTask()
                tasks = append(tasks, task)
                err = saveTasks(tasks)
                if err != nil {
                    fmt.Println("Error saving tasks:", err)
                }
                printTasks(tasks)
            case "edit":
                tasks = editTask(tasks)
                err = saveTasks(tasks)
                if err != nil {
                    fmt.Println("Error saving tasks:", err)
                }
                printTasks(tasks)
            case "delete":
                tasks = deleteTask(tasks)
                err = saveTasks(tasks)
                if err != nil {
                    fmt.Println("Error saving tasks:", err)
                }
                printTasks(tasks)
            case "quit":
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