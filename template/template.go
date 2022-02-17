package template

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"text/template"
)

// Execer is an interface that handles templates.
type Execer interface {
	// Exec parses and executes a string template against the specified
	// data object, returning the output string.
	Exec(text string, data interface{}) (string, error)
}

// New returns a new template execer to execute templates.
func New() Execer {
	return &store{
		templates: make(map[string]*template.Template),
	}
}

type store struct {
	templates map[string]*template.Template
}

func (s store) Exec(text string, data interface{}) (string, error) {
	h := hash(text)
	var (
		tpl *template.Template
		ok  bool
	)
	// find or create stored template in map
	tpl, ok = s.templates[h]
	if !ok {
		var err error
		if tpl, err = parse(h, text); err != nil {
			return "", err
		}
		s.templates[h] = tpl
	}
	ts, err := exec(tpl, data)
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

func exec(t *template.Template, data interface{}) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
