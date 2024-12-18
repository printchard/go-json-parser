package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	input string
	pos   int
}

func New(input string) Lexer {
	return Lexer{
		input: input,
		pos:   0,
	}
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) next() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	char := l.input[l.pos]
	l.pos++
	return char
}

func (l *Lexer) skipWhitespace() {
	for {
		if l.pos >= len(l.input) || !unicode.IsSpace(rune(l.peek())) {
			break
		}
		l.next()
	}
}

func (l *Lexer) parseString() (string, error) {
	l.next()
	var str strings.Builder
	for {
		if l.pos >= len(l.input) {
			return "", fmt.Errorf("invalid string: unterminated string")
		}
		char := l.next()

		if char == '"' {
			break
		}

		if char == '\\' {
			char = l.next()
		}
		str.WriteByte(char)
	}

	return str.String(), nil
}

func (l *Lexer) parseNumber() (string, error) {
	start := l.pos

	if l.peek() == '-' {
		l.next()
	}

	if l.peek() == '0' {
		l.next()
		if unicode.IsDigit(rune(l.peek())) {
			return "", fmt.Errorf("invalid number: leading 0")
		}
	} else if unicode.IsDigit(rune(l.peek())) {
		for unicode.IsDigit(rune(l.peek())) {
			l.next()
		}
	} else {
		return "", fmt.Errorf("invalid number: contains non number")
	}

	if l.peek() == '.' {
		l.next()
		if !unicode.IsDigit(rune(l.peek())) {
			return "", fmt.Errorf("invalid number: contains non number after decimal point")
		}
		for unicode.IsDigit(rune(l.peek())) {
			l.next()
		}
	}

	if l.peek() == 'e' || l.peek() == 'E' {
		l.next()
		if l.peek() == '+' || l.peek() == '-' {
			l.next()
		}
		if !unicode.IsDigit(rune(l.peek())) {
			return "", fmt.Errorf("invalid number: contains non number after exponent")
		}
		for unicode.IsDigit(rune(l.peek())) {
			l.next()
		}
	}

	return l.input[start:l.pos], nil
}

func (l *Lexer) matchLiteral(literal string) bool {
	if len(l.input)-l.pos < len(literal) {
		return false
	}

	if l.input[l.pos:l.pos+len(literal)] == literal {
		l.pos += len(literal)
		return true
	}

	return false
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		l.skipWhitespace()
		char := l.peek()

		if char == 0 {
			tokens = append(tokens, Token{Type: TokenEOF})
			return tokens, nil
		}

		switch char {
		case '{':
			tokens = append(tokens, Token{Type: TokenOpenBrace})
			l.next()
		case '}':
			tokens = append(tokens, Token{Type: TokenCloseBrace})
			l.next()
		case '[':
			tokens = append(tokens, Token{Type: TokenOpenBracket})
			l.next()
		case ']':
			tokens = append(tokens, Token{Type: TokenCloseBracket})
			l.next()
		case ':':
			tokens = append(tokens, Token{Type: TokenColon})
			l.next()
		case ',':
			tokens = append(tokens, Token{Type: TokenComma})
			l.next()
		case '"':
			str, err := l.parseString()
			if err != nil {
				return nil, fmt.Errorf("invalid input at: %s", err)
			}
			tokens = append(tokens, Token{Type: TokenString, Value: str})
		default:
			if unicode.IsDigit(rune(char)) || char == '-' {
				num, err := l.parseNumber()
				if err != nil {
					return nil, fmt.Errorf("invalid input: %s", err)
				}
				tokens = append(tokens, Token{Type: TokenNumber, Value: num})
			} else if l.matchLiteral("true") {
				tokens = append(tokens, Token{Type: TokenTrue})
			} else if l.matchLiteral("false") {
				tokens = append(tokens, Token{Type: TokenFalse})
			} else if l.matchLiteral("null") {
				tokens = append(tokens, Token{Type: TokenNull})
			} else {
				return nil, fmt.Errorf("invalid input: unexpected character")
			}
		}
	}
}
