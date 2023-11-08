package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "Illegal"
	EOF     = "EOF"

	// Identifiers + Literals
	IDENTIFIER = "Identifier" // add, x ,y, ...
	INT        = "Int"        // 123456
	STRING     = "String"     // "x", "y"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	EQUAL    = "=="
	NOTEQUAL = "!="

	LESSTHAN    = "<"
	GREATERTHAN = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LEFTPAREN    = "("
	RIGHTPAREN   = ")"
	LEFTBRACE    = "{"
	RIGHTBRACE   = "}"
	LEFTBRACKET  = "["
	RIGHTBRACKET = "]"

	// Keywords
	FUNCTION = "Function"
	LET      = "Let"
	TRUE     = "True"
	FALSE    = "False"
	IF       = "If"
	ELSE     = "Else"
	RETURN   = "Return"
)

// LookupIdent checks if the given identifier is a keyword and returns the corresponding TokenType.
func LookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "if":
		return IF
	case "else":
		return ELSE
	case "return":
		return RETURN
	default:
		return IDENTIFIER
	}
}
