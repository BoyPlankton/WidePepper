// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package parser implements a Pratt parser for the WidePepper language.
// It converts a stream of tokens from the lexer into an Abstract Syntax Tree (AST).
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"widepepper/ast"
	"widepepper/token"
)

// ParseError represents a parsing error with context information.
type ParseError struct {
	Token   token.Token
	Message string
}

// Error implements the error interface for ParseError.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %q: %s", e.Token.Literal, e.Message)
}

// ParseErrors is a collection of parsing errors.
type ParseErrors struct {
	errs []error
}

// Error implements the error interface for ParseErrors.
func (pe *ParseErrors) Error() string {
	if len(pe.errs) == 0 {
		return "no parse errors"
	}
	var messages []string
	for _, err := range pe.errs {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Add appends an error to the collection.
func (pe *ParseErrors) Add(err error) {
	if err != nil {
		pe.errs = append(pe.errs, err)
	}
}

// Len returns the number of errors in the collection.
func (pe *ParseErrors) Len() int {
	return len(pe.errs)
}

// All returns all errors as a slice.
func (pe *ParseErrors) All() []error {
	return pe.errs
}

// Precedence definition (Pratt Parsing)
const (
	_ int = iota
	LOWEST
	OR_AND      // ||, &&
	EQUALS      // ==, !=
	LESSGREATER // >, <, >=, <=
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -X, !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.OR:       OR_AND,
	token.AND:      OR_AND,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL, // For function calls
}

// Lexer interface for dependency injection
type Lexer interface {
	NextToken() token.Token
}

// Parser holds the lexer, the current token, and the lookahead token.
type Parser struct {
	l      Lexer
	errors *ParseErrors

	curToken  token.Token
	peekToken token.Token

	// Maps for registering prefix and infix parsing functions
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// NewParser creates a new Parser instance.
func NewParser(l Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: &ParseErrors{},
	}

	// Initialize maps
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// Register Prefix Parsing functions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STR, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)    // !
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)   // -
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression) // (expression)
	p.registerPrefix(token.CATNIP, p.parseBooleanLiteral)
	p.registerPrefix(token.MICE, p.parseBooleanLiteral)

	// Register Infix Parsing functions
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression) // For function calls

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// --- Parsing Helper Functions ---

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	err := &ParseError{
		Token:   p.peekToken,
		Message: fmt.Sprintf("expected %s, got %s", t, p.peekToken.Type),
	}
	p.errors.Add(err)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	err := &ParseError{
		Token:   p.curToken,
		Message: fmt.Sprintf("no prefix parsing function for %s", t),
	}
	p.errors.Add(err)
}

// --- Expression Parsing Functions ---

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		parseErr := &ParseError{
			Token:   p.curToken,
			Message: fmt.Sprintf("could not parse %q as float: %v", p.curToken.Literal, err),
		}
		p.errors.Add(parseErr)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.CATNIP)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // Consume '('

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseCallExpression handles function calls
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken, // LPAREN token
		Function: function,
	}

	exp.Arguments = p.parseCallArguments()

	return exp
}

// parseCallArguments parses the argument list inside parentheses
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	p.nextToken() // Move past '('

	// Handle empty argument list
	if p.curTokenIs(token.RPAREN) {
		return args
	}

	// Parse first argument
	args = append(args, p.parseExpression(LOWEST))

	// Parse remaining arguments
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Move to ','
		p.nextToken() // Move past ','

		args = append(args, p.parseExpression(LOWEST))
	}

	// Expect closing paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// --- Core Parsing Logic ---

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// Continue parsing infix expressions as long as the next token has a higher precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// --- Statement Parsing ---

// ParseProgram iterates through all statements in the source code.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement handles the different kinds of statements.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.YARN:
		return p.parseYarnStatement()
	case token.POUNCE:
		return p.parsePounceStatement()
	case token.TREAT:
		return p.parseReturnStatement()
	case token.PURR:
		return p.parsePurrStatement()
	case token.CHASE:
		return p.parseChaseStatement()
	case token.CLAW:
		return &ast.ClawStatement{Token: p.curToken}
	case token.STRETCH:
		return &ast.StretchStatement{Token: p.curToken}
	default:
		return nil
	}
}

// parseYarnStatement parses the `yarn variable = expression;` structure.
func (p *Parser) parseYarnStatement() *ast.YarnStatement {
	stmt := &ast.YarnStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken() // Consume '='

	// Parse the expression using LOWEST precedence
	stmt.Value = p.parseExpression(LOWEST)

	// Check for optional semicolon termination
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parses a block of statements enclosed in braces `{ ... }`
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // Consume '{'

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parsePounceStatement parses the `pounce (condition) { block } [hiss { block }]` structure
func (p *Parser) parsePounceStatement() *ast.PounceStatement {
	stmt := &ast.PounceStatement{Token: p.curToken}

	// Expect opening paren
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // Consume '('

	// Parse the condition expression
	stmt.Condition = p.parseExpression(LOWEST)

	// Expect closing paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Expect opening brace for consequence block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the consequence (main if block)
	stmt.Consequence = p.parseBlockStatement()

	// After parseBlockStatement, curToken is at RBRACE, move to next token
	p.nextToken()

	// Check for optional hiss (else) block
	if p.curTokenIs(token.HISS) {
		// Expect opening brace for alternative block
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		// Parse the alternative (else block)
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseReturnStatement parses the `treat [expression];` structure
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // Move past 'treat'

	// Check if there's no expression (treat followed by semicolon or end of block)
	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.RBRACE) {
		// No return value
		if p.curTokenIs(token.SEMICOLON) {
			// Already positioned at semicolon, no need to advance
		}
		return stmt
	}

	// Parse the return expression using LOWEST precedence
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// Check for optional semicolon termination
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parsePurrStatement parses the `purr function_name(param1, param2, ...) { body }` structure
func (p *Parser) parsePurrStatement() *ast.PurrStatement {
	stmt := &ast.PurrStatement{Token: p.curToken}

	// Expect function name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect opening paren
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse parameters
	stmt.Parameters = p.parseFunctionParameters()

	// Expect opening brace for function body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the function body
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parses the parameter list `param1, param2, ...` inside parentheses
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	p.nextToken() // Move past '('

	// Handle empty parameter list
	if p.curTokenIs(token.RPAREN) {
		return identifiers
	}

	// Parse first parameter
	if p.curTokenIs(token.IDENT) {
		identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	} else {
		p.peekError(token.IDENT)
	}

	// Parse remaining parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Move to ','
		p.nextToken() // Move past ','

		if !p.curTokenIs(token.IDENT) {
			p.peekError(token.IDENT)
			return nil
		}

		identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	// Expect closing paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseChaseStatement parses the `chase (condition) { body }` structure (while loop)
func (p *Parser) parseChaseStatement() *ast.ChaseStatement {
	stmt := &ast.ChaseStatement{Token: p.curToken}

	// Expect opening paren
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // Consume '('

	// Parse the condition expression
	stmt.Condition = p.parseExpression(LOWEST)

	// Expect closing paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Expect opening brace for loop body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the loop body
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// Errors returns the error collection, or nil if there are no errors.
// This implements idiomatic Go error handling.
func (p *Parser) Errors() error {
	if p.errors.Len() == 0 {
		return nil
	}
	return p.errors
}

// HasErrors returns true if any parsing errors occurred.
func (p *Parser) HasErrors() bool {
	return p.errors.Len() > 0
}
