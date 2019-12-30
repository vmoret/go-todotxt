// Command todotxt provides a cli to work with todo.text files.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pborman/getopt/v2"
	"github.com/vmoret/todotxt/pkg/todotxt"
)

const (
	// TodoDirEnv holds the environment key that holds the TODO directory.
	TodoDirEnv = "TODO_DIR"
	// DefaultTodoFile holds the default todo file
	DefaultTodoFile = "todo.txt"
)

var (
	force  = getopt.BoolLong("force", 'F', "Forces actions without confirmation or interactive input.")
	help   = getopt.Bool('h', "Display a short help message.")
	plain  = getopt.Bool('p', "Plain mode turns off colors.")
	notime = getopt.Bool('t', "Do not prepend the current date to a task when it's added.")
	file   = getopt.StringLong("file", 'f', DefaultTodoFile, "Add task to this file instead to default file.")
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

	switch action := strings.ToLower(args[0]); action {
	case "add":
		switch argc {
		}
	case "list":
		f, err := os.Open(path)
		switch {
		case os.IsNotExist(err):
			return
		case err != nil:
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		defer f.Close()
		tasks := make(todotxt.Tasks, 0)
		if err := tasks.Decode(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode tasks, %v", err)
			return
		}
		for i, t := range tasks {
			fmt.Fprintf(os.Stdout, "%d %s\n", i, t)
		}
	}
}
