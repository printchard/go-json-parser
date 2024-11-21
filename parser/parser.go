package parser

import (
	"fmt"
	"strconv"

	"github.com/printchard/go-json-parser/lexer"
)

type Parser struct {
	input  string
	tokens []lexer.Token
	pos    int
}

func New(input string) Parser {
	return Parser{input: input}
}

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{}
	}

	return p.tokens[p.pos]
}

func (p *Parser) next() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{}
	}
	token := p.tokens[p.pos]
	p.pos++

	return token
}

func (p *Parser) parseObject() (map[string]any, error) {
	p.next()
	obj := make(map[string]any)

	for {
		t := p.peek()
		if t.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF")
		}

		if t.Type == lexer.TokenCloseBrace {
			p.next()
			return obj, nil
		}

		if t.Type != lexer.TokenString {
			return nil, fmt.Errorf("expected string in key")
		}

		key := p.next()
		if key.Value == "" {
			return nil, fmt.Errorf("empty key")
		}
		if _, exists := obj[key.Value]; exists {
			return nil, fmt.Errorf("duplicate string %s", key.Value)
		}

		t = p.peek()
		if t.Type != lexer.TokenColon {
			return nil, fmt.Errorf("expected colon after key")
		}
		p.next()

		var val any
		var err error
		t = p.peek()
		switch t.Type {
		case lexer.TokenOpenBrace:
			val, err = p.parseObject()
		case lexer.TokenOpenBracket:
			val, err = p.parseArray()
		case lexer.TokenString:
			val = t.Value
			p.next()
		case lexer.TokenNumber:
			val, err = strconv.ParseFloat(t.Value, 64)
			p.next()
		case lexer.TokenTrue:
			val = true
			p.next()
		case lexer.TokenFalse:
			val = false
			p.next()
		case lexer.TokenNull:
			val = nil
			p.next()
		default:
			err = fmt.Errorf("invalid value after key")
		}

		if err != nil {
			return nil, err
		}
		obj[key.Value] = val

		t = p.peek()
		if t.Type == lexer.TokenComma {
			p.next()
			if p.peek().Type == lexer.TokenCloseBrace {
				return nil, fmt.Errorf("trailing comma")
			}
		} else if t.Type != lexer.TokenCloseBrace {
			return nil, fmt.Errorf("expected closing brace")
		}
	}
}

func (p *Parser) parseArray() ([]any, error) {
	p.next()
	arr := make([]any, 0)

	for {
		t := p.peek()

		if t.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF")
		}

		if t.Type == lexer.TokenCloseBracket {
			p.next()
			return arr, nil
		}

		var val any
		var err error
		switch t.Type {
		case lexer.TokenOpenBrace:
			val, err = p.parseObject()
		case lexer.TokenOpenBracket:
			val, err = p.parseArray()
		case lexer.TokenString:
			val = t.Value
			p.next()
		case lexer.TokenNumber:
			val, err = strconv.ParseFloat(t.Value, 64)
			p.next()
		case lexer.TokenTrue:
			val = true
			p.next()
		case lexer.TokenFalse:
			val = false
			p.next()
		case lexer.TokenNull:
			val = nil
			p.next()
		}

		if err != nil {
			return nil, err
		}
		arr = append(arr, val)

		t = p.peek()
		if t.Type == lexer.TokenComma {
			p.next()
			if p.peek().Type == lexer.TokenCloseBracket {
				return nil, fmt.Errorf("trailing comma")
			}
		} else if t.Type != lexer.TokenCloseBracket {
			return nil, fmt.Errorf("expected closing bracket")
		}
	}
}

func (p *Parser) Parse() (any, error) {
	l := lexer.New(p.input)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	p.tokens = tokens

	t := p.peek()
	if t.Type == lexer.TokenOpenBrace {
		return p.parseObject()
	} else if t.Type == lexer.TokenOpenBracket {
		return p.parseArray()
	}

	return nil, fmt.Errorf("invalid input")
}
