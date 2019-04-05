package data

import (
	"database/sql"
	"errors"
	"fmt"
)

// MetaConnInfo ...
type MetaConnInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

// MetaTable ...
type MetaTable struct {
	Schema string
	Table  string
}

// MetaData ...
type MetaData struct {
	Version string
	Conn    MetaConnInfo
	Tables  []MetaTable
}

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

// DataSource ...
type DataSource int

func (ds DataSource) String() string {
	return [...]string{"PostgreSQL", "MySQL", "File"}[ds]
}

const (
	// PostgreSQL ...
	PostgreSQL DataSource = iota
	// MySQL ...
	MySQL
	// File ...
	File
)

// GetTableDesc ...
func GetTableDesc(ds DataSource, meta MetaData) (ret []Table, err error) {
	switch ds {
	case PostgreSQL:
		return getPostgreTables(meta.Conn, meta.Tables)
	default:
		return nil, errors.New("data source " + ds.String() + " is not implemented")
	}
}

func getPostgreTables(conninfo MetaConnInfo, tables []MetaTable) (ret []Table, err error) {
	connstr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conninfo.Host, conninfo.Port, conninfo.User, conninfo.Password, conninfo.Dbname)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	var cols []Column
	ret = make([]Table, len(tables))
	for index, value := range tables {
		cols, err = getPostgreTableColumns(db, &value)
		if err != nil {
			return nil, err
		}
		ret[index].Name = value.Table
		ret[index].Columns = cols
	}

	return ret, nil
}

func getPostgreTableColumns(db *sql.DB, table *MetaTable) (ret []Column, err error) {
	rows, err := db.Query("select column_name, data_type, is_nullable from information_schema.columns where table_schema = $1 and table_name = $2",
		table.Schema, table.Table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var col Column
		if err = rows.Scan(&col.Name, &col.DataType, &col.IsNullable); err != nil {
			return nil, err
		}
		ret = append(ret, col)
	}
	return ret, nil
}
