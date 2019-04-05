package gen

import (
	"bytes"
	"errors"
	"go/format"
	"html/template"
	"os"
)

// Attribute ...
type Attribute struct {
	Name string
	Type string
}

// Entity ...
type Entity struct {
	Name  string
	Attrs []Attribute
}

// Language ...
type Language int

func (l Language) String() string {
	return [...]string{"Golang", "Java", "JavaScript"}[l]
}

const (
	// Golang ...
	Golang Language = iota
	// Java ...
	Java
	// JavaScript ...
	JavaScript
)

// Generate ...
func Generate(lang Language, entities []Entity, outfile string) (n int, err error) {
	switch lang {
	case Golang:
		return generateGoFile(entities, outfile)
	default:
		return 0, errors.New("code generator for " + lang.String() + " is not implemented")
	}
}

func generateGoFile(entities []Entity, outfile string) (n int, err error) {

	t, err := template.New("golang.tmpl").ParseFiles("golang.tmpl")
	if err != nil {
		return 0, err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, entities)
	if err != nil {
		return 0, err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, err
	}

	file, err := os.Create(outfile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write(out)
}
