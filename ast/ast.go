/*
This package handles the creation and traversal of an AST. 

The ASTNode and ASTVisitor interfaces are designd to be implemented by each distinct node type in the AST
The interfaces work together using the Visitor pattern. 

The ASTNode interface is implemented by the node types in the AST (as well as the ErrorNode type) such that
each node type implements the Accept() method. (In addition to any other methods defined within the interface).
The Accept() method takes an ASTVisitor as an argument and enables the visitor to call the Visit__() method
corresponding to the type of the current node. The ASTVisitor interface declares the set of methods that any 
type of visitor implementing the interface must define. 

In this case, the visitor is the interpreter. The interpreter initiates the visit process by calling Accept() on 
the root node. This is the entry point for the visitor to begin its operations on the AST-- in this case, traversing
through the nodes of the AST which has been created by the parser in such a way that the correct order of operations
are followed. 

The visit process is initiated by the visitor, where the visit methods are defined. In this case, the visitor
is the interpreter. So when the interpreter, the visitor, calls Accept() on a node, the node then calls the
Visit method corresponding to its node type, passing an instance of itself to the Visit method. E.g., if the 
visitor calls Accept on a BinaryOperation node, then the node will call VisitBinaryOperation(self) which in turn
allows the visitor to perform the operation(s) corresponding to that node type (E.g., left result + right result).
The visit methods corresponding to node types with children may also invoke Accept() on a node, thereby allowing
the interpreter to traverse the AST via recursive calls to Accept() by each node's child nodes until the terminal
node for each branch (a NumberLiteral) is reached.
*/
package ast

import (
    "calculator/token"
    "fmt"
)

type ASTNode interface {
    Accept(v ASTVisitor) (interface{}, error)
    String() string
}

type ASTVisitor interface {
    VisitBinaryOperation(node *BinaryOperation) (interface{}, error)
    VisitUnaryOperation(node *UnaryOperation) (interface{}, error)
    VisitNumberLiteral(node *NumberLiteral) (interface{}, error)
    VisitErrorNode(node *ErrorNode) (interface {}, error)
}

// BinaryOperation nodes: An ASTNode where the left and right child are the operands 
type BinaryOperation struct {
    LeftChild ASTNode
    Operator *token.Token
    RightChild ASTNode
}

func NewBinaryOperation(leftChild, rightChild ASTNode, operator *token.Token) ASTNode {
    binaryOperation := &BinaryOperation{
        LeftChild: leftChild,
        Operator: operator,
        RightChild: rightChild,
    }
    return binaryOperation
}

func (bo *BinaryOperation) Accept(v ASTVisitor) (interface{}, error) {
    return v.VisitBinaryOperation(bo)
}

func (bo *BinaryOperation) String() string {
    return fmt.Sprintf("(%v %s %v)", bo.LeftChild, bo.Operator.TokenType, bo.RightChild)
}

// UnaryOperation nodes: has one child node which is the expression immediately following 
type UnaryOperation struct {
    Operator *token.Token
    Expr ASTNode 
}

func NewUnaryOperation(operator *token.Token, expr ASTNode) ASTNode {
    return &UnaryOperation{Operator: operator, Expr: expr}
}

func (uo *UnaryOperation) Accept(v ASTVisitor) (interface{}, error) {
    return v.VisitUnaryOperation(uo)
}

func (uo *UnaryOperation) String() string {
    return fmt.Sprintf("(%s)(%v)", uo.Operator.TokenType, uo.Expr)//, nl.Value)
}

// NumberLiteral nodes: the leaf nodes of the AST, hold integer values 
type NumberLiteral struct { 
    Token *token.Token
    Value int
}

func NewNumberLiteral(token *token.Token) (ASTNode, error) {
    value, ok := token.Value.(int) // type assertion 
    if !ok {
        return nil, fmt.Errorf("ast.NewNumberLiteral(): token.TokenValue is not an int")
    }
    numberLiteral := &NumberLiteral{Token: token, Value: value}
    return numberLiteral,nil 
}

func (nl *NumberLiteral) Accept(v ASTVisitor) (interface{}, error) {
    return v.VisitNumberLiteral(nl)
}

func (nl *NumberLiteral) String() string {
    return fmt.Sprintf("%d", nl.Value)
}


// ErrorNodes are in the case that an error arises, the calling method may still return a node and
// the error message can be saved within the node in case later on, the specific error and location
// in the traversal can be recovered if required.  
type ErrorNode struct {
    ErrorType error
}

func NewErrorNode(errorType error) ASTNode {
    return &ErrorNode{ErrorType: errorType}
}

func (e *ErrorNode) Accept(v ASTVisitor) (interface{}, error) {
    return v.VisitErrorNode(e)
}

func (en *ErrorNode) String() string {
    return fmt.Sprintf("Node: %v",en.ErrorType)
}

