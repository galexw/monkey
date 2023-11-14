package ast

import "galexw/monkey/token"

type Node interface {
	TokenLiteral() string
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

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (i *ReturnStatement) statementNode()       {}
func (i *ReturnStatement) TokenLiteral() string { return i.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (i *ExpressionStatement) statementNode()       {}
func (i *ExpressionStatement) TokenLiteral() string { return i.Token.Literal }
