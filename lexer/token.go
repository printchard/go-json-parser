package lexer

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenOpenBrace
	TokenCloseBrace
	TokenOpenBracket
	TokenCloseBracket
	TokenColon
	TokenComma
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenOpenBrace:
		return "OpenBrace"
	case TokenCloseBrace:
		return "CloseBrace"
	case TokenOpenBracket:
		return "OpenBracket"
	case TokenCloseBracket:
		return "CloseBracket"
	case TokenColon:
		return "Colon"
	case TokenComma:
		return "Comma"
	case TokenString:
		return "string"
	case TokenNumber:
		return "number"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNull:
		return "null"
	default:
		return "Unknown TokenType"
	}
}

type Token struct {
	Type  TokenType
	Value string
}
