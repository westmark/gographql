package gographql_test

import (
	"reflect"
	"strings"
	"testing"

	ql "gographql"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s     string
		block *ql.Block
		err   string
	}{
		// Empty block statement
		{
			s:     `{}`,
			block: &ql.Block{},
		},

		// Single field
		{
			s: `{
        foo
      }`,
			block: &ql.Block{Fields: []*ql.Field{&ql.Field{Key: "foo"}}},
		},

		// Multiple fields
		{
			s: `{
        foo,
        bar
      }`,
			block: &ql.Block{Fields: []*ql.Field{&ql.Field{Key: "foo"}, &ql.Field{Key: "bar"}}},
		},

		// Multiple fields and model
		{
			s: `{
        foo,
        bar,
        user(id:1) {}

      }`,
			block: &ql.Block{Fields: []*ql.Field{
				&ql.Field{Key: "foo"},
				&ql.Field{Key: "bar"},
				&ql.Field{Key: "user", Model: &ql.Model{
					Key:       "user",
					QueryArgs: []*ql.QueryArg{&ql.QueryArg{Key: "id", Value: int64(1)}},
					Block:     &ql.Block{Fields: []*ql.Field{}},
				}},
			}},
		},

		{
			s: `{
        foo,
        bar,
        user(id:1) {
					name,
					age
				},
				store(address:"street", zip:"1337", active:true) {}

      }`,
			block: &ql.Block{Fields: []*ql.Field{
				&ql.Field{Key: "foo"},
				&ql.Field{Key: "bar"},
				&ql.Field{Key: "user", Model: &ql.Model{
					Key:       "user",
					QueryArgs: []*ql.QueryArg{&ql.QueryArg{Key: "id", Value: int64(1)}},
					Block: &ql.Block{Fields: []*ql.Field{
						&ql.Field{Key: "name"},
						&ql.Field{Key: "age"},
					}},
				}},
				&ql.Field{Key: "store", Model: &ql.Model{
					Key: "store",
					QueryArgs: []*ql.QueryArg{
						&ql.QueryArg{Key: "address", Value: "street"},
						&ql.QueryArg{Key: "zip", Value: "1337"},
						&ql.QueryArg{Key: "active", Value: true},
					},
					Block: &ql.Block{Fields: []*ql.Field{}},
				}},
			}},
		},

		/*
			// Multi-field statement
			{
				s: `SELECT first_name, last_name, age FROM my_table`,
				stmt: &ql.Block{
					Fields:    []string{"first_name", "last_name", "age"},
					TableName: "my_table",
				},
			},

			// Select all statement
			{
				s: `SELECT * FROM my_table`,
				stmt: &ql.Block{
					Fields:    []string{"*"},
					TableName: "my_table",
				},
			},

			// Errors
			{s: `foo`, err: `found "foo", expected SELECT`},
			{s: `SELECT !`, err: `found "!", expected field`},
			{s: `SELECT field xxx`, err: `found "xxx", expected FROM`},
			{s: `SELECT field FROM *`, err: `found "*", expected table name`},
		*/
	}

	for i, tt := range tests {
		block, err := ql.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !compareBlocks(tt.block, block) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.block, block)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func compareQueryArgs(a, b *ql.QueryArg) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	return (a.Key == b.Key && a.Value == b.Value)
}

func compareModels(a, b *ql.Model) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	if a == nil && b == nil {
		return true
	}

	if len(a.QueryArgs) != len(b.QueryArgs) {
		return false
	}

	for i := range a.QueryArgs {
		if !compareQueryArgs(a.QueryArgs[i], b.QueryArgs[i]) {
			return false
		}
	}

	return compareBlocks(a.Block, b.Block)
}

func compareFields(a, b *ql.Field) bool {
	if a.Key != b.Key {
		return false
	}

	return compareModels(a.Model, b.Model)
}

func compareBlocks(a, b *ql.Block) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	if a == nil && b == nil {
		return true
	}

	if len(a.Fields) != len(b.Fields) {
		return false
	}

	for i := range a.Fields {
		if !compareFields(a.Fields[i], b.Fields[i]) {
			return false
		}
	}

	return true
}
