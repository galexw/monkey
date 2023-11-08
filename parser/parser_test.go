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

	// Test the first let statement
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
}

// Helper function for TestLetStatements
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. Got %q", s.TokenLiteral())
		return false
	}

	// Type assertion
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *LetStatement. Got %T", s)
		return false
	}

	// Test the identifier
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. Got %s", name, letStmt.Name.Value)
		return false
	}

	// Test the identifier's token literal
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s. Got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
