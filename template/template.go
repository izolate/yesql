package template

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"sync"
	"text/template"
)

var m = &sync.Map{}

// Executer is an interface for template execution.
type Executer interface {
	// Execute parses and executes a string template against the specified
	// data object, returning the output string.
	Execute(template string, data any) (string, error)
}

// New returns a new template execer to execute templates.
func New() Executer {
	return &store{m}
}

type store struct {
	m *sync.Map
}

func (s store) Execute(text string, data any) (string, error) {
	// generate unique hash for template string
	h := hash(text)

	var (
		tpl *template.Template
		err error
	)

	// either find the stored template in the sync map,
	// or store it in the map if it doesn't already exist.
	val, ok := s.m.Load(h)
	if ok {
		tpl, _ = val.(*template.Template)
	} else {
		tpl, err = parse(h, text)
		if err != nil {
			return "", err
		}
		s.m.Store(h, tpl)
	}

	ts, err := execute(tpl, data)
	if err != nil {
		return "", err
	}
	return ts, nil
}

func hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func parse(hash, text string) (*template.Template, error) {
	t, err := template.New(hash).Parse(text)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func execute(t *template.Template, data any) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
