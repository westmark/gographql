package gographql

// QueryArg represents a query argument for a Model
type QueryArg struct {
	Key   string
	Value interface{}
}

// Field represents a field in a Block
type Field struct {
	Key   string
	Model *Model
}

// Block represents a block of Fields
type Block struct {
	Fields []*Field
}

// Model represents a model definition complete with a query and a Block
type Model struct {
	Block     *Block
	Key       string
	QueryArgs []*QueryArg
}

// Token represents one or several runes matching allowed symbols in GraphQL
type Token int

const (
	// Illegal covers all Illegal tokens
	Illegal Token = iota

	// EOF - end of file
	EOF

	// WS - any whitespace
	WS

	// Literals

	// Ident - Identifiers. Can also be booleans in the case of query args
	Ident

	// String - A string token
	String

	// Int - An int64 token
	Int

	// Float - A float64 token
	Float

	// Boolean - A boolean token
	Boolean

	// Misc tokens

	// LeftCurly - {
	LeftCurly

	// RightCurly - }
	RightCurly

	// LeftParenthesis - (
	LeftParenthesis

	// RightParenthesis - )
	RightParenthesis

	// Comma - ,
	Comma

	// Colon - :
	Colon
)
