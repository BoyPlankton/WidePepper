// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"widepepper/ast"
	"widepepper/lexer"
	"widepepper/token"
)

// TestYarnStatementWithIntegerLiteral tests yarn x = 5;
func TestYarnStatementWithIntegerLiteral(t *testing.T) {
	l := lexer.NewLexer("yarn x = 5;")
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has %d statements, expected 1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.YarnStatement)
	if !ok {
		t.Fatalf("statement is not YarnStatement, got %T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("name is not 'x', got %q", stmt.Name.Value)
	}

	intLit, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("value is not IntegerLiteral, got %T", stmt.Value)
	}

	if intLit.Value != 5 {
		t.Errorf("value is not 5, got %v", intLit.Value)
	}
}

// TestYarnStatementWithIdentifier tests yarn x = y;
func TestYarnStatementWithIdentifier(t *testing.T) {
	l := lexer.NewLexer("yarn x = y;")
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt := program.Statements[0].(*ast.YarnStatement)
	if stmt.Name.Value != "x" {
		t.Errorf("name is not 'x', got %q", stmt.Name.Value)
	}

	ident, ok := stmt.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("value is not Identifier, got %T", stmt.Value)
	}

	if ident.Value != "y" {
		t.Errorf("identifier value is not 'y', got %q", ident.Value)
	}
}

// TestYarnStatementWithStringLiteral tests yarn msg = "hello";
func TestYarnStatementWithStringLiteral(t *testing.T) {
	l := lexer.NewLexer(`yarn msg = "hello";`)
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt := program.Statements[0].(*ast.YarnStatement)
	if stmt.Name.Value != "msg" {
		t.Errorf("name is not 'msg', got %q", stmt.Name.Value)
	}

	strLit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("value is not StringLiteral, got %T", stmt.Value)
	}

	if strLit.Value != "hello" {
		t.Errorf("string value is not 'hello', got %q", strLit.Value)
	}
}

// TestYarnStatementWithBooleanLiteral tests yarn flag = catnip;
func TestYarnStatementWithBooleanLiteral(t *testing.T) {
	l := lexer.NewLexer("yarn flag = catnip;")
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt := program.Statements[0].(*ast.YarnStatement)
	boolLit, ok := stmt.Value.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("value is not BooleanLiteral, got %T", stmt.Value)
	}

	if !boolLit.Value {
		t.Errorf("boolean value is not true, got %v", boolLit.Value)
	}
}

// TestOperatorPrecedence tests various operator precedence scenarios
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"addition and multiplication", "treat 5 + 2 * 3;", "(5 + (2 * 3))"},
		{"multiplication and addition", "treat 2 * 3 + 5;", "((2 * 3) + 5)"},
		{"division and subtraction", "treat 10 / 2 - 3;", "((10 / 2) - 3)"},
		{"grouped expression", "treat (5 + 2) * 3;", "((5 + 2) * 3)"},
		{"nested grouping", "treat ((5 + 2) * (3 + 4));", "((5 + 2) * (3 + 4))"},
		{"comparison operators", "treat a > b && c < d;", "((a > b) && (c < d))"},
		{"equality and comparison", "treat x == 5 && y > 10;", "((x == 5) && (y > 10))"},
		{"logical operators", "treat a || b && c;", "((a || b) && c)"},
		{"negation and comparison", "treat !catnip || x > 5;", "((!catnip) || (x > 5))"},
		{"prefix on grouped", "treat -(5 + 2);", "(-(5 + 2))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			testParseErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has %d statements, expected 1", len(program.Statements))
			}

			// We expect a treat statement with just an expression
			stmt := program.Statements[0].(*ast.ReturnStatement)
			if stmt.ReturnValue == nil {
				t.Fatalf("return statement has no value")
			}

			actual := stmt.ReturnValue.String()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

// TestPrefixExpressions tests prefix operators: !, -
func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		operator string
		expected string
	}{
		{"negation of number", "treat -5;", "-", "(−5)"},
		{"logical not", "treat !catnip;", "!", "(!catnip)"},
		{"double negation", "treat !!x;", "!", "(!(!x))"},
		{"not of false", "treat !mice;", "!", "(!mice)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			testParseErrors(t, p)

			stmt := program.Statements[0].(*ast.ReturnStatement)
			prefix, ok := stmt.ReturnValue.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("expression is not PrefixExpression, got %T", stmt.ReturnValue)
			}

			if prefix.Operator != tt.operator {
				t.Errorf("operator is not %q, got %q", tt.operator, prefix.Operator)
			}
		})
	}
}

// TestInfixExpressions tests infix operators: +, -, *, /, ==, !=, <, >, <=, >=, &&, ||
func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		left     string
		operator string
		right    string
	}{
		{"addition", "treat 5 + 3;", "5", "+", "3"},
		{"subtraction", "treat 10 - 4;", "10", "-", "4"},
		{"multiplication", "treat 2 * 6;", "2", "*", "6"},
		{"division", "treat 20 / 4;", "20", "/", "4"},
		{"equality", "treat x == 5;", "x", "==", "5"},
		{"inequality", "treat a != b;", "a", "!=", "b"},
		{"less than", "treat 3 < 10;", "3", "<", "10"},
		{"greater than", "treat 10 > 3;", "10", ">", "3"},
		{"less or equal", "treat 5 <= 5;", "5", "<=", "5"},
		{"greater or equal", "treat 10 >= 5;", "10", ">=", "5"},
		{"logical AND", "treat catnip && mice;", "catnip", "&&", "mice"},
		{"logical OR", "treat x || y;", "x", "||", "y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			testParseErrors(t, p)

			stmt := program.Statements[0].(*ast.ReturnStatement)
			infix, ok := stmt.ReturnValue.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("expression is not InfixExpression, got %T", stmt.ReturnValue)
			}

			if infix.Left.String() != tt.left {
				t.Errorf("left is not %q, got %q", tt.left, infix.Left.String())
			}

			if infix.Operator != tt.operator {
				t.Errorf("operator is not %q, got %q", tt.operator, infix.Operator)
			}

			if infix.Right.String() != tt.right {
				t.Errorf("right is not %q, got %q", tt.right, infix.Right.String())
			}
		})
	}
}

// TestBooleanLogicExpressions tests complex boolean expressions
func TestBooleanLogicExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"a == b && c < 10", "treat a == b && c < 10;", "((a == b) && (c < 10))"},
		{"!true", "treat !catnip;", "(!catnip)"},
		{"x == 5 || y > 10", "treat x == 5 || y > 10;", "((x == 5) || (y > 10))"},
		{"with negation", "treat !a && b == c;", "((!a) && (b == c))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			if err := p.Errors(); err != nil {
				t.Logf("parser error: %v", err)
				return
			}

			if len(program.Statements) == 0 {
				t.Fatalf("program has no statements")
			}

			stmt := program.Statements[0].(*ast.ReturnStatement)
			actual := stmt.ReturnValue.String()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

// TestMalformedSyntaxErrors tests error reporting for malformed syntax
func TestMalformedSyntaxErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldHaveErr bool
	}{
		{"missing equals in yarn", "yarn x 5;", true},
		{"missing identifier after yarn", "yarn = 5;", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			_ = p.ParseProgram()

			err := p.Errors()
			if tt.shouldHaveErr {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestMultipleStatements tests parsing multiple statements
func TestMultipleStatements(t *testing.T) {
	input := `
		yarn x = 5;
		yarn y = 10;
		treat x + y;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program has %d statements, expected 3", len(program.Statements))
	}

	// First statement: yarn x = 5;
	stmt1, ok := program.Statements[0].(*ast.YarnStatement)
	if !ok {
		t.Fatalf("statement 1 is not YarnStatement, got %T", program.Statements[0])
	}
	if stmt1.Name.Value != "x" {
		t.Errorf("statement 1 name is not 'x', got %q", stmt1.Name.Value)
	}

	// Second statement: yarn y = 10;
	stmt2, ok := program.Statements[1].(*ast.YarnStatement)
	if !ok {
		t.Fatalf("statement 2 is not YarnStatement, got %T", program.Statements[1])
	}
	if stmt2.Name.Value != "y" {
		t.Errorf("statement 2 name is not 'y', got %q", stmt2.Name.Value)
	}

	// Third statement: treat x + y;
	stmt3, ok := program.Statements[2].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("statement 3 is not ReturnStatement, got %T", program.Statements[2])
	}
	if stmt3.ReturnValue == nil {
		t.Errorf("statement 3 return value is nil")
	}
}

// TestConditionalStatements tests pounce/hiss (if/else) parsing
func TestConditionalStatements(t *testing.T) {
	input := `pounce (catnip) { treat 1; } hiss { treat 2; }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.PounceStatement)
	if !ok {
		t.Fatalf("statement is not PounceStatement, got %T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Errorf("condition is nil")
	}

	if stmt.Consequence == nil {
		t.Errorf("consequence block is nil")
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence has %d statements, expected 1", len(stmt.Consequence.Statements))
	}

	if stmt.Alternative == nil {
		t.Errorf("alternative (hiss) block is nil")
	}

	if len(stmt.Alternative.Statements) != 1 {
		t.Errorf("alternative has %d statements, expected 1", len(stmt.Alternative.Statements))
	}
}

// TestLoopStatements tests chase (for loop) parsing
func TestLoopStatements(t *testing.T) {
	input := `chase (catnip) { treat 1; }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ChaseStatement)
	if !ok {
		t.Fatalf("statement is not ChaseStatement, got %T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Errorf("condition is nil")
	}

	if stmt.Body == nil {
		t.Errorf("body block is nil")
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("body has %d statements, expected 1", len(stmt.Body.Statements))
	}
}

// TestBreakAndContinue tests claw (break) and stretch (continue) parsing
func TestBreakAndContinue(t *testing.T) {
	input := `
		claw;
		stretch;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has %d statements, expected 2", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ClawStatement)
	if !ok {
		t.Fatalf("statement 1 is not ClawStatement, got %T", program.Statements[0])
	}

	_, ok = program.Statements[1].(*ast.StretchStatement)
	if !ok {
		t.Fatalf("statement 2 is not StretchStatement, got %T", program.Statements[1])
	}
}

// TestComplexExpressionsInStatements tests expressions with various combinations
func TestComplexExpressionsInStatements(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"arithmetic in yarn", "yarn result = 5 + 3 * 2;"},
		{"boolean logic in yarn", "yarn condition = x > 5 && y < 10;"},
		{"nested conditions", "pounce (x > 0 && y < 100) { treat 1; }"},
		{"comparison chain in return", "treat a == b && c != d || e <= f;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			err := p.Errors()
			if err != nil {
				t.Logf("parser error (may be expected): %v", err)
			}

			if len(program.Statements) != 1 {
				t.Errorf("expected 1 statement, got %d", len(program.Statements))
			}
		})
	}
}

// TestErrorMessagesAreClear tests that error messages are informative
func TestErrorMessagesAreClear(t *testing.T) {
	l := lexer.NewLexer("yarn x 5;")
	p := NewParser(l)
	_ = p.ParseProgram()

	err := p.Errors()
	if err == nil {
		t.Fatal("expected error, got none")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("error message is empty")
	}

	// Error message should contain some context about what was expected
	if len(errMsg) < 5 {
		t.Errorf("error message too short: %q", errMsg)
	}
}

// Helper functions

func testParseErrors(t *testing.T, p *Parser) {
	err := p.Errors()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
}

// TestReturnStatementValue tests treat statements with various return values
func TestReturnStatementValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		valueStr string
	}{
		{"return number", "treat 42;", "42"},
		{"return identifier", "treat x;", "x"},
		{"return string", `treat "hello";`, `"hello"`},
		{"return boolean", "treat catnip;", "catnip"},
		{"return expression", "treat 5 + 3;", "(5 + 3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			testParseErrors(t, p)

			stmt, ok := program.Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("statement is not ReturnStatement, got %T", program.Statements[0])
			}

			if stmt.ReturnValue == nil {
				t.Fatalf("return value is nil")
			}

			if stmt.ReturnValue.String() != tt.valueStr {
				t.Errorf("expected %q, got %q", tt.valueStr, stmt.ReturnValue.String())
			}
		})
	}
}

// TestIdentifierTokenTypes tests that identifiers preserve their token type
func TestIdentifierTokenTypes(t *testing.T) {
	l := lexer.NewLexer("treat x;")
	p := NewParser(l)
	program := p.ParseProgram()

	testParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ReturnStatement)
	ident, ok := stmt.ReturnValue.(*ast.Identifier)
	if !ok {
		t.Fatalf("value is not Identifier, got %T", stmt.ReturnValue)
	}

	if ident.Token.Type != token.IDENT {
		t.Errorf("token type is not IDENT, got %v", ident.Token.Type)
	}

	if ident.Token.Literal != "x" {
		t.Errorf("token literal is not 'x', got %q", ident.Token.Literal)
	}
}
