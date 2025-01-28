package lexer

/* 
The lexer reads the raw input and creates tokens for the parser to evaluate. If a character not in the alphabet is
encountered, then it returns an error. Otherwise it returns a token with a token type and token value. 
*/

import (
     "calculator/token"
     "fmt"
     "unicode"
     "strconv"
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

type Lexer struct {
    Input string 
    Position int
    CurrentChar byte
}

func NewLexer(input string) *Lexer {
    var currentChar byte = 0
    if len(input) > 0 {
        currentChar = input[0]
    } 
    return &Lexer{Input: input, Position: 0, CurrentChar: currentChar}
    
}

func (lex *Lexer) GetNextChar() {
    lex.Position += 1
    if lex.Position > len(lex.Input) - 1 {
        lex.CurrentChar = 0
     } else {
        lex.CurrentChar = lex.Input[lex.Position]
    }
    return 
}

func (lex *Lexer) SkipWhiteSpace() {
    for lex.CurrentChar != 0 && unicode.IsSpace(rune(lex.CurrentChar)) {
        lex.GetNextChar()
    }
    return
}

// parse multi-digit integer 
func (lex *Lexer) Integer() (int, error) {
    integerString := ""
    for lex.CurrentChar != 0 && unicode.IsDigit(rune(lex.CurrentChar)) {
        integerString += string(lex.CurrentChar)
        lex.GetNextChar()
    }
    integerA, err := strconv.ParseInt(integerString, 10, 32)
    integer := int(integerA)
//    integer, err := strconv.Atoi(integerString)
    if err != nil {
        return 0, fmt.Errorf("failed to convert input to integer %w",err)
    }
    return integer, nil    
}

// Creates a token based on current character in the input, if the character read is not in the alphabet,
// then an error is returned, otherwise it returns the token.
func (lex *Lexer) GetNextToken() (*token.Token, error) {    
    for lex.CurrentChar != 0 {
        switch {
        case unicode.IsSpace(rune(lex.CurrentChar)):
            lex.SkipWhiteSpace()

        case unicode.IsDigit(rune(lex.CurrentChar)):
            integer, err := lex.Integer()
            if err != nil {
                return token.NewToken("",0), err
            }
            return token.NewToken(INTEGER, integer), nil

        case lex.CurrentChar == '+':
            lex.GetNextChar()
            return token.NewToken(PLUS, '+'), nil
        
        case lex.CurrentChar == '-':
            lex.GetNextChar()
            //fmt.Printf("lexer.GetNextToken(): returning MINUS token\n")
            return token.NewToken(MINUS, '-'), nil

        case lex.CurrentChar == '*':
            lex.GetNextChar()
            //fmt.Printf("lexer.GetNextToken(): (MUL, '*')\n")
            return token.NewToken(MUL, '-'), nil
            
        case lex.CurrentChar == '/':
            lex.GetNextChar()
            return token.NewToken(DIV, '/'), nil

        case lex.CurrentChar == '(':
            lex.GetNextChar()
            //fmt.Printf("lexer.GetNextToken(): (LPAR, '(')\n")
            return token.NewToken(LPAR, '('), nil
            
        case lex.CurrentChar == ')':
            lex.GetNextChar()
            //fmt.Printf("lexer.GetNextToken(): (RPAR, ')')\n")
            return token.NewToken(RPAR, ')'), nil

        default:
            return token.NewToken("",0), fmt.Errorf("lexer.GetNextToken(): invalid character: %c", lex.CurrentChar)
         }
     }
    //fmt.Printf("lexer.GetNextToken(): returning EOF token\n")
    return token.NewToken(EOF, 0), nil
}
