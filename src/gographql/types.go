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
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals,
	IDENT // field names
	STRING
	INT
	FLOAT
	BOOLEAN

	// Misc chars
	LEFT_CURLY        // {
	RIGHT_CURLY       // }
	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	COMMA             // ,
	COLON             // :
)
