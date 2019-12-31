package todotxt

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/vmoret/todotxt/pkg/todotxt/date"
	"github.com/vmoret/todotxt/pkg/todotxt/priority"
	"github.com/vmoret/todotxt/pkg/todotxt/token"
	"gopkg.in/gookit/color.v1"
)

var (
	// Regular expresions used to parse project, context and key/value tags
	projectTagRe  = regexp.MustCompile(` (\+\S*)`)
	contextTagRe  = regexp.MustCompile(` (\@\S*)`)
	keyValueTagRe = regexp.MustCompile(` (\S+[^:]):(\S+)`)

	// Regular expresion used to parse priority
	priorityRe = regexp.MustCompile(`^^\(([A-Z])\){1} `)

	// Regular expresion used to parse date
	dateRe = regexp.MustCompile(`^([0-9]{2,4}-[0-9]{2}-[0-9]{2}){1} `)
)

// Task represents a task.
type Task struct {
	Completed      bool
	Priority       priority.Priority
	CompletionDate date.Date
	CreationDate   date.Date
	Description    token.Tokens
}

// Validate a task
func (t *Task) Validate() error {
	if err := t.Priority.Validate(); err != nil {
		return err
	}
	if t.CompletionDate.IsZero() && !t.Completed {
		t.CompletionDate = date.ZeroDate
	}
	if t.Completed && t.CompletionDate.IsZero() && t.CreationDate.IsZero() {
		t.CompletionDate = date.ZeroDate
	}
	return nil
}

// UnmarshalText implements TextUnmarshaler
func (t *Task) UnmarshalText(text []byte) error {
	if bytes.HasPrefix(text, []byte("x ")) {
		t.Completed = true
		text = text[2:]
	}
	matches := priorityRe.FindAllSubmatch(text, -1)
	if len(matches) == 1 {
		t.Priority = priority.Priority([]rune(string(matches[0][1]))[0])
		text = text[len(matches[0][0]):]
	}
	dates := [2]date.Date{date.ZeroDate, date.ZeroDate}
	for i := range dates {
		matches := dateRe.FindAllSubmatch(text, -1)
		if len(matches) == 1 {
			dates[i] = date.Parse(string(matches[0][1]))
			text = text[len(matches[0][0]):]
		}
	}
	if t.Completed {
		t.CompletionDate = dates[0]
		t.CreationDate = dates[1]
	} else {
		t.CreationDate = dates[0]
	}
	t.Description.UnmarshalText(text)
	return nil
}

// MarshalText implements TextMarshaler
func (t *Task) MarshalText() (text []byte, err error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString("")
	if t.Completed {
		buf.WriteString("x ")
	}
	buf.WriteString(t.Priority.String())
	buf.WriteString(t.CompletionDate.String())
	buf.WriteString(t.CreationDate.String())
	b, _ := t.Description.MarshalText()
	buf.Write(b)
	return buf.Bytes(), nil
}

var completedRender = color.FgGray.Render

func (t *Task) String() string {
	b, err := t.MarshalText()
	if err != nil {
		return err.Error()
	}
	if t.Completed {
		return completedRender(string(b))
	}
	return t.Priority.Render(string(b))
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
	}
	return nil
}
