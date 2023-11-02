package lexer

import (
	"galexw/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer { // Input here is actually the source code in Monkey
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token { // ch is the character that is being read
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	// TODO: Implement next token method
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Operators
	case '=':
		tok = newToken(token.Assign, l.ch)
	case '+':
		tok = newToken(token.Plus, l.ch)
	case '-':
		tok = newToken(token.Minus, l.ch)
	case '!':
		tok = newToken(token.Bang, l.ch)
	case '*':
		tok = newToken(token.Asterisk, l.ch)
	case '/':
		tok = newToken(token.Slash, l.ch)
	case '<':
		tok = newToken(token.LessThan, l.ch)
	case '>':
		tok = newToken(token.GreaterThan, l.ch)

	// Delimiters
	case ';':
		tok = newToken(token.Semicolon, l.ch)
	case '(':
		tok = newToken(token.LeftParen, l.ch)
	case ')':
		tok = newToken(token.RightParen, l.ch)
	case ',':
		tok = newToken(token.Comma, l.ch)
	case '{':
		tok = newToken(token.LeftBrace, l.ch)
	case '}':
		tok = newToken(token.RightBrace, l.ch)
	case '[':
		tok = newToken(token.LeftBracket, l.ch)
	case ']':
		tok = newToken(token.RightBracket, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // If the character is a letter
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) { // If the character is a digit
			tok.Type = token.Int
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch) // If the character is not a letter or a digit
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' { // While the character is a whitespace
		l.readChar()
	}
}

// Helper functions
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		ch == '_' // Allow underscores in identifiers
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) { // While the character is a letter
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) { // While the character is a digit
		l.readChar()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' // If the character is a digit
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}
