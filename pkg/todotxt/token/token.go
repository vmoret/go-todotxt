package token

import (
	"bytes"
)

// Type represents a token type
type Type int

// Supported tokens
const (
	Word Type = iota
	ProjectTag
	ContextTag
	KeyValueTag
)

// Token represents a token.
type Token struct {
	Key, Value []byte
	Type       Type
}

var (
	projectTag  = []byte("+")
	contextTag  = []byte("@")
	keyValueTag = []byte(":")
	separator   = []byte(" ")
)

// UnmarshalText implements TextUnmarshaler
func (t *Token) UnmarshalText(text []byte) error {
	switch {
	case bytes.HasPrefix(text, projectTag):
		t.Value = bytes.TrimPrefix(text, projectTag)
		t.Type = ProjectTag
	case bytes.HasPrefix(text, contextTag):
		t.Value = bytes.TrimPrefix(text, contextTag)
		t.Type = ContextTag
	case bytes.Count(text, keyValueTag) == 1:
		a := bytes.Split(text, keyValueTag)
		t.Key = a[0]
		t.Value = a[1]
		t.Type = KeyValueTag
	default:
		t.Value = text
		t.Type = Word
	}
	return nil
}

// MarshalText implements TextMarshaler
func (t *Token) MarshalText() (text []byte, err error) {
	switch t.Type {
	case ProjectTag:
		text = append(projectTag, t.Value...)
	case ContextTag:
		text = append(contextTag, t.Value...)
	case KeyValueTag:
		text = append(keyValueTag, t.Value...)
	default:
		text = t.Value
	}
	return
}

// Tokens represents a collection of tokens.
type Tokens []*Token

// UnmarshalText implements TextUnmarshaler
func (ts *Tokens) UnmarshalText(text []byte) error {
	a := bytes.Split(text, separator)
	temp := make(Tokens, len(a))
	for i, s := range a {
		token := &Token{}
		token.UnmarshalText([]byte(s))
		temp[i] = token
	}
	*ts = temp
	return nil
}

// MarshalText implements TextMarshaler
func (ts Tokens) MarshalText() (text []byte, err error) {
	a := make([][]byte, len(ts))
	for i, t := range ts {
		a[i], _ = t.MarshalText()
	}
	return bytes.Join(a, separator), nil
}
