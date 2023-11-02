package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	// Identifiers + Literals
	Identifier = "Identifier" // add, x ,y, ...
	Int        = "Int"        // 123456
	String     = "String"     // "x", "y"

	// Operators
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	Equal    = "=="
	NotEqual = "!="

	LessThan    = "<"
	GreaterThan = ">"

	// Delimiters
	Comma     = ","
	Semicolon = ";"
	Colon     = ":"

	LeftParen    = "("
	RightParen   = ")"
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"

	// Keywords
	Function = "Function"
	Let      = "Let"
	True     = "True"
	False    = "False"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
)

// LookupIdent checks if the given identifier is a keyword and returns the corresponding TokenType.
func LookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return Function
	case "let":
		return Let
	case "true":
		return True
	case "false":
		return False
	case "if":
		return If
	case "else":
		return Else
	case "return":
		return Return
	default:
		return Identifier
	}
}
