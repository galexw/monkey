package parser

import (
	"galexw/monkey/ast"
	"galexw/monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		// Type assertion
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
	checkParserErrors(t, p)
}

func TestFailLetStatements(t *testing.T) {
	input := `
	let x 5;
	let = 10;
	let 838383;
	`
	l := lexer.New(input)
	p := New(l)

	p.ParseProgram()
	if len(p.Errors()) != 3 {
		t.Fatalf("ParseProgram() should have returned errors")
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		if stmt.TokenLiteral() != "return" {
			t.Fatalf("stmt.TokenLiteral() not 'return'. Got %q", stmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // Checking the program statement is an expression statement
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier) // Checking the expression is an identifier
	if !ok {
		t.Fatalf("stmt.Expression is not an ast.Identifier. Got %T", stmt.Expression)
	}

	if ident.Value != "foobar" { // Checking the identifier value is "foobar"
		t.Errorf("ident.Value not %s. Got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" { // Checking the identifier token literal is "foobar"
		t.Errorf("ident.TokenLiteral() not %s. Got %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // Checking the program statement is an expression statement
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral) // Checking the expression is an identifier
	if !ok {
		t.Fatalf("stmt.Expression is not an ast.IntegerLiteral. Got %T", stmt.Expression)
	}

	if ident.Value != 5 { // Checking the identifier value is 5
		t.Errorf("ident.Value not %d. Got %d", 5, ident.Value)
	}

	if ident.TokenLiteral() != "5" { // Checking the identifier token literal is "5"
		t.Errorf("ident.TokenLiteral() not %s. Got %s", "5", ident.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Fatalf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Fatalf("parser error: %q", msg)
	}
}

// Why's it not testing the actual expression value? I think its TBD
// My guess is there's a bunch of expression types
// In the sneak peek on pg 48 there's OperatorExpression, IntegerLiteral etc
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. Got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *LetStatement. Got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. Got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s. Got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
