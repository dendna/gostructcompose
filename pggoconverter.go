package gostructcompose

import (
	"strings"
)

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
		// TODO: implement remaining data types
		dstType = "STUB"
		// return "", fmt.Errorf("unknown PostgreSQL datatype: %v", srcType)
	}

	// TODO: decide about array types
	if strings.ToUpper(nullable) == "YES" {
		dstType = "*" + dstType
	}

	return dstType, nil
}
