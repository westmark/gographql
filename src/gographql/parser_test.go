package gographql_test

import (
	"fmt"
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
		} else if tt.err == "" && !reflect.DeepEqual(tt.block, block) {
			printBlock(tt.block)
			printBlock(block)
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

func printBlock(block *ql.Block) {
	fmt.Println("<block>")
	for _, field := range block.Fields {
		fmt.Println(field.Key)
		if field.Model != nil {
			fmt.Println(field.Model.Key)

			for _, qa := range field.Model.QueryArgs {
				fmt.Print(qa.Key + " ")
				fmt.Println(qa.Value)
				_, ok := qa.Value.(int64)
				fmt.Println(ok)
			}

			if field.Model.Block != nil {
				printBlock(field.Model.Block)
			}
		}
	}
}
