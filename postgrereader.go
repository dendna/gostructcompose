package gostructcompose

import (
	"database/sql"
	"errors"
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
	ret = make([]Table, 0, len(items))

	for _, value := range items {
		schema, table, err := pr.splitFullName(value.FullName, ".")
		if err != nil {
			return nil, err
		}

		if table == "*" {

			if err = func() error {
				rows, err := pr.db.Query("select table_name "+
					" from information_schema.tables "+
					" where table_type = 'BASE TABLE' and table_schema = $1", schema)
				if err != nil {
					return err
				}
				defer rows.Close()

				for rows.Next() {
					if err = rows.Scan(&table); err != nil {
						return err
					}
					cols, err = pr.getColumns(schema, table)
					if err != nil {
						return err
					}
					ret = append(ret, Table{Name: table, Columns: cols})
				}
				return nil
			}(); err != nil {
				return nil, err
			}

		} else {
			cols, err = pr.getColumns(schema, table)
			if err != nil {
				return nil, err
			}
			ret = append(ret, Table{Name: table, Columns: cols})
		}
	}

	// fmt.Println(ret)

	return ret, nil
}

func (pr *PostgreReader) getColumns(schema, table string) (ret []Column, err error) {
	rows, err := pr.db.Query("select column_name, data_type, is_nullable "+
		" from information_schema.columns "+
		" where table_schema = $1 and table_name = $2",
		schema, table)
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

func (pr *PostgreReader) splitFullName(fullname string, sep string) (schema, table string, err error) {
	const errEmpty string = "table name cannot be empty"

	if fullname == "" {
		return "", "", errors.New(errEmpty)
	}

	if !strings.Contains(fullname, sep) {
		return "public", fullname, nil
	}

	schema = strings.Split(fullname, sep)[0]
	table = strings.Split(fullname, sep)[1]

	if table == "" {
		return "", "", errors.New(errEmpty)
	}

	if schema == "" {
		schema = "public"
	}

	return schema, table, nil
}
