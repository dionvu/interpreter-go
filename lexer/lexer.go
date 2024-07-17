package lexer

// PG 24

import (
	"monkey-lang/token"
	"monkey-lang/utils"
)

type Lexer struct {
	input   string
	pos     int  // curr position in input (points to curr char)
	readPos int  // current reading position in input (after current char)
	ch      byte // current char under examination, represented in byte
}

// Returns new instance of a lexer
func New(input string) *Lexer {
	// GO makes sure lexer doesn't get deallocated since it still has a refernce to it
	// ch = 0 is the null byte
	lexer := Lexer{input: input, pos: -1, readPos: 0, ch: 0}

	lexer.readChar()

	return &lexer
}

// Returns a token with the correct token literal, and token type.
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.ch {

	// Bracket and parentheses
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
	case '/':
		tok = newToken(token.SLASH, lexer.ch)

	// Operators
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case '-':
		tok = newToken(token.MINUS, lexer.ch)
	case '*':
		tok = newToken(token.ASTERISK, lexer.ch)
	case '<':
		tok = newToken(token.LT, lexer.ch)
	case '>':
		tok = newToken(token.GT, lexer.ch)
	case '!':
		if lexer.peekChar() == '=' {

			lexer.readChar()

			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, '!')
		}
	case '=':
		if lexer.peekChar() == '=' {

			lexer.readChar()

			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, '=')
		}

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	// Handles identifiers and keywords
	//
	default:
		if utils.IsLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			// Neccessary to avoid the later readChar() after the switch
			return tok

		} else if utils.IsDigit(lexer.ch) {
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

// Moves pos and readPos until end of input
func (lexer *Lexer) readChar() {
	// Next char to read has exceeded the length
	if len(lexer.input) <= lexer.readPos {
		lexer.ch = 0 // Indicates ascii code for null or haven't read yet
	} else {
		lexer.ch = lexer.input[lexer.readPos]
	}

	lexer.pos = lexer.readPos

	lexer.readPos += 1
}

func (lexer *Lexer) readNumber() string {
	starting_pos := lexer.pos
	for utils.IsDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[starting_pos:lexer.pos]
}

func (lexer *Lexer) peekChar() byte {
	if lexer.pos >= len(lexer.input) {
		// null
		return 0
	} else {
		return lexer.input[lexer.readPos]
	}
}

// Reads until end of the word
// Returns the keyword / iden
func (lexer *Lexer) readIdentifier() string {
	start_pos := lexer.pos

	for utils.IsLetter(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[start_pos:lexer.pos]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}
