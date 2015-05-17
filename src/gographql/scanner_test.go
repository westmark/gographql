package gographql_test

import (
	"strings"
	"testing"

	ql "gographql"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok ql.Token
		lit interface{}
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: ql.EOF, lit: ""},
		{s: `#`, tok: ql.ILLEGAL, lit: `#`},
		{s: ` `, tok: ql.WS, lit: " "},
		{s: "\t", tok: ql.WS, lit: "\t"},
		{s: "\n", tok: ql.WS, lit: "\n"},

		// Misc characters
		{s: `{`, tok: ql.LEFT_CURLY, lit: "{"},
		{s: `}`, tok: ql.RIGHT_CURLY, lit: "}"},
		{s: `(`, tok: ql.LEFT_PARENTHESIS, lit: "("},
		{s: `)`, tok: ql.RIGHT_PARENTHESIS, lit: ")"},
		{s: `:`, tok: ql.COLON, lit: ":"},

		// Identifiers
		{s: `foo`, tok: ql.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: ql.IDENT, lit: `Zx12_3U_`},

		{s: `"foo"`, tok: ql.STRING, lit: `foo`},

		{s: `20`, tok: ql.INT, lit: "20"},

		{s: `20.5`, tok: ql.FLOAT, lit: "20.5"},

		// Keywords
	}

	for i, tt := range tests {
		s := ql.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
