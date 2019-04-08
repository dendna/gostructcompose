package gostructcompose

import "testing"

func TestConvert(t *testing.T) {
	tests := []struct {
		src      string
		nullable string
		dest     string
		wanterr  bool
	}{
		{src: "INTEGER", nullable: "", dest: "", wanterr: true},
		{src: "integer", nullable: "NO", dest: "int", wanterr: false},
		{src: "integer", nullable: "YES", dest: "*int", wanterr: false},
		{src: "double precision", nullable: "NO", dest: "float64", wanterr: false},
		{src: "double precision", nullable: "YES", dest: "*float64", wanterr: false},
		{src: "text", nullable: "NO", dest: "string", wanterr: false},
		{src: "text", nullable: "YES", dest: "*string", wanterr: false},
		{src: "boolean", nullable: "NO", dest: "bool", wanterr: false},
		{src: "boolean", nullable: "YES", dest: "*bool", wanterr: false},
		// TODO: add remaining data types
	}

	var converter PgGoConverter

	for _, test := range tests {
		dst, err := converter.Convert(test.src, test.nullable)
		if (err != nil) != test.wanterr {
			t.Errorf("PgGoConverter: unexpected error: %v", test)
			continue
		}

		if dst != test.dest {
			t.Errorf("PgGoConverter: expected: %v actual: %v", test.dest, dst)
		}
	}

}
