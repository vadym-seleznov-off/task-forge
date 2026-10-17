package main

import (
	"fmt"
	"strings"
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

func printTasks(tasks []Task) {
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	tasks := []Task{}
	id := 1	

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
			newTask.ID = id
			newTask.Description = description
			newTask.Status = StatusTodo

			tasks = append(tasks, newTask)
			fmt.Printf("Task created with ID: %d\n", id); id += 1
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
