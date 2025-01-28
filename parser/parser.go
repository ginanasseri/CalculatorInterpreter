/*
The p: checks syntax and builds AST for the interpreter to execute the Visitor pattern on to produce the 
result of the expression
*/

package parser

import (
    "calculator/ast"
    "calculator/lexer"
    "calculator/token"
    "calculator/nestingstack"
    "fmt"
)

const (
    INTEGER = "INTEGER"
    PLUS    = "PLUS"
    MINUS   = "MINUS"
    DIV     = "DIV"
    MUL     = "MUL"
    LPAR    = "LPAR"
    RPAR    = "RPAR"
    EOF     = "EOF"
)

type Parser struct {
    Lex *lexer.Lexer
    CurrentToken *token.Token
    Stack *nestingstack.NestingStack
}

func NewParser(lex *lexer.Lexer) (*Parser, error) {
    currentToken, err := lex.GetNextToken()
    stack := nestingstack.NewNestingStack() 
    return &Parser{Lex: lex, CurrentToken: currentToken, Stack: stack}, err
}


// Ensure current token type is consistent with the expected type, all parentheses are balanced,
// and there are no two successive integers with no operation between them. 
func (p *Parser) Consume(expectedType string) error {
    
    // check for type mismatch
    if p.CurrentToken.TokenType != expectedType {
        return fmt.Errorf("parser.Consume(): syntax error: expected %s but received %s",
            expectedType, p.CurrentToken.TokenType)
    }

    // if current token is LPAR then push an LPAR to the nesting stack 
    if p.CurrentToken.TokenType == LPAR {
        p.Stack.Push(token.Token{LPAR, '('})
    }
    // if current token is RPAR, then pop an LPAR from the nesting stack
    if p.CurrentToken.TokenType == RPAR {
        t, err := p.Stack.Peek() 
        if err != nil {
            return fmt.Errorf("parser.Parse(): unexpected ')'") // missing opening '('
        }
        if t.TokenType != LPAR {
            return fmt.Errorf("parser.Parse(): invalid character on stack") 
        }
        t, err = p.Stack.Pop()
        if err != nil {
            return err
        }
    }
    previousToken := p.CurrentToken // save current token before getting next token
    var err error
    p.CurrentToken, err = p.Lex.GetNextToken()
    if err != nil {
        return err
    }

    // check if an RPAR was read without a matching LPAR (nesting stack is empty)
    if p.CurrentToken.TokenType == RPAR {
        if p.Stack.IsEmpty() {
            return fmt.Errorf("parser.Parse(): unexpected ')'")
        }
    }
    // check for integers separated by white space 
    if previousToken.TokenType == INTEGER {
        if p.CurrentToken.TokenType == INTEGER {
            return fmt.Errorf("parser.Consume(): syntax error: missing op between integers")
        }
    }
    return nil // all tests passed
}


// Factor(): returns an ASTNode of type: UnaryOperation, INTEGER, or an Expr() subtree
func (p *Parser) Factor() (ast.ASTNode, error) {
   
    // factor: (PLUS| MINUS)|LPAR expr RPAR|INTEGER
    token := p.CurrentToken  
    switch token.TokenType {
    case PLUS:
        // unary operation: +
        if err := p.Consume(PLUS); err != nil {
            return ast.NewErrorNode(err), err // consume failed, returns error node and error 
        }
        unaryChild, err := p.Factor() // expression following the unary operation
        if err != nil {
            return ast.NewErrorNode(err), err
        }
        unaryNode := ast.NewUnaryOperation(token, unaryChild)
        return unaryNode, nil

    case MINUS:
        // unary operation: - 
        if err := p.Consume(MINUS); err != nil {
            return ast.NewErrorNode(err), err
        }
        unaryChild, err := p.Factor() 
        if err != nil {
            return ast.NewErrorNode(err),err
        }
        unaryNode := ast.NewUnaryOperation(token, unaryChild)
        return unaryNode, nil

    case INTEGER: 
        if err := p.Consume(INTEGER); err != nil {
            return ast.NewErrorNode(err), err
        }
        // token was correct type, create integer node
        integerNode, err := ast.NewNumberLiteral(token)
        if err != nil {
            return ast.NewErrorNode(err), err // if type assertion fails 
        }
        // token was correct type and type assertion passed, return integer node 
        return integerNode, nil
     
    case LPAR:
         // ( expr ) 
        if err := p.Consume(LPAR); err != nil {
            return ast.NewErrorNode(err), err
        }
        // get the expression within the ()'s
        subTreeRoot, err := p.Expr()
        if err != nil {
            return ast.NewErrorNode(err), err
        }
        // ensure closing parenthesis on return 
        if err := p.Consume(RPAR); err != nil {
            return ast.NewErrorNode(err), err
        }
        // all tests passed: returns the expression within the ()'s
        return subTreeRoot, nil

     default:
        err := fmt.Errorf("parser.Factor(): unexpected %s",p.CurrentToken.TokenType)
        return ast.NewErrorNode(err),err
    }
}


// Term(): returns an ASTNode: a subtree with MUL or DIV as the root, an INTEGER leaf node, or UnaryOp
func (p *Parser) Term() (ast.ASTNode, error) {

    // term: factor((MUL|DIV)factor)*
    leftChild, err := p.Factor()  
    if err != nil {
        return ast.NewErrorNode(err), err 
    }
    for p.CurrentToken.TokenType == MUL || p.CurrentToken.TokenType == DIV { 
        token := p.CurrentToken

        // get operation type 
        switch p.CurrentToken.TokenType {
        case MUL:
            if err := p.Consume(MUL); err != nil { 
                return ast.NewErrorNode(err), err
            }
        case DIV:
            if err := p.Consume(DIV); err != nil {
                return ast.NewErrorNode(err), err
            }
        default:
            return ast.NewErrorNode(err), fmt.Errorf("parser.Term() reached default case")
        }
        // get rightChild (integer leaf node or addition/subtraction subtree)
        rightChild, err := p.Factor()
        if err != nil {
            return ast.NewErrorNode(err), err
        }
        // create new node with leftChild from previous iteration and rightChild from current iteration 
        leftChild = ast.NewBinaryOperation(leftChild, rightChild, token)

    }
    //fmt.Println(leftChild)
    return leftChild, nil 
}


// Expr(): Returns an ASTNode: a subtree with PLUS or MINUS as the root 
func (p *Parser) Expr() (ast.ASTNode, error) {

    // expr: term((PLUS|MINUS)term)*    
    leftChild, err := p.Term()  
    if err != nil {
        return ast.NewErrorNode(err), err 
    }
    for p.CurrentToken.TokenType == PLUS || p.CurrentToken.TokenType == MINUS {       
        token := p.CurrentToken // preserve current token
        switch p.CurrentToken.TokenType {
        case PLUS:
            if err := p.Consume(PLUS); err != nil {
                return ast.NewErrorNode(err), err
            }
        case MINUS:
            if err := p.Consume(MINUS); err != nil {
                return ast.NewErrorNode(err), err
            }
        default:
            return ast.NewErrorNode(err), fmt.Errorf("parser.Expr(): unexpected %s",token.TokenType)
        }
        rightChild, err := p.Term()
        if err != nil {
            return ast.NewErrorNode(err), err
        }
        leftChild = ast.NewBinaryOperation(leftChild, rightChild, token) //recursively build LC return node
    }
    return leftChild, nil 
}

// final return point to Interpreter: returns root of AST to interpreter 
func (p *Parser) Parse() (ast.ASTNode, error) {
    rootNode, err := p.Expr()
    if err != nil {
        return ast.NewErrorNode(err), err
    }

   // ~~ for debugging ~~ 
   // fmt.Println(rootNode)
   // fmt.Printf("parser.Parse(): return stack size: %d\n",parser.Stack.StackSize())
   // fmt.Printf("parser.Parse(): return token:")    

    // make sure stack is empty 
    if !p.Stack.IsEmpty() {
        err := fmt.Errorf("parser.Parse(): missing opening or closing parentheses: parentheses not balanced")
        return ast.NewErrorNode(err),err
    }    
    // make sure all input was parsed 
    if p.CurrentToken.TokenType != EOF {
        err := fmt.Errorf("parser.Parse(): unexpected input at end of expression")
        return ast.NewErrorNode(err),err
    }
    // everything went well: return the AST of the input to the interpreter 
    return rootNode, err
}
