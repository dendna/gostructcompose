package gostructcompose

import "testing"

func TestSplitFullName(t *testing.T) {
	tests := []struct {
		fullname string
		sep      string
		schema   string
		table    string
		wanterr  bool
	}{
		{fullname: "", sep: ".", schema: "", table: "", wanterr: true},
		{fullname: "schema_name.", sep: ".", schema: "", table: "", wanterr: true},
		// {fullname: ".schema_name.", sep: ".", schema: "", table: "", wanterr: true},
		{fullname: "table_name", sep: ".", schema: "public", table: "table_name", wanterr: false},
		{fullname: ".table_name", sep: ".", schema: "public", table: "table_name", wanterr: false},
		{fullname: "schema_name.table_name", sep: ".", schema: "schema_name", table: "table_name", wanterr: false},
	}

	var pr PostgreReader
	for _, test := range tests {
		schema, table, err := pr.splitFullName(test.fullname, test.sep)
		if (err != nil) != test.wanterr {
			t.Errorf("splitFullName: unexpected error: %v", test)
			continue
		}

		if schema != test.schema || table != test.table {
			t.Errorf("splitFullName: expected: %v.%v actual: %v.%v", test.schema, test.table, schema, table)
		}
	}
}
