/*
The interpreter defines the and executes the Visit methods for each node type as it traverses the AST. If the input was
successfully parsed, then the input expression is calculated and returned. Otherwise, an error is returned. 
*/

package interpreter

import (
    "calculator/parser"
    "calculator/ast"
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


// the interpreter 
type Interpreter struct {
    Parser *parser.Parser
}

func NewInterpreter(parser *parser.Parser) *Interpreter { 
    return &Interpreter{Parser: parser}
}

func (interp *Interpreter) VisitBinaryOperation(node *ast.BinaryOperation) (interface{}, error) {
    leftResult, err := node.LeftChild.Accept(interp) // recursively evaluate left child 
    if err != nil {
        return 0, err
    }
    leftValue, ok := leftResult.(int) // type assertion on left result value
    if !ok {
        return 0, 
        fmt.Errorf("interpreter.VisitBinaryOperation(): leftResult evaluation returned non-integer value")
    }
    rightResult, err := node.RightChild.Accept(interp) // recursively evaluate right child 
    if err != nil {
        return nil, err
    }
    rightValue, ok := rightResult.(int) // type assertion on right result value 
    if !ok {
        if errorNode, isErrorNode := rightResult.(*ast.ErrorNode); isErrorNode {
            return 0, fmt.Errorf("interpreter encountered an error: %s", errorNode.ErrorType)
        }
       return nil, 
       fmt.Errorf("interpreter.VisitBinaryOperation(): rightResult returned a non-integer value: %v",rightResult)
    }

    // perform operation corresponding to BinaryNodeOperation operator type 
    switch node.Operator.TokenType {
    case PLUS:
        return leftValue + rightValue, nil
    case MINUS:
        return leftValue - rightValue, nil
    case MUL:
        return leftValue * rightValue, nil
    case DIV:
        if rightValue == 0 {
            return nil, fmt.Errorf("interpreter.VisitBinaryOperatrion(): division by zero") // div by 0 check
        }
        return leftValue / rightValue, nil
    default:
        return nil, fmt.Errorf("interpreter.VisitBinaryOperatrion(): default case reached")
    }
}

// Visit NumberLiteral: return integer value of the node
func (interp *Interpreter) VisitNumberLiteral(nl *ast.NumberLiteral) (interface{}, error) {
    return nl.Value, nil
}


// Visit UnaryOperation: recursively evaluates its child node and, if the operator type was negative 
// then it returns the result multiplied by -1. Otherwise it returns the result unmodified. 
func (interp *Interpreter) VisitUnaryOperation(node *ast.UnaryOperation) (interface{}, error) {
    
    exprResult, err := node.Expr.Accept(interp) // recursively evaluate child node 
    if err != nil {
        return 0, err
    }
    exprValue, ok := exprResult.(int) // type assertion 
    if !ok {
        return 0,
        fmt.Errorf("interpreter.UnaryOperation(): leftResult evaluation returned non-integer value")
    }
    switch node.Operator.TokenType {
    case PLUS:
        return exprValue, nil
    case MINUS:
        return -exprValue, nil // multiply result by -1 and return it 
    default:
        return 0, fmt.Errorf("VisitUnaryOperation(): default case reached")
    }
}

// Returns the error message in the node. Here for consistency 
func (interp *Interpreter) VisitErrorNode(en *ast.ErrorNode) (interface{}, error) {
    return en.ErrorType, en.ErrorType
}

// Interpret tree: initializes visitor pattern with the root node. 
func (interp *Interpreter) Interpret() (int, error) {
    root, err := interp.Parser.Parse()
    if err != nil {
        return 0, err // error returned from parser.Parse()
    }
    if root == nil {
        return 0, fmt.Errorf("interpreter.Interpret(): parser.Parse(): parsed an empty expression")
    }
    result, err := root.Accept(interp)
    if err != nil {
        return 0, err // error returned somewhere in interpretation
    }
    finalResult, ok := result.(int) // type assertion
    if !ok {
        if errorNode, isErrorNode := root.(*ast.ErrorNode); isErrorNode { // interpreter returned error node
            return 0, fmt.Errorf("interpreter encountered an error: %s", errorNode.ErrorType)
        }
        return 0, fmt.Errorf("interpreter.Interpret(): final result is a non-integer value")
    }
    return finalResult, nil
}
