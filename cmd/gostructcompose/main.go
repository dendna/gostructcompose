package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dendna/gostructcompose/pkg/common"
	"github.com/dendna/gostructcompose/pkg/data"
	"github.com/dendna/gostructcompose/pkg/gen"
	_ "github.com/lib/pq"
)

// -----------------------------

func main() {

	file := flag.String("file", "meta.json", "")
	flag.Parse()

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(2)
	}

	var meta data.MetaData
	err = json.Unmarshal(content, &meta)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", *file, err)
		os.Exit(2)
	}

	// TODO: put DataSource in JSON file
	tables, err := data.GetTableDesc(data.PostgreSQL, meta)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(2)
	}

	entities, err := common.TransformTables(tables)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(2)
	}

	_, err = gen.Generate(gen.Golang, entities, strings.TrimSuffix(*file, filepath.Ext(*file))+".go")
	if err != nil {
		fmt.Printf("Generating error: %v", err)
		os.Exit(2)
	}
}
