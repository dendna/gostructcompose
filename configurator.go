package gostructcompose

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	// CurVer is a current version of config file
	CurVer = "1.1"
)

type configJSON struct {
	Version  string
	InType   string
	OutType  string
	ConvType string
	Dest     string
	Location string
	Items    []Item
}

// Configurator ...
type Configurator struct {
	Location string
	Dest     string
	Items    []Item

	Reader    MetaReader
	Writer    EntityWriter
	Converter TypeConverter
}

// LoadFromJSON ...
func (c *Configurator) LoadFromJSON(filename string) error {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var cfg configJSON
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		// fmt.Printf("Error while parsing %v: %v", filename, err)
		return err
	}

	// fmt.Println(cfg)

	if cfg.Version != CurVer {
		return fmt.Errorf("incorrect config version: %v. Expected version %v", cfg.Version, CurVer)
	}

	c.Location = cfg.Location
	c.Dest = cfg.Dest
	c.Items = make([]Item, len(cfg.Items))
	copy(c.Items, cfg.Items)

	err = c.setReader(cfg.InType)
	if err != nil {
		return err
	}

	err = c.setWriter(cfg.OutType)
	if err != nil {
		return err
	}

	err = c.setConverter(cfg.ConvType)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurator) setReader(intype string) error {
	switch intype {
	case "postgresql":
		c.Reader = new(PostgreReader)
	default:
		return fmt.Errorf("incorrect input type %v", intype)
	}

	return nil
}

func (c *Configurator) setWriter(outtype string) error {
	switch outtype {
	case "golang":
		c.Writer = new(GolangWriter)
	default:
		return fmt.Errorf("incorrect output type %v", outtype)
	}

	return nil
}

func (c *Configurator) setConverter(convtype string) error {
	switch convtype {
	case "postgretogolang":
		c.Converter = new(PgGoConverter)
	default:
		return fmt.Errorf("incorrect converter type %v", convtype)
	}

	return nil
}
