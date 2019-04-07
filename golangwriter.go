package gostructcompose

import (
	"bytes"
	"go/format"
	"html/template"
	"os"
)

// GolangWriter implements EntityWriter interface to generate Golang source code
type GolangWriter struct{}

// Write implement EntityWriter interface
func (gw GolangWriter) Write(dest string, entities []Entity) error {
	t, err := template.New("golang.tmpl").ParseFiles("golang.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, entities)
	if err != nil {
		return err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(out)

	return err
}
