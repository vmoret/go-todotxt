package todotxt

import (
	"fmt"
	"os"
)

// ReadFile reads tasks from file on given path.
func ReadFile(path string) (tasks Tasks, err error) {
	f, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		return
	case err != nil:
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	defer f.Close()
	err = tasks.Decode(f)
	return
}

// WriteFile writes tasks to file on given path.
func WriteFile(path string, tasks Tasks) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	err = tasks.Encode(f)
	return
}
