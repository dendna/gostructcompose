package gostructcompose

import (
	"database/sql"
	"strings"

	// postgresql vendor driver
	_ "github.com/lib/pq"
)

// PostgreReader ...
type PostgreReader struct {
	db *sql.DB
}

// Read implements MetaReader interface for using with PostgreSQL database
func (pr PostgreReader) Read(location string, items []Item) (ret []Table, err error) {

	err = pr.connect(location)
	if err != nil {
		return nil, err
	}
	defer pr.db.Close()

	return pr.getTables(items)
}

func (pr *PostgreReader) connect(connstr string) (err error) {
	pr.db, err = sql.Open("postgres", connstr)
	if err != nil {
		return err
	}

	err = pr.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgreReader) getTables(items []Item) (ret []Table, err error) {
	var cols []Column
	ret = make([]Table, len(items))
	for index, value := range items {
		// TODO: error handling strings.Split ?
		cols, err = pr.getColumns(strings.Split(value.FullName, "."))
		if err != nil {
			return nil, err
		}
		ret[index].Name = strings.Split(value.FullName, ".")[1]
		ret[index].Columns = cols
	}

	// fmt.Println(ret)

	return ret, nil
}

func (pr *PostgreReader) getColumns(fullname []string) (ret []Column, err error) {
	rows, err := pr.db.Query("select column_name, data_type, is_nullable from information_schema.columns where table_schema = $1 and table_name = $2",
		fullname[0], fullname[1])
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
