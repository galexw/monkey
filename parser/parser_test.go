package parser

import (
	"fmt"
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
	if len(p.Errors()) == 0 {
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

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not an ast.PrefixExpression. Got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("ident.Operator not %s. Got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.RightExpression, tt.value) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. Got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not an ast.InfixExpression. Got %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.LeftExpression, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Errorf("ident.Operator not %s. Got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.RightExpression, tt.rightValue) {
			return
		}
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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *IntegerLiteral. Got %T", il)
		return false
	}

	if intLit.Value != value {
		t.Errorf("intLit.Value not %d. Got %d", value, intLit.Value)
		return false
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intLit.TokenLiteral() not %s. Got %s", fmt.Sprintf("%d", value), intLit.TokenLiteral())
		return false
	}

	return true
}
