# **WidePepper Scripting Language**

WidePepper is a dynamically typed, embeddable scripting language with a cat-themed syntax. Its design philosophy centers around simplicity, readability, and ensuring users are always entertained by its feline reserved words.

![WidePepper "Spread It"](https://github.com/KrosTheProto/WidePepper-test/blob/main/widepepper.jpg?raw=true)

## **Table of Contents**

1. [Quick Start](#quick-start)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Examples](#examples)
5. [Language Specification](#language-specification)
6. [Development](#development)

---

## **Quick Start**

### Running a WidePepper Script

```bash
# Run a script file
widepepper script.wp

# Show tokens (for debugging)
widepepper -tokens script.wp

# Read from stdin
cat script.wp | widepepper -stdin

# Show help
widepepper -help
```

### Hello World

Create a file named `hello.wp`:

```widepepper
yarn greeting = "Hello, WidePepper!";
meow greeting;
```

Run it:

```bash
widepepper hello.wp
```

---

## **Installation**

### Prerequisites

- Go 1.16 or higher
- Git

### Building from Source

```bash
git clone https://github.com/BoyPlankton/WidePepper.git
cd WidePepper
make build
```

The compiled binary will be in `./bin/widepepper`.

### Installing Globally

```bash
make install
```

This installs the `widepepper` command to your `$GOPATH/bin` directory.

---

## **Usage**

### Command-Line Options

```
widepepper [OPTIONS] [SCRIPT.wp]

Options:
  -help       Show help message
  -version    Show version information
  -tokens     Display lexer tokens instead of parsed AST
  -stdin      Read script from standard input
```

### Examples

```bash
# Parse and show AST
widepepper script.wp

# Show tokens for debugging
widepepper -tokens script.wp

# Pipe script from another command
echo 'meow "test";' | widepepper -stdin

# Display version
widepepper -version
```

---

## **Examples**

WidePepper includes several example scripts in the `examples/` directory:

| Example | Description |
|---------|-------------|
| `01_hello_world.wp` | Simplest program - declare a variable and print it |
| `02_variables.wp` | Variables and arithmetic operators |
| `03_conditionals.wp` | If/else statements (pounce/hiss) |
| `04_loops.wp` | Looping with chase, continue (stretch), and break (claw) |
| `05_functions.wp` | Function definition and calling (purr) |
| `06_operators.wp` | Arithmetic, comparison, and logical operators |
| `07_complex_example.wp` | Complete program combining multiple features |

### Running Examples

```bash
# Run an example with make
make run-example EXAMPLE=01_hello_world
make run-example EXAMPLE=02_variables

# Or run directly
widepepper examples/01_hello_world.wp
widepepper examples/04_loops.wp

# Show tokens for an example
make run-tokens EXAMPLE=02_variables
```

---

## **Language Specification**

### **Lexical Elements**

#### **Reserved Words (Keywords)**

All reserved words in WidePepper are case-sensitive and cat-themed.

| Keyword | Purpose | Traditional Equivalent |
|---------|---------|----------------------|
| `yarn` | Variable declaration | `var`, `let` |
| `purr` | Function definition | `function`, `def` |
| `pounce` | Conditional (if) | `if` |
| `hiss` | Else branch | `else` |
| `chase` | Loop (while) | `while`, `for` |
| `treat` | Return value | `return` |
| `meow` | Output/Print | `print`, `log` |
| `catnip` | Boolean true | `true` |
| `mice` | Boolean false | `false` |
| `claw` | Break loop | `break` |
| `stretch` | Continue loop | `continue` |

#### **Operators**

| Symbol | Purpose |
|--------|---------|
| `=` | Assignment |
| `==`, `!=`, `<`, `>`, `<=`, `>=` | Comparison |
| `+`, `-`, `*`, `/` | Arithmetic |
| `&&`, `\|\|`, `!` | Logical (AND, OR, NOT) |
| `;` | Statement terminator |
| `{}` | Code blocks |
| `()` | Grouping, function calls |
| `[]` | Array indexing |
| `.` | Member access |

#### **Comments**

```widepepper
// Single-line comment (Whiskers style)

/* Multi-line comment (Nap time style)
   Perfect for long explanations
*/
```

### **Data Types**

WidePepper supports several core data types:

| Type | Description | Example |
|------|-------------|---------|
| **Number** | Integer or float | `42`, `3.14` |
| **String** | Text in double quotes | `"Hello, cat"` |
| **Boolean** | Logical value | `catnip` (true), `mice` (false) |
| **Array** | Ordered collection | `["toy", "ball", 3]` |
| **Null** | No value (reserved) | `EmptyBowl` |

### **Variables**

Variables are declared with the `yarn` keyword and are dynamically typed.

```widepepper
// Declaration and initialization
yarn cat_name = "Whiskers";
yarn age = 5;
yarn is_happy = catnip;

// Reassignment (yarn omitted after first declaration)
age = age + 1;

// Output
meow cat_name;
```

### **Functions**

Functions are defined with the `purr` keyword and use `treat` to return values.

```widepepper
// Function definition
purr greet(name) {
    meow "Hello, " + name + "!";
}

// Function with return value
purr add(a, b) {
    yarn sum = a + b;
    treat sum;
}

// Using functions
greet("Mittens");
yarn result = add(5, 3);
meow result;
```

### **Control Flow**

#### **Conditionals (pounce/hiss)**

```widepepper
yarn age = 7;

pounce (age < 3) {
    meow "Kitten!";
} hiss {
    meow "Adult cat!";
}

// Nested conditionals
pounce (age > 10) {
    meow "Senior kitty!";
} hiss pounce (age > 3) {
    meow "Adult cat!";
} hiss {
    meow "Young cat!";
}
```

#### **Loops (chase)**

```widepepper
// Loop with counter
yarn count = 0;
chase (count < 5) {
    meow count;
    count = count + 1;
}

// With continue (stretch)
yarn i = 0;
chase (i < 10) {
    i = i + 1;
    pounce (i == 5) {
        stretch;  // Skip this iteration
    }
    meow i;
}

// With break (claw)
yarn found = mice;
chase (found == mice) {
    pounce (some_condition) {
        found = catnip;
        claw;  // Exit loop
    }
}
```

### **Operators**

#### **Arithmetic**

```widepepper
yarn a = 10;
yarn b = 3;

meow a + b;      // 13
meow a - b;      // 7
meow a * b;      // 30
meow a / b;      // 3.333...
```

#### **Comparison**

```widepepper
yarn x = 5;
yarn y = 10;

meow x < y;      // catnip (true)
meow x > y;      // mice (false)
meow x == y;     // mice
meow x != y;     // catnip
```

#### **Logical**

```widepepper
yarn is_sunny = catnip;
yarn is_warm = catnip;

meow is_sunny && is_warm;    // catnip
meow is_sunny || is_warm;    // catnip
meow !is_sunny;              // mice
```

#### **Operator Precedence** (highest to lowest)

1. Grouping: `()`
2. Unary: `!`, `-`
3. Multiplication/Division: `*`, `/`
4. Addition/Subtraction: `+`, `-`
5. Comparison: `<`, `>`, `<=`, `>=`
6. Equality: `==`, `!=`
7. Logical AND: `&&`
8. Logical OR: `||`

### **Input/Output**

The `meow` keyword outputs values to the console.

```widepepper
meow "Hello";           // String
meow 42;                // Number
meow catnip;            // Boolean
meow 5 + 3;             // Expression
```

---

## **Development**

### **Build Targets**

The project includes a `Makefile` with convenient development targets:

```bash
# Build the project
make build

# Run all tests
make test

# Run tests with coverage
make test-cover

# Display coverage percentage
make coverage

# Run a specific example
make run-example EXAMPLE=02_variables

# Show tokens for debugging
make run-tokens EXAMPLE=02_variables

# Format code
make fmt

# Run linter
make vet

# Clean build artifacts
make clean

# View all available targets
make help
```

### **Project Structure**

```
WidePepper/
├── ast/              # Abstract Syntax Tree definitions
├── lexer/            # Lexical analyzer
├── parser/           # Parser implementation
├── token/            # Token definitions
├── examples/         # Example WidePepper scripts (.wp files)
├── main.go           # CLI entry point
├── go.mod            # Go module definition
├── Makefile          # Build and development targets
└── README.md         # This file
```

### **Testing**

Run tests with coverage:

```bash
make test-cover
```

View detailed coverage:

```bash
make coverage
```

Generate HTML coverage report:

```bash
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

### **Current Status**

- ✅ Lexer: 95.8% coverage, 21 comprehensive tests
- ✅ Parser: 69.4% coverage, 17 comprehensive tests
- ✅ AST: 57.6% coverage, 13+ comprehensive tests
- ✅ Total: 50+ passing tests across all packages

### **Architecture**

The compiler is organized into modular packages:

1. **Token Package** (`token/`): Defines token types used by the lexer
2. **Lexer Package** (`lexer/`): Tokenizes source code into a stream of tokens
3. **AST Package** (`ast/`): Defines Abstract Syntax Tree node types
4. **Parser Package** (`parser/`): Parses tokens into an AST
5. **Main Package** (`main.go`): CLI and script execution

---

## **Contributing**

Contributions are welcome! Areas for enhancement include:

- Full AST evaluation/interpretation
- Additional built-in functions
- Better error messages
- Standard library (strings, math, etc.)
- REPL mode
- Optimization passes

---

## **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy scripting with WidePepper! 🐱**  
