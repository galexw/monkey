package parser

import (
	"fmt"
	"galexw/monkey/ast"
	"galexw/monkey/lexer"
	"galexw/monkey/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	// As of now - I do not understand what's the point of these maps of token
	// types to functions. I'll probably understand it when I get to the
	// expression parsing part.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letToken := p.curToken

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	identifierPtr := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Implement parsing expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &ast.LetStatement{
		Token: letToken,
		Name:  identifierPtr,
		Value: nil,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnToken := p.curToken

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &ast.ReturnStatement{
		Token:       returnToken,
		ReturnValue: nil, // TODO: Parse return expression
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(expectedTokenType token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected token %s, got %s", expectedTokenType, p.peekToken.Type))
}
