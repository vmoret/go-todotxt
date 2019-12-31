package priority

import (
	"errors"
	"fmt"

	"gopkg.in/gookit/color.v1"
)

var (
	// ErrInvalidPriority is returned when an invalid priority is submitted
	ErrInvalidPriority = errors.New("invalid priority")
)

// Priority represents a task priority
type Priority rune

// ZeroPriority indicates no priority
var ZeroPriority Priority = Priority(91)

// Validate validates the priority rune
func (p Priority) Validate() error {
	if p == ZeroPriority {
		return nil
	}
	if p < 65 || p > 90 {
		return ErrInvalidPriority
	}
	return nil
}

// IsZero indicates this priority is zero.
func (p Priority) IsZero() bool { return p == ZeroPriority }

// Bytes returns a bytes representation of this priority
func (p Priority) Bytes() []byte {
	return []byte(fmt.Sprintf("%c", p))
}

func (p Priority) String() string {
	if p == ZeroPriority {
		return ""
	}
	return fmt.Sprintf("(%c) ", p)
}

// Render renders given a using the proper color renderer
func (p Priority) Render(a ...interface{}) string {
	return renderers[p](a...)
}

var (
	renderers map[Priority]func(a ...interface{}) string
)

func init() {
	renderers = make(map[Priority]func(a ...interface{}) string, 0)
	renderers[ZeroPriority] = color.FgBlack.Render
	renderers[Priority('A')] = color.FgYellow.Render
	renderers[Priority('B')] = color.FgGreen.Render
	renderers[Priority('C')] = color.FgLightBlue.Render
	for i := 68; i <= 90; i++ {
		renderers[Priority(i)] = color.FgBlack.Render
	}
}
