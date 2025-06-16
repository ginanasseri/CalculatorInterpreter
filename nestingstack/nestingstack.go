package nestingstack

/*

Used by the parser to ensure parentheses are balanced. 

*/

import (
    "calculator/token"
    "fmt"
)

type NestingStack struct {
    Stack []token.Token
}

// NestingStack constuctor: returns pointer a new NestingStack instance
func NewNestingStack() *NestingStack {
    return &NestingStack{Stack: make([]token.Token,0)}
}

func (ns *NestingStack) StackSize() int {
    return len(ns.Stack)
}

func (ns *NestingStack) IsEmpty() bool {
    return len(ns.Stack) == 0
}

// Remove element from top of stack and return it 
func (ns *NestingStack) Pop() (token.Token, error) {
    if ns.IsEmpty() {
        return token.Token{},
            fmt.Errorf("nestingstack.Pop(): cannot pop from empty stack") // returns nil token and error
    }
    index := len(ns.Stack) - 1
    element := ns.Stack[index]
    ns.Stack = ns.Stack[:index] // removes last element in stack (top of stack)
    return element, nil
}

// Return element at top of stack, but don't remove it. 
func (ns *NestingStack) Peek() (token.Token, error) {
    if ns.IsEmpty() {
        return token.Token{}, 
            fmt.Errorf("nestingstack.Peek(): stack is empty") // returns nil token and error message
    }
    return ns.Stack[len(ns.Stack) - 1], nil
}

// Add element to top of stack. 
func(ns *NestingStack) Push(value token.Token) {
    ns.Stack = append(ns.Stack, value)
}
