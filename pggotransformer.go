package gostructcompose

import "errors"

// PgGoConverter ...
type PgGoConverter struct{}

// Convert implements TypeConverter interface to convert PostgreSQL types into Golang ones
func (pgc PgGoConverter) Convert(srcType string, nullable string) (dstType string, err error) {
	switch srcType {
	case "integer", "bigint", "smallint":
		dstType = "int"
	case "double precision", "numeric", "real":
		dstType = "float64"
	case "text", "character", "character varying":
		dstType = "string"
	case "date", "time", "timestamp":
		dstType = "time.Time"
	case "boolean":
		dstType = "bool"
	case "USER-DEFINED":
		dstType = "[]byte"
	default:
		return "", errors.New("unknown PostgreSQL datatype")
	}

	// TODO: decide about array types
	if nullable == "YES" {
		dstType = "*" + dstType
	}

	return dstType, nil
}
