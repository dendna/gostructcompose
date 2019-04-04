package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

// ConnInfo ...
type ConnInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

// DBTable ...
type DBTable struct {
	Schema string
	Table  string
}

// FileData ...
type FileData struct {
	Version string
	Conn    ConnInfo
	Tables  []DBTable
}

// -----------------------------

// DBColumn contains column metadata
type DBColumn struct {
	Name       string
	Datatype   string
	Isnullable string
}

// ------------------------------

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

// -------------------------------

func main() {

	file := flag.String("file", "meta.json", "")
	flag.Parse()
	// fmt.Printf("file: %s\n", *file)

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal("File: ", err)
	}
	// fmt.Printf(*file+" file content: %s", content)

	var meta FileData
	err = json.Unmarshal(content, &meta)
	if err != nil {
		log.Fatal("JSON: ", err)
	}
	// fmt.Printf("%+v", meta)

	// conninfo := "host=192.168.0.73 user=postgreadmin1 password=postgreadmin dbname=moscow sslmode=disable"
	// conninfo := "postgres://postgreadmin:postgreadmin@192.168.0.73/moscow?sslmode=verify-full"
	conninfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		meta.Conn.Host, meta.Conn.Port, meta.Conn.User, meta.Conn.Password, meta.Conn.Dbname)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatal("DB connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("DB ping: ", err)
	}

	var tablecols *[]DBColumn
	var entities []Entity
	for _, value := range meta.Tables {
		// fmt.Println(value)
		tablecols = getTableColumns(db, &value)
		// fmt.Println(tablecols)

		entity := getEntity(value.Table, tablecols)
		entities = append(entities, *entity)
		// fmt.Println(tablestruct)
	}
	// fmt.Println(tablestructs)
	generateGoFile(strings.TrimSuffix(*file, filepath.Ext(*file))+".go", &entities)
}

func getTableColumns(db *sql.DB, table *DBTable) *[]DBColumn {
	rows, err := db.Query("select column_name, data_type, is_nullable from information_schema.columns where table_schema = $1 and table_name = $2",
		table.Schema, table.Table)
	if err != nil {
		log.Fatal("DB query: ", err)
	}
	defer rows.Close()

	var tablecols []DBColumn
	for rows.Next() {
		var tablecol DBColumn
		if err = rows.Scan(&tablecol.Name, &tablecol.Datatype, &tablecol.Isnullable); err != nil {
			log.Fatal("Rows: ", err)
		}
		tablecols = append(tablecols, tablecol)
		//fmt.Println(tablecol.Name)
	}
	return &tablecols
}

func getEntity(tablename string, cols *[]DBColumn) *Entity {

	var attrs = make([]Attribute, len(*cols))
	for index, value := range *cols {
		attrs[index].Name = value.Name
		attrs[index].Type = getGoTypeByPg(value.Datatype, value.Isnullable)
	}

	var entity = Entity{tablename, attrs}
	return &entity
}

func generateGoFile(filename string, entities *[]Entity) {
	t := template.Must(template.New("go.tmpl").ParseFiles("go.tmpl"))

	var buf bytes.Buffer
	err := t.Execute(&buf, entities)
	if err != nil {
		log.Fatal("Generate: ", err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Creating file: ", err)
	}
	file.Write(out)
	defer file.Close()

}

func getGoTypeByPg(pgType string, nullable string) (ret string) {
	switch pgType {
	case "integer", "bigint", "smallint":
		ret = "int"
	case "double precision", "numeric", "real":
		ret = "float64"
	case "text", "character", "character varying":
		if nullable == "YES" {
			ret = "string"
		} else {
			ret = "*string"
		}
	case "date", "time", "timestamp":
		ret = "time.Time"
	case "boolean":
		ret = "bool"
	case "USER-DEFINED":
		ret = "[]byte"
	default:
		return "! Unknown"
	}

	return ret
}
