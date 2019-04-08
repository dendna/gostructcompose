package gostructcompose

import (
	"errors"
)

// Column contains column metadata
type Column struct {
	Name       string
	DataType   string
	IsNullable string
}

// Table ...
type Table struct {
	Name    string
	Columns []Column
}

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

// Item is a fully qualified name of the source data
// example: "schema_name.table_name" or "filename.extension"
type Item struct {
	FullName string
}

// MetaReader ...
type MetaReader interface {
	Read(location string, items []Item) (ret []Table, err error)
}

// EntityWriter  ...
type EntityWriter interface {
	Write(dest string, entities []Entity) error
}

// TypeConverter ...
type TypeConverter interface {
	Convert(srcType string, nullable string) (dstType string, err error)
}

// Generator ...
type Generator struct {
	location string
	items    []Item
	dest     string

	reader    MetaReader
	writer    EntityWriter // TODO: we can use several different EntityWriter at once
	converter TypeConverter
}

// NewGenerator is a constructor ...
func NewGenerator(location string, items []Item, dest string, mr MetaReader, ew EntityWriter, tc TypeConverter) (*Generator, error) {
	// func NewGenerator(location string, items []Item, dest string, mr MetaReader, ew EntityWriter, tc TypeConverter) *Generator {
	if location == "" || items == nil || dest == "" || mr == nil || ew == nil || tc == nil {
		return nil, errors.New("generator creating: empty parameter")
	}

	g := &Generator{
		location:  location,
		dest:      dest,
		reader:    mr,
		writer:    ew,
		converter: tc,
	}
	// TODO: make exception handling ?
	g.items = make([]Item, len(items))
	copy(g.items, items)

	return g, nil
}

// Generate is a main method for the entities generating
func (g *Generator) Generate() error {

	// read meta data from a source
	tables, err := g.reader.Read(g.location, g.items)
	if err != nil {
		return err
	}

	// transform source-specific data into the universal structure
	entities, err := g.transformTables(tables)
	if err != nil {
		return err
	}

	// generate output data depending on the writer implementation (file, db, etc.)
	return g.writer.Write(g.dest, entities)

}

// transformTables transforms data from Table into Entity structure
func (g *Generator) transformTables(tables []Table) (ret []Entity, err error) {
	ret = make([]Entity, len(tables))
	for index, value := range tables {
		entity, err := g.getEntity(&value)
		if err != nil {
			return nil, err
		}
		ret[index] = *entity
	}
	return ret, nil
}

func (g *Generator) getEntity(table *Table) (*Entity, error) {
	if table.Name == "" {
		return nil, errors.New("item name cannot be empty")
	}
	var attrs = make([]Attribute, len(table.Columns))
	var err error
	for index, value := range table.Columns {
		attrs[index].Name = value.Name
		attrs[index].Type, err = g.converter.Convert(value.DataType, value.IsNullable)
		if err != nil {
			return nil, err
		}
	}

	return &Entity{Name: table.Name, Attrs: attrs}, nil
}
