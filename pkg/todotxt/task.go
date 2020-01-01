package todotxt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/vmoret/todotxt/pkg/todotxt/date"
	"github.com/vmoret/todotxt/pkg/todotxt/priority"
	"github.com/vmoret/todotxt/pkg/todotxt/token"
	"gopkg.in/gookit/color.v1"
)

var (
	lineFeed = []byte("\n")
	// Regular expresion used to parse priority
	priorityRe = regexp.MustCompile(`^^\(([A-Z])\){1} `)

	// Regular expresion used to parse date
	dateRe = regexp.MustCompile(`^([0-9]{2,4}-[0-9]{2}-[0-9]{2}){1} `)
)

// DecodeTask decodes a Task from given string s.
func DecodeTask(s string) (task *Task, err error) {
	task = new(Task)
	err = task.UnmarshalText([]byte(s))
	return task, err
}

// Task represents a task.
type Task struct {
	completed      bool
	priority       priority.Priority
	completionDate date.Date
	creationDate   date.Date
	description    token.Tokens
}

// NewTask creates a new task from string s.
func NewTask(s string) (*Task, error) {
	task, err := DecodeTask(s)
	task.creationDate = date.Now()
	return task, err
}

// MarkCompleted marks the task as completed.
func (t *Task) MarkCompleted() error {
	t.completed = true
	if !t.creationDate.IsZero() {
		t.completionDate = date.Now()
	}
	if !t.priority.IsZero() {
		t.description = append(t.description, &token.Token{
			Type: token.KeyValueTag, Key: []byte("pri"), Value: t.priority.Bytes()})
		t.priority = priority.ZeroPriority
	}
	return nil
}

// SetPriority sets the priority of the task
func (t *Task) SetPriority(p priority.Priority) error {
	t.priority = p
	return nil
}

// Description returns the task description.
func (t *Task) Description() string {
	b, err := t.description.MarshalText()
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// SetDescription sets the task description.
func (t *Task) SetDescription(s string) error {
	return t.description.UnmarshalText([]byte(s))
}

// Validate a task
func (t *Task) Validate() error {
	if err := t.priority.Validate(); err != nil {
		return err
	}
	if !t.completionDate.IsZero() && !t.completed {
		t.completionDate = date.ZeroDate
	}
	if t.completed && t.completionDate.IsZero() && t.creationDate.IsZero() {
		t.completionDate = date.ZeroDate
	}
	return nil
}

// UnmarshalText implements TextUnmarshaler
func (t *Task) UnmarshalText(text []byte) error {
	if bytes.HasPrefix(text, []byte("x ")) {
		t.completed = true
		text = text[2:]
	}
	matches := priorityRe.FindAllSubmatch(text, -1)
	if len(matches) == 1 {
		t.priority = priority.Priority([]rune(string(matches[0][1]))[0])
		text = text[len(matches[0][0]):]
	} else {
		t.priority = priority.ZeroPriority
	}
	dates := [2]date.Date{date.ZeroDate, date.ZeroDate}
	for i := range dates {
		matches := dateRe.FindAllSubmatch(text, -1)
		if len(matches) == 1 {
			dates[i] = date.Parse(string(matches[0][1]))
			text = text[len(matches[0][0]):]
		}
	}
	if t.completed {
		t.completionDate = dates[0]
		t.creationDate = dates[1]
	} else {
		t.creationDate = dates[0]
	}
	t.description.UnmarshalText(text)
	return nil
}

// MarshalText implements TextMarshaler
func (t *Task) MarshalText() (text []byte, err error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString("")
	if t.completed {
		buf.WriteString("x ")
	}
	buf.WriteString(t.priority.String())
	buf.WriteString(t.completionDate.String())
	buf.WriteString(t.creationDate.String())
	b, _ := t.description.MarshalText()
	buf.Write(b)
	return buf.Bytes(), nil
}

var completedRender = color.FgGray.Render

func (t *Task) String() string {
	b, err := t.MarshalText()
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// Render renders the task as string.
func (t *Task) Render() string {
	s := t.String()
	if t.completed {
		return completedRender(s)
	}
	return t.priority.Render(s)
}

// Tasks represents a collection of tasks.
type Tasks []*Task

// Decode the tasks from given reader
func (tasks *Tasks) Decode(r io.Reader) (err error) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	temp := make(Tasks, 0)
	for s.Scan() {
		t := new(Task)
		if err := t.UnmarshalText(s.Bytes()); err != nil {
			return err
		}
		temp = append(temp, t)
	}
	*tasks = temp
	return
}

// Encode the tasks to given writer
func (tasks Tasks) Encode(w io.Writer) error {
	for _, t := range tasks {
		b, err := t.MarshalText()
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}
		_, err = w.Write(lineFeed)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add adds given task to the collection
func (tasks *Tasks) Add(task *Task) {
	*tasks = append(*tasks, task)
}

// Fprint prints the tasks to given writer
func (tasks Tasks) Fprint(w io.Writer) {
	for i, t := range tasks {
		fmt.Fprintf(w, "%d %s\n", i+1, t.Render())
	}
}

// ByString implements sort.Interface based on String field.
type ByString Tasks

func (a ByString) Len() int           { return len(a) }
func (a ByString) Less(i, j int) bool { return a[i].String() < a[j].String() }
func (a ByString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
