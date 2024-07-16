package lexer

import (
	"monkey-lang/token"
)

type Lexer struct {
	input   string
	pos     int  // curr position in input (points to curr char)
	readPos int  // current reading position in input (after current char)
	ch      byte // current char under examination, represented in byte
}

func New(input string) *Lexer {
	// GO makes sure lexer doesn't get deallocated since it still has a refernce to it
	lexer := Lexer{input: input}

	lexer.readChar()

	return &lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.ch {
	// Delimiters
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)

	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)

	// Operators
	case '=':
		tok = newToken(token.ASSIGN, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Determining if it is a keyword or iden
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			// Neccessary to avoid the later readChar() after the switch
			return tok

		} else if isDigit(lexer.ch) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) readNumber() string {
	starting_pos := lexer.pos
	for isDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[starting_pos:lexer.pos]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

// Reads until end of the word
// Returns the keyword / iden
func (lexer *Lexer) readIdentifier() string {
	start_pos := lexer.pos

	for isLetter(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[start_pos:lexer.pos]
}

// True if alpha or _
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Moves pos and readPos until end of input
func (lexer *Lexer) readChar() {
	if len(lexer.input) <= lexer.readPos {
		lexer.ch = 0 // Indicates ascii code for null or haven't read yet
	} else {
		lexer.ch = lexer.input[lexer.readPos]
	}

	lexer.pos = lexer.readPos

	lexer.readPos += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
