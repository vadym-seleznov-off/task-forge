package main

import (
	"fmt"
	"bufio"
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
	for _, task := range tasks {
		fmt.Printf("[%d] %s (%s)\n", task.ID, task.Title, task.Status)
	}
}

func main() {
	tasks := []Task{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input error:", err)
		}

		request := scanner.Text()
		
		switch request {
		case "help":
			fmt.Println("Commands are: help, add, list, exit")
		case "add":
			var newTask Task
			
			fmt.Print("Title: ")
			scanner.Scan()
			
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input error:", err)
			}

			title := scanner.Text()

			fmt.Print("Description: ")
			scanner.Scan()
			
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input error:", err)
			}

			description := scanner.Text()

			id := len(tasks) + 1

			newTask.Title = title
			newTask.ID = id
			newTask.Description = description
			newTask.Status = StatusTodo

			tasks = append(tasks, newTask)
			fmt.Printf("Task created with ID: %d\n", id)
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
