package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
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

func main() {
	r := gin.Default()

	r.GET("/tasks", func(c *gin.Context) {
		data, err := ioutil.ReadFile("tasks.json")
		if err != nil {
			// If the file doesn't exist, return an empty list
			if os.IsNotExist(err) {
				c.JSON(http.StatusOK, []Task{})
				return
			}

			// For other errors, return a server error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading tasks"})
			return
		}

		var tasks TaskList
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
			return
		}

		// Return the list of tasks
		c.JSON(http.StatusOK, tasks.Tasks)
	})

	r.POST("/tasks", func(c *gin.Context) {
		// Define a new task
		var newTask Task

		// Bind the JSON body of the request to the newTask variable
		if err := c.BindJSON(&newTask); err != nil {
			// If the JSON is not properly formatted, return a bad request status
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load the existing tasks
		data, err := ioutil.ReadFile("tasks.json")
		if err != nil {
			if !os.IsNotExist(err) {
				// If there's an error (other than the file not existing), return a server error response
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading tasks"})
				return
			}
		}

		var tasks TaskList
		if len(data) > 0 {
			// If there are existing tasks, unmarshal the data
			err = json.Unmarshal(data, &tasks)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
				return
			}
		}

		// Append the new task to the task list
		tasks.Tasks = append(tasks.Tasks, newTask)

		// Save the updated task list
		data, err = json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		err = ioutil.WriteFile("tasks.json", data, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		// Return the new task
		c.JSON(http.StatusOK, newTask)
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		// Get the ID of the task to update
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		// Load the existing tasks
		data, err := ioutil.ReadFile("tasks.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading tasks"})
			return
		}

		var tasks TaskList
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
			return
		}

		// Check if the task exists
		if id > len(tasks.Tasks) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// Define a new task
		var updatedTask Task

		// Bind the JSON body of the request to the updatedTask variable
		if err := c.BindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the task
		tasks.Tasks[id-1] = updatedTask

		// Save the updated task list
		data, err = json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		err = ioutil.WriteFile("tasks.json", data, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		// Return the updated task
		c.JSON(http.StatusOK, updatedTask)
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		// Get the ID of the task to delete
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		// Load the existing tasks
		data, err := ioutil.ReadFile("tasks.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading tasks"})
			return
		}

		var tasks TaskList
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tasks"})
			return
		}

		// Check if the task exists
		if id > len(tasks.Tasks) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// Delete the task
		tasks.Tasks = append(tasks.Tasks[:id-1], tasks.Tasks[id:]...)

		// Save the updated task list
		data, err = json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		err = ioutil.WriteFile("tasks.json", data, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving tasks"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
