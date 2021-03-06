// Command todotxt provides a cli to work with todo.text files.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
	"github.com/vmoret/todotxt/pkg/todotxt"
	"github.com/vmoret/todotxt/pkg/todotxt/priority"
)

const (
	// TodoDirEnv holds the environment key that holds the TODO directory.
	TodoDirEnv = "TODO_DIR"
	// DefaultTodoFile holds the default todo file
	DefaultTodoFile = "todo.txt"
)

var (
	force = getopt.BoolLong("force", 'F', "Forces actions without confirmation or interactive input.")
	help  = getopt.Bool('h', "Display a short help message.")
	plain = getopt.Bool('p', "Plain mode turns off colors.")
	file  = getopt.StringLong("file", 'f', DefaultTodoFile, "Add task to this file instead to default file.")
)

func main() {
	// Parse the program arguments
	getopt.Parse()
	if *help {
		getopt.Usage()
		return
	}

	todoDir := os.Getenv(TodoDirEnv)
	path := filepath.Join(todoDir, *file)

	// Get the remaining positional arguments
	args := getopt.Args()

	if len(args) == 0 {
		args = []string{"list"}
	}
	argc := len(args)

	tasks, err := todotxt.ReadFile(path)
	handleError(err)

	switch action := strings.ToLower(args[0]); action {
	case "add":
		task, err := todotxt.NewTask(strings.Join(args[1:], " "))
		handleError(err)
		tasks.Add(task)
		err = todotxt.WriteFile(path, tasks)
		handleError(err)
		tasks.Fprint(os.Stdout)

	case "append":
		if argc < 3 {
			fmt.Println("Missing task number and/or task description")
			os.Exit(1)
		}
		i, err := getTaskNr(args, 1)
		handleError(err)
		s := strings.Join(args[2:], " ")
		err = tasks[i].SetDescription(tasks[i].Description() + " " + s)
		handleError(err)
		err = todotxt.WriteFile(path, tasks)
		handleError(err)

	case "do":
		if argc == 1 {
			fmt.Println("Missing task number")
			os.Exit(1)
		}
		i, err := getTaskNr(args, 1)
		handleError(err)
		tasks[i].MarkCompleted()
		err = todotxt.WriteFile(path, tasks)
		handleError(err)

	case "pri":
		if argc < 3 {
			fmt.Println("Missing task number and/or priority")
			os.Exit(1)
		}
		i, err := getTaskNr(args, 1)
		handleError(err)
		tasks[i].SetPriority(priority.Priority(args[2][0]))
		err = todotxt.WriteFile(path, tasks)
		handleError(err)

	case "depri":
		if argc == 1 {
			fmt.Println("Missing task number")
			os.Exit(1)
		}
		i, err := getTaskNr(args, 1)
		handleError(err)
		tasks[i].SetPriority(priority.ZeroPriority)
		err = todotxt.WriteFile(path, tasks)
		handleError(err)

	case "sort":
		sort.Sort(todotxt.ByString(tasks))
		err = todotxt.WriteFile(path, tasks)
		handleError(err)

	case "list":
		tasks.Fprint(os.Stdout)
	}
}

func getTaskNr(args []string, i int) (int, error) {
	i, err := strconv.Atoi(args[i])
	if err != nil {
		return -1, err
	}
	return i + 1, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
