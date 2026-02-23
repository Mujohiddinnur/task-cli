package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const fileName = "tasks.json"

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	json.Unmarshal(data, &tasks)
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	return os.WriteFile(fileName, data, 0644)
}

func addTask(desc string) {
	tasks, _ := loadTasks()

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	now := time.Now().Format(time.RFC3339)

	task := Task{
		ID:          id,
		Description: desc,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, task)
	saveTasks(tasks)

	fmt.Printf("Task added successfully (ID: %d)\n", id)
}

func listTasks(filter string) {
	tasks, _ := loadTasks()

	for _, t := range tasks {
		if filter == "" || t.Status == filter {
			fmt.Printf("[%d] %s (%s)\n", t.ID, t.Description, t.Status)
		}
	}
}

func markTask(id int, status string) {
	tasks, _ := loadTasks()

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			saveTasks(tasks)
			fmt.Println("Task updated.")
			return
		}
	}
	fmt.Println("Task not found.")
}

func deleteTask(id int) {
	tasks, _ := loadTasks()
	newTasks := []Task{}

	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}

	saveTasks(newTasks)
	fmt.Println("Task deleted.")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command>")
		return
	}

	cmd := os.Args[1]

	switch cmd {

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Provide task description.")
			return
		}
		addTask(os.Args[2])

	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		listTasks(filter)

	case "mark-done":
		id, _ := strconv.Atoi(os.Args[2])
		markTask(id, "done")

	case "mark-in-progress":
		id, _ := strconv.Atoi(os.Args[2])
		markTask(id, "in-progress")

	case "delete":
		id, _ := strconv.Atoi(os.Args[2])
		deleteTask(id)

	default:
		fmt.Println("Unknown command")
	}
}