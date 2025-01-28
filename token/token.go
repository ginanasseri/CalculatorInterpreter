package token

import (
    "fmt"
)

// token struct, has type and value
type Token struct {
    TokenType string
    Value     interface{}
}

// creates a new token, returns ptr to the new token
func NewToken(tokenType string, value interface{}) *Token {
    return &Token{TokenType: tokenType, Value: value}
}

// print statement for token struct 
func (t Token) String() string {
    return fmt.Sprintf("Token{%s, %v}", 
        t.TokenType, t.Value)
}

