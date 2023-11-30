package ast

import (
	"bytes"
	"galexw/monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// Still figuring out what this is trying to do
// But I'll probably understand it when the whole thing pieces together
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil { // To remove later when we have expressions
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (i *Boolean) expressionNode()      {}
func (i *Boolean) TokenLiteral() string { return i.Token.Literal }
func (i *Boolean) String() string       { return i.Token.Literal }

// <prefix><expression>
type PrefixExpression struct {
	Token           token.Token
	Operator        string
	RightExpression Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

// <expression><infix><expression>
type InfixExpression struct {
	Token           token.Token
	LeftExpression  Expression
	Operator        string
	RightExpression Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.LeftExpression.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression // A condition is an expression that produces a boolean value
	Consequence *Block
	Alternative *Block
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

type Block struct {
	Token      token.Token // The { token
	Statements []Statement
}

func (i *Block) statementNode()       {}
func (i *Block) TokenLiteral() string { return i.Token.Literal }
func (i *Block) String() string {
	var out bytes.Buffer
	for _, s := range i.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *Block
}

func (fe *FunctionLiteral) expressionNode()      {}
func (fe *FunctionLiteral) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("fn")
	out.WriteString("(")
	for i, p := range fe.Parameters {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") {")
	out.WriteString(fe.Body.String())
	out.WriteString("}")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (i *ReturnStatement) statementNode()       {}
func (i *ReturnStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(i.TokenLiteral() + " ")

	if i.ReturnValue != nil { // To remove later when we have expressions
		out.WriteString(i.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (i *ExpressionStatement) statementNode()       {}
func (i *ExpressionStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ExpressionStatement) String() string {
	return i.Expression.String() + ";"
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (i *CallExpression) expressionNode()      {}
func (i *CallExpression) TokenLiteral() string { return i.Token.Literal }
func (i *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(i.Function.String())
	out.WriteString("(")
	for i, p := range i.Arguments {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(")")
	return out.String()
}
