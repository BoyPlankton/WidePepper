// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package token defines all token types and constants for the WidePepper language.
// It provides the fundamental lexical units used by the lexer and parser.
package token

// TokenType is a string that represents the type of a token.
// It's used to classify different lexical elements of the WidePepper language.
type TokenType string

// String returns the string representation of the token type.
func (t TokenType) String() string {
	return string(t)
}

// Tokens represent all lexical elements of the WidePepper language
const (
	ILLEGAL   TokenType = "ILLEGAL"
	EOF       TokenType = "EOF"
	IDENT     TokenType = "IDENT" // identifiers like x, y, my_variable
	INT       TokenType = "INT"   // integer or float literals (e.g., 5, 3.14)
	STR       TokenType = "STR"   // string literals (e.g., "hello")
	ASSIGN    TokenType = "="
	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	BANG      TokenType = "!"
	ASTERISK  TokenType = "*"
	SLASH     TokenType = "/"
	MODULO    TokenType = "%"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	DOT       TokenType = "."

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	LBRACK TokenType = "["
	RBRACK TokenType = "]"

	// Keywords (WidePepper specific)
	YARN    TokenType = "yarn"    // var declaration
	PURR    TokenType = "purr"    // function definition
	POUNCE  TokenType = "pounce"  // if statement
	HISS    TokenType = "hiss"    // else statement
	TREAT   TokenType = "treat"   // return statement
	MEOW    TokenType = "meow"    // print statement
	CATNIP  TokenType = "catnip"  // true
	MICE    TokenType = "mice"    // false
	CHASE   TokenType = "chase"   // while loop
	CLAW    TokenType = "claw"    // break
	STRETCH TokenType = "stretch" // continue

	// Operators
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LTE    TokenType = "<="
	GTE    TokenType = ">="
	AND    TokenType = "&&"
	OR     TokenType = "||"
)

// Token is a single token from the lexer.
type Token struct {
	Type    TokenType
	Literal string
}

// keywords is the map of reserved words and their token types.
var keywords = map[string]TokenType{
	"yarn":    YARN,
	"purr":    PURR,
	"pounce":  POUNCE,
	"hiss":    HISS,
	"treat":   TREAT,
	"meow":    MEOW,
	"catnip":  CATNIP,
	"mice":    MICE,
	"chase":   CHASE,
	"claw":    CLAW,
	"stretch": STRETCH,
}

// LookupIdent returns the token type of an identifier or keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
