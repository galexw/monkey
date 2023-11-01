package lexer

import "galexw/monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer { // Input here is actually the source code in Monkey
	l := &Lexer{input: input}
	return l
}

func (l *Lexer) NextToken() token.Token {
	// TODO: Implement next token method
	return token.Token{}
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
