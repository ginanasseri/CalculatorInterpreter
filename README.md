# Calculator

**Author:** Gina Nasseri

A calculator interpreter that handles addition, subtraction, multiplication, and division. The alphabet accepted is `{'+', '-', '*', '\',' (' , ')'}` along with integer strings of any length. The input is parsed according to context-free grammar rules, constructing an abstract syntax tree (AST), which is then interpreted to evaluate expressions. 

Usage: run `go run main.go` to start the calculator

The packages in this calculator:
- `token`: defines the token type
- `lexer`: creates tokens from the input 
- `parser`: checks token syntax and builds AST
- `ast`: contains the ASTNode and ASTVisitor interfaces and node methods
- `interpreter`: traverses the AST provided by the parser and calculates the result 
- `nestingstack`: used to ensure parentheses are balanced. 

Also included: `main_test.go` that extensively tests for syntax error cases and ensures order of operations is followed. Use `go test` to run.
