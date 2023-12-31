package parser

import (
	"fmt"
	"galexw/monkey/ast"
	"galexw/monkey/lexer"
	"galexw/monkey/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
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

var precedences = map[token.TokenType]int{
	token.EQUAL:       EQUALS,
	token.NOTEQUAL:    EQUALS,
	token.LESSTHAN:    LESSGREATER,
	token.GREATERTHAN: LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.LEFTPAREN:   CALL,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	// I don't understand why identifiers and integer literals are using
	// prefix parse functions at the moment
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LEFTPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOTEQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESSTHAN, p.parseInfixExpression)
	p.registerInfix(token.GREATERTHAN, p.parseInfixExpression)

	// This is really cool
	p.registerInfix(token.LEFTPAREN, p.parseCallExpression)

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
		program.Statements = append(program.Statements, stmt)

		// Need to organize where to end a statement
		// Some statements end with a semicolon, some don't
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
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: nil,
	}

	// Still trying to understand what's LOWEST
	// "That's going to make more sense in a short while, I promise"
	expressionStatement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExpression := prefix()

	// Can visualize this as if the operator on the right is stronger,
	// it will absorb the left expression into the right side
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression) // This is where the magic happens
	}

	return leftExpression
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{
		Token: p.curToken,
	}

	// Parsing condition
	if !p.expectPeek(token.LEFTPAREN) {
		return nil
	}

	p.nextToken()
	condition := p.parseExpression(LOWEST)
	ifExpression.Condition = condition

	if !p.expectPeek(token.RIGHTPAREN) {
		return nil
	}

	if !p.expectPeek(token.LEFTBRACE) {
		return nil
	}

	consequenceBlock := p.parseBlock()
	ifExpression.Consequence = consequenceBlock

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LEFTBRACE) {
			return nil
		}

		alternativeBlock := p.parseBlock()
		ifExpression.Alternative = alternativeBlock
	}

	return ifExpression
}

func (p *Parser) parseBlock() *ast.Block {
	block := &ast.Block{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RIGHTBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	itl := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	itl.Value = value
	return itl
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHTPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	// Will soon understand why we need PREFIX
	pe.RightExpression = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	// We first parse the left expression, now we would be on the infix operator,
	// so we move to the next token, and then parse the right expression
	ie := &ast.InfixExpression{
		Token:          p.curToken,
		Operator:       p.curToken.Literal,
		LeftExpression: leftExpression,
	}

	precedence := p.curPrecedence()

	p.nextToken()

	ie.RightExpression = p.parseExpression(precedence)

	return ie
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	identifier := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	letStatement.Name = identifier

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	expression := p.parseExpression(LOWEST)
	letStatement.Value = expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()

	returnStatement.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return returnStatement
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LEFTPAREN) {
		return nil
	}

	functionLiteral.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFTBRACE) {
		return nil
	}

	functionLiteral.Body = p.parseBlock()

	return functionLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RIGHTPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Nice way to implement parsing a list of identifiers
	identifier := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		identifier := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RIGHTPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(leftExpression ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  leftExpression, // This is the identifier for the function
		Arguments: p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	if p.peekTokenIs(token.RIGHTPAREN) {
		p.nextToken()
		return arguments
	}

	p.nextToken()

	arguments = append(arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RIGHTPAREN) {
		return nil
	}

	return arguments
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

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekError(expectedTokenType token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected token %s, got %s", expectedTokenType, p.peekToken.Type))
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("No prefix parse function for token %s", t))
}
