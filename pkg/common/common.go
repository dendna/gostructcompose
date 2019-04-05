package common

import (
	"errors"

	"github.com/dendna/gostructcompose/pkg/data"

	"github.com/dendna/gostructcompose/pkg/gen"
)

// TransformTables transforms data.Table into gen.Entity structure
func TransformTables(tables []data.Table) (ret []gen.Entity, err error) {
	ret = make([]gen.Entity, len(tables))
	for index, value := range tables {
		entity, err := getEntity(value)
		if err != nil {
			return nil, err
		}
		ret[index] = *entity
	}
	return ret, nil
}

func getEntity(table data.Table) (*gen.Entity, error) {
	var attrs = make([]gen.Attribute, len(table.Columns))
	var err error
	for index, value := range table.Columns {
		attrs[index].Name = value.Name
		attrs[index].Type, err = getGoTypeByPg(value.DataType, value.IsNullable)
		if err != nil {
			return nil, err
		}
	}

	return &gen.Entity{Name: table.Name, Attrs: attrs}, nil
}

func getGoTypeByPg(pgType string, nullable string) (ret string, err error) {
	switch pgType {
	case "integer", "bigint", "smallint":
		ret = "int"
	case "double precision", "numeric", "real":
		ret = "float64"
	case "text", "character", "character varying":
		ret = "string"
	case "date", "time", "timestamp":
		ret = "time.Time"
	case "boolean":
		ret = "bool"
	case "USER-DEFINED":
		ret = "[]byte"
	default:
		return "", errors.New("unknown PostgreSQL datatype")
	}

	// TODO: decide about slices
	if nullable == "YES" {
		ret = "*" + ret
	}

	return ret, nil
}
