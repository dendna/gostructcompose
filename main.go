package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

// Conninfo ...
type Conninfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

// Table ...
type Table struct {
	Schema string
	Table  string
}

// Filedata ...
type Filedata struct {
	Version string
	Conn    Conninfo
	Tables  []Table
}

// Tablecol ...
type Tablecol struct {
	Name       string
	Datatype   string
	Isnullable string
}

func main() {

	file := flag.String("file", "meta65.json", "")
	flag.Parse()
	// fmt.Printf("file: %s\n", *file)

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal("File: ", err)
	}
	// fmt.Printf(*file+" file content: %s", content)

	var meta Filedata
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

	var tablecols *[]Tablecol
	for _, value := range meta.Tables {
		// fmt.Println(value)
		tablecols = getTableStruct(db, &value)
		fmt.Println(tablecols)
	}
}

func getTableStruct(db *sql.DB, table *Table) *[]Tablecol {
	rows, err := db.Query("select column_name, data_type, is_nullable from information_schema.columns where table_schema = $1 and table_name = $2", table.Schema, table.Table)
	if err != nil {
		log.Fatal("DB query: ", err)
	}
	defer rows.Close()

	var tablecols []Tablecol
	for rows.Next() {
		var tablecol Tablecol
		if err = rows.Scan(&tablecol.Name, &tablecol.Datatype, &tablecol.Isnullable); err != nil {
			log.Fatal("Rows: ", err)
		}
		tablecols = append(tablecols, tablecol)
		//fmt.Println(tablecol.Name)
	}
	return &tablecols
}
