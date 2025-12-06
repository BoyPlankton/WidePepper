// Copyright (c) 2025 BoyPlankton. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"widepepper/lexer"
	"widepepper/parser"
)

const version = "1.0.0"

func main() {
	// Define command-line flags
	showHelp := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version")
	showTokens := flag.Bool("tokens", false, "Show tokens instead of AST")
	readStdin := flag.Bool("stdin", false, "Read script from stdin")

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("WidePepper v%s\n", version)
		os.Exit(0)
	}

	var input string
	var err error

	if *readStdin {
		// Read from stdin
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(bytes)
	} else {
		// Get script file from arguments
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No script file provided\n")
			fmt.Fprintf(os.Stderr, "Usage: widepepper <script.wp> [options]\n")
			fmt.Fprintf(os.Stderr, "Use -help for more information\n")
			os.Exit(1)
		}

		scriptFile := args[0]
		input, err = readScriptFile(scriptFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading script file: %v\n", err)
			os.Exit(1)
		}
	}

	// Process the script
	if *showTokens {
		tokenizeScript(input)
	} else {
		parseScript(input)
	}
}

// readScriptFile reads a WidePepper script file (.wp)
func readScriptFile(filename string) (string, error) {
	// Validate file extension
	if filepath.Ext(filename) != ".wp" {
		fmt.Fprintf(os.Stderr, "Warning: Script file should have .wp extension\n")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// tokenizeScript shows the tokens produced by the lexer
func tokenizeScript(input string) {
	fmt.Println("--- WidePepper Lexer Tokenizing ---")
	fmt.Println()

	l := lexer.NewLexer(input)
	tokenCount := 0

	for tok := l.NextToken(); tok.Type != "EOF"; tok = l.NextToken() {
		fmt.Printf("%-12s | %q\n", tok.Type, tok.Literal)
		tokenCount++
	}

	fmt.Printf("\nTotal tokens: %d\n", tokenCount)
}

// parseScript parses a WidePepper script and produces an AST
func parseScript(input string) {
	fmt.Println("--- WidePepper Parser (AST Generation) ---")
	fmt.Println()

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if err := p.Errors(); err != nil {
		fmt.Fprintf(os.Stderr, "Parser Error:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %d statement(s):\n\n", len(program.Statements))
	for i, stmt := range program.Statements {
		fmt.Printf("[%d] %s\n", i+1, stmt.String())
	}
}

// printHelp displays usage information
func printHelp() {
	help := `WidePepper v` + version + ` - A Cat-Themed Scripting Language

Usage: widepepper [OPTIONS] [SCRIPT.wp]

Options:
  -help       Show this help message
  -version    Show version information
  -tokens     Display lexer tokens instead of parsed AST
  -stdin      Read script from standard input instead of a file

Examples:
  widepepper script.wp              # Parse and display AST
  widepepper -tokens script.wp      # Display tokens
  cat script.wp | widepepper -stdin # Read from pipe
  widepepper -help                  # Show this help

File Format:
  WidePepper scripts should use the .wp file extension.

For more information and syntax documentation, visit:
  https://github.com/BoyPlankton/WidePepper
`
	fmt.Print(help)
}
