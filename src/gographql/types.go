package gographql

// QueryArg represents a query argument for a Model
type QueryArg struct {
	Key   string
	Value interface{}
}

// Field represents a field in a block
type Field struct {
	Key   string
	Model *Model
}

// Block represents a {} block
type Block struct {
	Fields []*Field
}

// Model represents a model definition
type Model struct {
	Block     *Block
	Key       string
	QueryArgs []*QueryArg
}

// Token represents a character
type Token int

const (
	// ILLEGAL covers all illegal tokens
	ILLEGAL Token = iota
	// EOF - end-of-file
	EOF
	// WS - whitespace
	WS

	// Literals

	// IDENT - Identifiers. Can also be booleans in the case of query args
	IDENT
	// STRING - string
	STRING

	// INT - int64
	INT

	// FLOAT - float64
	FLOAT

	// BOOLEAN - booleancs
	BOOLEAN

	// Misc chars

	// LEFT_CURLY - {
	LEFT_CURLY

	// RIGHT_CURLY - }
	RIGHT_CURLY

	// LEFT_PARENTHESIS - (
	LEFT_PARENTHESIS

	// RIGHT_PARENTHESIS - )
	RIGHT_PARENTHESIS

	// COMMA - ,
	COMMA

	// COLON - :
	COLON
)
