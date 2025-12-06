// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package ast

import (
	"testing"

	"widepepper/token"
)

func TestIdentifierTokenLiteral(t *testing.T) {
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "foobar"},
		Value: "foobar",
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %q. got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestIdentifierString(t *testing.T) {
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "myVar"},
		Value: "myVar",
	}

	if ident.String() != "myVar" {
		t.Errorf("ident.String() not %q. got=%q", "myVar", ident.String())
	}
}

func TestIntegerLiteralString(t *testing.T) {
	lit := &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "5"},
		Value: 5,
	}

	if lit.String() != "5" {
		t.Errorf("lit.String() not %q. got=%q", "5", lit.String())
	}
}

func TestStringLiteralString(t *testing.T) {
	lit := &StringLiteral{
		Token: token.Token{Type: token.STR, Literal: `"hello"`},
		Value: "hello",
	}

	if lit.String() != `"hello"` {
		t.Errorf("lit.String() not %q. got=%q", `"hello"`, lit.String())
	}
}

func TestBooleanLiteralString(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		literal  string
		expected string
	}{
		{"true", true, "catnip", "catnip"},
		{"false", false, "mice", "mice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := &BooleanLiteral{
				Token: token.Token{Type: token.CATNIP, Literal: tt.literal},
				Value: tt.value,
			}

			if lit.String() != tt.expected {
				t.Errorf("lit.String() not %q. got=%q", tt.expected, lit.String())
			}
		})
	}
}

func TestPrefixExpressionString(t *testing.T) {
	expr := &PrefixExpression{
		Token:    token.Token{Type: token.BANG, Literal: "!"},
		Operator: "!",
		Right: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x"},
			Value: "x",
		},
	}

	if expr.String() != "(!x)" {
		t.Errorf("expr.String() not %q. got=%q", "(!x)", expr.String())
	}
}

func TestInfixExpressionString(t *testing.T) {
	expr := &InfixExpression{
		Token: token.Token{Type: token.PLUS, Literal: "+"},
		Left: &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "5"},
			Value: 5,
		},
		Operator: "+",
		Right: &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "10"},
			Value: 10,
		},
	}

	if expr.String() != "(5 + 10)" {
		t.Errorf("expr.String() not %q. got=%q", "(5 + 10)", expr.String())
	}
}

func TestCallExpressionString(t *testing.T) {
	expr := &CallExpression{
		Token: token.Token{Type: token.LPAREN, Literal: "("},
		Function: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "add"},
			Value: "add",
		},
		Arguments: []Expression{
			&IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "1"},
				Value: 1,
			},
			&IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "2"},
				Value: 2,
			},
		},
	}

	if expr.String() != "add(1, 2)" {
		t.Errorf("expr.String() not %q. got=%q", "add(1, 2)", expr.String())
	}
}

func TestYarnStatementString(t *testing.T) {
	stmt := &YarnStatement{
		Token: token.Token{Type: token.YARN, Literal: "yarn"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "myVar"},
			Value: "myVar",
		},
		Value: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
			Value: "anotherVar",
		},
	}

	if stmt.String() != "yarn myVar = anotherVar;" {
		t.Errorf("stmt.String() not %q. got=%q", "yarn myVar = anotherVar;", stmt.String())
	}
}

func TestBlockStatementString(t *testing.T) {
	stmt := &BlockStatement{
		Token: token.Token{Type: token.LBRACE, Literal: "{"},
		Statements: []Statement{
			&ReturnStatement{
				Token: token.Token{Type: token.TREAT, Literal: "treat"},
				ReturnValue: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
		},
	}

	result := stmt.String()
	if result == "" {
		t.Error("BlockStatement.String() returned empty string")
	}
}

func TestReturnStatementString(t *testing.T) {
	stmt := &ReturnStatement{
		Token: token.Token{Type: token.TREAT, Literal: "treat"},
		ReturnValue: &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "42"},
			Value: 42,
		},
	}

	if stmt.String() != "treat 42;" {
		t.Errorf("stmt.String() not %q. got=%q", "treat 42;", stmt.String())
	}
}

func TestClawStatementString(t *testing.T) {
	stmt := &ClawStatement{
		Token: token.Token{Type: token.CLAW, Literal: "claw"},
	}

	if stmt.String() != "claw;" {
		t.Errorf("stmt.String() not %q. got=%q", "claw;", stmt.String())
	}
}

func TestStretchStatementString(t *testing.T) {
	stmt := &StretchStatement{
		Token: token.Token{Type: token.STRETCH, Literal: "stretch"},
	}

	if stmt.String() != "stretch;" {
		t.Errorf("stmt.String() not %q. got=%q", "stretch;", stmt.String())
	}
}

func TestProgramString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&YarnStatement{
				Token: token.Token{Type: token.YARN, Literal: "yarn"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
		},
	}

	if program.String() == "" {
		t.Error("Program.String() returned empty string")
	}
}
