package gographql

import (
	"fmt"
	"io"
	"strconv"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) parseQueryArgs() ([]*QueryArg, error) {

	var queryArgs []*QueryArg
	var err error

	if tok, lit := p.scanIgnoreWhitespace(); tok != LeftParenthesis {
		return nil, fmt.Errorf("found %q, expected LeftParenthesis", lit)
	}

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != Ident {
			return nil, fmt.Errorf("found %q, expected Ident", lit)
		}

		key := lit

		if tok, lit := p.scanIgnoreWhitespace(); tok != Colon {
			return nil, fmt.Errorf("found %q, expected Colon", lit)
		}

		tok, lit = p.scanIgnoreWhitespace()
		if tok != Ident && tok != String && tok != Int && tok != Float {
			return nil, fmt.Errorf("found %q, expected Boolean, String, Int or Float", lit)
		}

		var value interface{}
		value = lit

		err = nil
		if tok == Ident {
			if lit == "true" {
				value = true
			} else if lit == "false" {
				value = false
			} else {
				return nil, fmt.Errorf("found %q, expected Boolean, String, Int or Float", lit)
			}
		} else if tok == Int {
			value, err = strconv.ParseInt(lit, 10, 64)
		} else if tok == Float {
			value, err = strconv.ParseFloat(lit, 64)
		}

		if err != nil {
			return nil, err
		}

		queryArgs = append(queryArgs, &QueryArg{key, value})

		if tok, _ := p.scanIgnoreWhitespace(); tok != Comma {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != RightParenthesis {
		return nil, fmt.Errorf("found %q, expected RightParenthesis", lit)
	}

	return queryArgs, nil
}

func (p *Parser) parseField() (*Field, error) {

	field := &Field{}
	var err error

	tok, lit := p.scanIgnoreWhitespace()
	if tok != Ident {
		return nil, fmt.Errorf("found %q, expected Ident", lit)
	}

	field.Key = lit

	tok, lit = p.scanIgnoreWhitespace()
	p.unscan()

	if tok != LeftParenthesis {
		return field, nil
	}

	// We have a ModelBlock
	field.Model = &Model{Key: field.Key}

	if field.Model.QueryArgs, err = p.parseQueryArgs(); err != nil {
		return nil, err
	}

	tok, lit = p.scanIgnoreWhitespace()
	p.unscan()
	if tok == LeftCurly {
		field.Model.Block, err = p.parseBlock()
		if err != nil {
			return nil, err
		}

	}

	return field, nil

}

func (p *Parser) parseBlock() (*Block, error) {
	block := &Block{}

	var tok Token
	var lit string

	if tok, lit = p.scanIgnoreWhitespace(); tok != LeftCurly {
		return nil, fmt.Errorf("found %q, expected LeftCurly", lit)
	}

	tok, _ = p.scanIgnoreWhitespace()
	if tok != RightCurly {
		p.unscan()
	}

	for tok != RightCurly {
		field, err := p.parseField()

		if err != nil {
			return nil, err
		}

		block.Fields = append(block.Fields, field)
		tok, lit = p.scanIgnoreWhitespace()
		if tok != Comma && tok != RightCurly {
			return nil, fmt.Errorf("found %q, expected Comma or RightCurly", lit)
		}
	}

	return block, nil
}

// Parse returns the top level block
func (p *Parser) Parse() (*Block, error) {

	return p.parseBlock()

}
