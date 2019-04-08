package gostructcompose

import (
	"reflect"
	"testing"
)

func TestNewGenerator(t *testing.T) {

	tests := []struct {
		loc   string
		items []Item
		dest  string
		mr    MetaReader
		ew    EntityWriter
		tc    TypeConverter
		//err   error
		wanterr bool
	}{

		/* {loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), err: nil},
		{loc: "", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), err: errors.New("failed")},
		{loc: "http://", items: nil, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), err: errors.New("failed")},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), err: errors.New("failed")},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: nil, ew: new(GolangWriter), tc: new(PgGoConverter), err: errors.New("failed")},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: nil, tc: new(PgGoConverter), err: errors.New("failed")},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: nil, err: errors.New("failed")}, */
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), wanterr: false},
		{loc: "", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), wanterr: true},
		{loc: "http://", items: nil, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), wanterr: true},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "", mr: new(PostgreReader), ew: new(GolangWriter), tc: new(PgGoConverter), wanterr: true},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: nil, ew: new(GolangWriter), tc: new(PgGoConverter), wanterr: true},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: nil, tc: new(PgGoConverter), wanterr: true},
		{loc: "http://", items: []Item{{FullName: "public.table"}}, dest: "outfile.go", mr: new(PostgreReader), ew: new(GolangWriter), tc: nil, wanterr: true},
	}

	for _, test := range tests {
		_, err := NewGenerator(test.loc, test.items, test.dest, test.mr, test.ew, test.tc)
		if (err != nil) != test.wanterr {
			//if (test.err == nil && err != nil) || (test.err != nil && err == nil) {
			t.Errorf("NewGenerator: %v", test)
		}
	}

}

func TestGetEntity(t *testing.T) {
	tests := []struct {
		t       Table
		e       Entity
		wanterr bool
	}{
		{Table{Name: "", Columns: []Column{}},
			Entity{Name: "", Attrs: []Attribute{}}, true},
		{Table{Name: "table_1", Columns: []Column{}},
			Entity{Name: "table_1", Attrs: []Attribute{}}, false},
		{Table{Name: "table_1", Columns: []Column{{Name: "col_1", DataType: "integer", IsNullable: "YES"}}},
			Entity{Name: "table_1", Attrs: []Attribute{{Name: "col_1", Type: "*int"}}}, false},
		{Table{Name: "table_1", Columns: []Column{{Name: "col_1", DataType: "integer", IsNullable: "YES"}, {Name: "col_2", DataType: "integer", IsNullable: "YES"}}},
			Entity{Name: "table_1", Attrs: []Attribute{{Name: "col_1", Type: "*int"}, {Name: "col_2", Type: "*int"}}}, false},
	}

	g, err := NewGenerator("http://", []Item{{FullName: "public.table"}}, "outfile.go", new(PostgreReader), new(GolangWriter), new(PgGoConverter))
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, test := range tests {
		e, err := g.getEntity(&test.t)
		if (err != nil) != test.wanterr {
			t.Errorf("getEntity: unexpected error: src: %v dest: %v", test.t, test.e)
			continue
		}

		if e != nil {
			if !reflect.DeepEqual(*e, test.e) {
				t.Errorf("getEntity !DeepEqual: src: %v dest: %v", test.t, test.e)
			}
		}

	}

}
