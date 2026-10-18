package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func (t Task) IsDone() bool {
	return t.Status == StatusDone
}

func (t *Task) MarkDone() {
	t.Status = StatusDone
}

type TaskList []Task

func (tasks TaskList) getById(id int) (*Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}

	return nil, errors.New("task not found")
}

func (tasks *TaskList) deleteById(id int) error {
	for i := range *tasks {
		if (*tasks)[i].ID == id {
			*tasks = slices.Delete(*tasks, i, i+1)
			return nil
		}
	}

	return errors.New("task not found")
}

func printTasks(tasks TaskList) {
	if len(tasks) == 0 {
		fmt.Println("Task list is empty at this moment")
		return
	}

	for _, task := range tasks {
		fmt.Printf("[%d] %s (%s)\n", task.ID, task.Title, task.Status)
	}
}

func readInput(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	return scanner.Text(), nil
}

func readTaskId(scanner *bufio.Scanner) (int, error) {
	fmt.Print("Provide ID of task: ")
	input, err := readInput(scanner)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, err
		}

		return 0, err
	}

	input = strings.TrimSpace(input)

	id, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func writeTasks(tasks TaskList) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var tasks TaskList

	data, err := os.ReadFile("tasks.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Could not load tasks from a file tasks.json")
			tasks = TaskList{}
		} else {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return
		}
	} else {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while tried to parse json:", err)
			return
		}

		fmt.Println("Successfully load tasks from tasks.json file")
	}

	scanner := bufio.NewScanner(os.Stdin)
	nextID := 1

	for {
		fmt.Print("> ")
		request, err := readInput(scanner)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nGoodbye!")
				return
			}

			fmt.Fprintln(os.Stderr, "reading standard input error:", err)
			return
		}

		request = strings.TrimSpace(request)

		switch request {
		case "help":
			fmt.Println("Commands are: help, add, list, exit, done, delete")
		case "delete":
			taskID, err := readTaskId(scanner)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}

			if err = tasks.deleteById(taskID); err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}

			fmt.Println("Successfully deleted task with ID:", taskID)
			err = writeTasks(tasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Something went wrong when tried to delete task:", err)
			}
		case "done":
			taskID, err := readTaskId(scanner)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}

			t, err := tasks.getById(taskID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: ", err)
				continue
			}

			t.MarkDone()

			fmt.Println("Task with ID:", taskID, " is now marked as (done) ")

			err = writeTasks(tasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Something went wrong when tried to update task:", err)
			}
		case "add":
			var newTask Task

			fmt.Print("Title: ")
			title, err := readInput(scanner)

			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("\nGoodbye!")
					return
				}

				fmt.Fprintln(os.Stderr, "input error:", err)
				return
			}

			title = strings.TrimSpace(title)

			if title == "" {
				fmt.Fprintln(os.Stderr, "Title of task cannot be empty!")
				break
			}

			fmt.Print("Description: ")
			description, err := readInput(scanner)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("\nGoodbye!")
					return
				}

				fmt.Fprintln(os.Stderr, "input error:", err)
				return
			}

			description = strings.TrimSpace(description)

			newTask.Title = title
			newTask.ID = nextID
			newTask.Description = description
			newTask.Status = StatusTodo

			tasks = append(tasks, newTask)
			fmt.Printf("Task created with ID: %d\n", nextID)
			nextID += 1

			err = writeTasks(tasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Something went wrong when tried to add a new task:", err)
			}
		case "list":
			printTasks(tasks)
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Fprintln(os.Stderr, "Unkown request:", request)
		}
	}

}
