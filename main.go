package main

import (
	"fmt"
	"strings"
	"slices"
	"strconv"
	"bufio"
	"errors"
	"io"
	"os"
)


type Status string

const (
	StatusTodo 	Status = "todo"
	StatusDoing Status = "doing"
	StatusDone 	Status = "done"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      Status
}

func (t Task) IsDone() bool {
	return t.Status == StatusDone
}

func (t *Task) MarkDone() {
	t.Status = StatusDone
}

type TaskList []Task

func (tasks TaskList) getById(id int) (*Task, error){
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
			*tasks = slices.Delete(*tasks, i, i + 1)
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

	id, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var tasks TaskList 
	nextID := 1	

	for {
		fmt.Print("> ")
		request, err := readInput(scanner)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input error:", err)
			return
		}
	
		request = strings.TrimSpace(request)

		switch request {
		case "help":
			fmt.Println("Commands are: help, add, list, exit")
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
		case "done":
			taskID, err := readTaskId(scanner)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}

			t, err := tasks.getById(taskID)
			if err != nil{
				fmt.Fprintln(os.Stderr, "ERROR: ", err)
				continue
			}

			t.MarkDone()
			
			fmt.Println("Task with ID:", taskID, " is now marked as (done) ")
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
			fmt.Printf("Task created with ID: %d\n", nextID); nextID += 1
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
