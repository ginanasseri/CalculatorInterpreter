/* 

Calculator which takes an arithmetic expression as input, evaluates it, and provides the result as output.
Accepts '+', '-', '*', '/', '(', ')', and integer strings. If input includes a character not in the 
alphabet of the interpreter, an error is returned indicating unrecognized character was in the input. 
If input includes a structural syntax error, a syntax error is returned. Otherwise, the result is returned.
*/

package main

import (
    "calculator/interpreter"
    "calculator/lexer"
    "calculator/parser"
    "fmt"
    "os"
    "bufio"

)




// errorTest checks if there was a scanning error 
func errorTest(err error) {
    if err != nil {
        fmt.Printf("Error reading input: %v\n", err)
        os.Exit(1)
    }
}



func main() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("-------------------------------------\n... Starting calculator... (Q = exit)")
    for {
        fmt.Print(">> ")
        if !scanner.Scan() {
            break // EOF or error
        }
        input := scanner.Text()
        if input == "" {
            continue
        }
        if input == "q" || input == "Q" {
            break
        }
 
        lexer := lexer.NewLexer(input)

        parser, err := parser.NewParser(lexer)
        if err != nil {
            fmt.Printf("%v\n",err)
            continue
        }
        interp := interpreter.NewInterpreter(parser)

        result, err1 := interp.Interpret()
        if err1 != nil {
            fmt.Printf("%v\n",err1)
            continue
        }
        fmt.Printf("result: %v\n",result)
    }
    errorTest(scanner.Err())
}
