// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package lexer

import (
	"testing"
	"widepepper/token"
)

// TestNextToken tests the basic tokenization of the lexer.
func TestNextToken(t *testing.T) {
	input := `yarn x = 5;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YARN, "yarn"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestKeywords tests that all WidePepper keywords are recognized.
func TestKeywords(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"yarn", token.YARN, "yarn"},
		{"purr", token.PURR, "purr"},
		{"pounce", token.POUNCE, "pounce"},
		{"hiss", token.HISS, "hiss"},
		{"treat", token.TREAT, "treat"},
		{"meow", token.MEOW, "meow"},
		{"catnip", token.CATNIP, "catnip"},
		{"mice", token.MICE, "mice"},
		{"chase", token.CHASE, "chase"},
		{"claw", token.CLAW, "claw"},
		{"stretch", token.STRETCH, "stretch"},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("keyword %q: expected type %q, got %q",
				tt.input, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("keyword %q: expected literal %q, got %q",
				tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestIdentifiers tests that identifiers are correctly tokenized.
func TestIdentifiers(t *testing.T) {
	input := `my_var foobar _leading trailing_`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "my_var"},
		{token.IDENT, "foobar"},
		{token.IDENT, "_leading"},
		{token.IDENT, "trailing_"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestIntegerLiterals tests integer and float literal tokenization.
func TestIntegerLiterals(t *testing.T) {
	input := `5 42 3.14 0 999.999`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "5"},
		{token.INT, "42"},
		{token.INT, "3.14"},
		{token.INT, "0"},
		{token.INT, "999.999"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestStringLiterals tests string literal tokenization.
func TestStringLiterals(t *testing.T) {
	input := `"hello" "world" "foo bar" ""`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STR, "hello"},
		{token.STR, "world"},
		{token.STR, "foo bar"},
		{token.STR, ""},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestSingleCharOperators tests single-character operators.
func TestSingleCharOperators(t *testing.T) {
	input := `+ - * / % = ! < > , ; . ( ) { } [ ]`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.MODULO, "%"},
		{token.ASSIGN, "="},
		{token.BANG, "!"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.DOT, "."},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACK, "["},
		{token.RBRACK, "]"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestMultiCharOperators tests multi-character operators.
func TestMultiCharOperators(t *testing.T) {
	input := `== != <= >= && ||`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestSingleLineComments tests that single-line comments are skipped.
func TestSingleLineComments(t *testing.T) {
	input := `yarn x = 5; // this is a comment
yarn y = 10;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YARN, "yarn"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.YARN, "yarn"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestMultiLineComments tests that multi-line comments are skipped.
func TestMultiLineComments(t *testing.T) {
	input := `yarn x = 5; /* multi
line
comment */ yarn y = 10;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YARN, "yarn"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.YARN, "yarn"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestWhitespaceSkipping tests that whitespace is correctly skipped.
func TestWhitespaceSkipping(t *testing.T) {
	input := "yarn   x\t=\n5  ;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YARN, "yarn"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestComplexExpression tests a more complex expression with mixed tokens.
func TestComplexExpression(t *testing.T) {
	input := `purr add(a, b) {
		treat a + b;
	}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PURR, "purr"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.TREAT, "treat"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestConditionalExpression tests tokenization of conditional statements.
func TestConditionalExpression(t *testing.T) {
	input := `pounce (age > 5 && name != "cat") {
		meow "adult";
	} hiss {
		meow "young";
	}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.POUNCE, "pounce"},
		{token.LPAREN, "("},
		{token.IDENT, "age"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "name"},
		{token.NOT_EQ, "!="},
		{token.STR, "cat"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.MEOW, "meow"},
		{token.STR, "adult"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.HISS, "hiss"},
		{token.LBRACE, "{"},
		{token.MEOW, "meow"},
		{token.STR, "young"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestLoopExpression tests tokenization of loop statements.
func TestLoopExpression(t *testing.T) {
	input := `chase (x < 10) {
		x = x + 1;
		pounce (x == 5) {
			claw;
		} hiss {
			stretch;
		}
	}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CHASE, "chase"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.POUNCE, "pounce"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.CLAW, "claw"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.HISS, "hiss"},
		{token.LBRACE, "{"},
		{token.STRETCH, "stretch"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestBooleanLiterals tests tokenization of boolean values.
func TestBooleanLiterals(t *testing.T) {
	input := `catnip mice true false`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CATNIP, "catnip"},
		{token.MICE, "mice"},
		{token.IDENT, "true"},
		{token.IDENT, "false"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestArithmeticOperations tests tokenization of arithmetic expressions.
func TestArithmeticOperations(t *testing.T) {
	input := `5 + 3 * 2 - 1 / 4 % 2`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.MINUS, "-"},
		{token.INT, "1"},
		{token.SLASH, "/"},
		{token.INT, "4"},
		{token.MODULO, "%"},
		{token.INT, "2"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestMixedComments tests both single-line and multi-line comments together.
func TestMixedComments(t *testing.T) {
	input := `yarn x = 5; // single line comment
/* multi-line
   comment here */
yarn y = 10; // another comment`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YARN, "yarn"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.YARN, "yarn"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestEmptyInput tests lexing an empty input.
func TestEmptyInput(t *testing.T) {
	input := ``

	l := NewLexer(input)
	tok := l.NextToken()

	if tok.Type != token.EOF {
		t.Fatalf("expected EOF, got %q", tok.Type)
	}
}

// TestOnlyWhitespace tests lexing only whitespace.
func TestOnlyWhitespace(t *testing.T) {
	input := `   
	
  `

	l := NewLexer(input)
	tok := l.NextToken()

	if tok.Type != token.EOF {
		t.Fatalf("expected EOF, got %q", tok.Type)
	}
}

// TestNegativeNumbers tests that negative numbers are properly tokenized.
func TestNegativeNumbers(t *testing.T) {
	input := `- 5 - 3.14`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MINUS, "-"},
		{token.INT, "5"},
		{token.MINUS, "-"},
		{token.INT, "3.14"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestStringContent tests strings with various content.
func TestStringContent(t *testing.T) {
	input := `"hello world" "123" "special!@#"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STR, "hello world"},
		{token.STR, "123"},
		{token.STR, "special!@#"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// BenchmarkLexer benchmarks lexer performance on a realistic input.
func BenchmarkLexer(b *testing.B) {
	input := `purr fibonacci(n) {
		pounce (n <= 1) {
			treat n;
		}
		yarn a = 0;
		yarn b = 1;
		chase (n > 1) {
			yarn temp = a + b;
			a = b;
			b = temp;
			n = n - 1;
		}
		treat b;
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := NewLexer(input)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		}
	}
}
