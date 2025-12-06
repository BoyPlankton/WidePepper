// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package ast defines the Abstract Syntax Tree nodes and structure for the WidePepper language.
// It represents the hierarchical structure of parsed WidePepper programs.
package ast

import (
	"fmt"
	"strings"
	"widepepper/token"
)

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is the base interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression is the base interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// --- AST: Expressions ---

// Identifier Expression
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral Expression (used for all numbers, including floats)
type IntegerLiteral struct {
	Token token.Token // the INT token
	Value float64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral Expression
type StringLiteral struct {
	Token token.Token // the STR token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return fmt.Sprintf(`"%s"`, sl.Value) }

// BooleanLiteral Expression
type BooleanLiteral struct {
	Token token.Token // the CATNIP or MICE token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// Prefix Expression (e.g., -5, !catnip)
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g., ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Infix Expression (e.g., 5 + 5, a > b)
type InfixExpression struct {
	Token    token.Token // The operator token, e.g., +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Call Expression (e.g., add(2, 3), function_name(arg1, arg2, ...))
type CallExpression struct {
	Token     token.Token  // the LPAREN token
	Function  Expression   // the function name (identifier)
	Arguments []Expression // the arguments
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out strings.Builder
	out.WriteString(ce.Function.String())
	out.WriteString("(")

	for i, arg := range ce.Arguments {
		out.WriteString(arg.String())
		if i < len(ce.Arguments)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")
	return out.String()
}

// --- AST: Statements ---

// Variable Declaration Statement (`yarn`)
type YarnStatement struct {
	Token token.Token // the YARN token
	Name  *Identifier
	Value Expression
}

func (ys *YarnStatement) statementNode()       {}
func (ys *YarnStatement) TokenLiteral() string { return ys.Token.Literal }
func (ys *YarnStatement) String() string {
	var out strings.Builder
	out.WriteString(ys.TokenLiteral() + " ")
	out.WriteString(ys.Name.String())
	out.WriteString(" = ")

	if ys.Value != nil {
		out.WriteString(ys.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Block Statement (`{...}` block of statements)
type BlockStatement struct {
	Token      token.Token // the LBRACE token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
		out.WriteString(" ")
	}
	out.WriteString("}")
	return out.String()
}

// Pounce Statement (`pounce (condition) { block } [hiss { block }]`)
type PounceStatement struct {
	Token       token.Token // the POUNCE token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // Optional hiss block
}

func (ps *PounceStatement) statementNode()       {}
func (ps *PounceStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PounceStatement) String() string {
	var out strings.Builder
	out.WriteString("pounce (")
	out.WriteString(ps.Condition.String())
	out.WriteString(") ")
	out.WriteString(ps.Consequence.String())

	if ps.Alternative != nil {
		out.WriteString(" hiss ")
		out.WriteString(ps.Alternative.String())
	}

	return out.String()
}

// Return Statement (`treat expression;`)
type ReturnStatement struct {
	Token       token.Token // the TREAT token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString(rs.TokenLiteral())
	if rs.ReturnValue != nil {
		out.WriteString(" ")
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// Function Definition Statement (`purr function_name(param1, param2, ...) { body }`)
type PurrStatement struct {
	Token      token.Token   // the PURR token
	Name       *Identifier   // function name
	Parameters []*Identifier // list of parameters
	Body       *BlockStatement
}

func (ps *PurrStatement) statementNode()       {}
func (ps *PurrStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PurrStatement) String() string {
	var out strings.Builder
	out.WriteString(ps.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ps.Name.String())
	out.WriteString("(")

	for i, param := range ps.Parameters {
		out.WriteString(param.String())
		if i < len(ps.Parameters)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(") ")
	out.WriteString(ps.Body.String())

	return out.String()
}

// Chase Statement (`chase (condition) { body }` - while loop)
type ChaseStatement struct {
	Token     token.Token // the CHASE token
	Condition Expression
	Body      *BlockStatement
}

func (cs *ChaseStatement) statementNode()       {}
func (cs *ChaseStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ChaseStatement) String() string {
	var out strings.Builder
	out.WriteString("chase (")
	out.WriteString(cs.Condition.String())
	out.WriteString(") ")
	out.WriteString(cs.Body.String())
	return out.String()
}

// Break Statement (`claw;`)
type ClawStatement struct {
	Token token.Token // the CLAW token
}

func (cs *ClawStatement) statementNode()       {}
func (cs *ClawStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ClawStatement) String() string {
	return cs.Token.Literal + ";"
}

// Continue Statement (`stretch;`)
type StretchStatement struct {
	Token token.Token // the STRETCH token
}

func (ss *StretchStatement) statementNode()       {}
func (ss *StretchStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StretchStatement) String() string {
	return ss.Token.Literal + ";"
}
