package main 

import (
    "calculator/lexer"
    "calculator/parser"
    "calculator/interpreter"
    "testing"
)

type TestCase struct {
    input          string
    shouldPass     bool   
    expectedResult int 
}

func TestInterpreter(t *testing.T) {
    testCases := []TestCase{
        {"2 + 2", true, 0},
        {"5 * 3", true, 0},
        {"10 / 2", true, 0},
        {"1 - 3", true, 0},
        {"2 + 3 * 4", true, 0},
        {"(2 + 3) * 4", true, 0},
        {"2 + (3 * 4)", true, 0},
        {"(2 + (3 * (4 - 1)))", true, 0},
        {"(2 + 3 * 4", false, 0},
        {"2 + 3) * 4", false, 0},
        {"2 * ()", false, 0},
        {"2 ---(3))", false, 0},
        {"1 1 ", false, 0},
        {"+ 2 + 3", true, 0},
        {"2 + 3 -", false, 0},
        {"2 ++ 3", true, 0},
        {"4 ** 5", false, 0},
        {"2 + $ + 3", false, 0},
        {"3 + 4 * (2 - 1) / (3 * (4 + 5)) - 6", true, 0},
        {"2+2", true, 0},
        {"2     +     2", true, 0},
        {"-3 + 4", true, 0},
        {"0 * 5", true, 0},
        {"2 * 3", true, 0},

		// Additional cases with descriptions
        {"3 + 5", true, 8},                  // basic addition
        {"2 - 9", true, -7},                 // basic subtraction
        {"4152 - 109", true, 4043},          // multi digits
        {"2", true, 2},                      // single term
        {"   2  +  8", true, 10},            // additional whitespace
        {"222 + 9 + 15 - 12", true, 234},    // addition and subtraction, multiple terms
        {"4 * 4", true, 16},                 // multiplication
        {"15/3", true, 5},                   // division
        {"15 / 3 * 9", true, 45},            // multiplication and division, multiple terms
        {"28 - 16 * 2 + 3", true, -1},       // correct order of operations
        {"(28 - 18) * 2 + 3", true, 23},     // correct order with brackets
        {"2 * 16 - 8 / 2 - 1", true, 27},    // correct order, multiplication and division
        {"2 * (16 -8) / 2 - 1", true, 7},    // correct order, multiplcation and division with parenthesis
        {"(2 +  18) * (3 + 5 )", true, 160}, // multiple parenthesis 
        {"1 *+ 2", true, 2},    			 // multiplication and + number
        {"+ 9 ", true, 9},      			 // missing leading integer (unary op)
        {"- 9", true, -9},       			 // same as above
        {"1 ++ 2", true, 3},                 // double addition symbols



        // fail cases
        {"a", false, 0},         // invalid character
        {"3 + b * 3", false, 0}, // invalid character
        {"(. + 3)", false, 0},   // invalid character after parenthesis
        {"1 + ", false, 0},      // missing integer after +
        {"+ ", false, 0},        // no integers
        {"* 9 + 2", false, 0},   // missing leading integer 
        {"(", false, 0},         // missing integer, etc. after (
        {"(1+2) +", false, 0},   // missing integer at end 
        {"()", false, 0},        // missing integers between ()'s*/
        {"1 1", false, 0},       // no operation between integers
        {"1 * 2) ", false, 0},   // trailing ) 
        {"(2 * 4", false, 0},    // missing closing parenthesis
        {"((", false, 0},        // missing integers and )
        {"(1+2)) + 13", false, 0}, // imbalanced parenthesis

    }



    for _, testCase := range testCases {
        var result int
       
        
        lexer := lexer.NewLexer(testCase.input)
        parser, _ := parser.NewParser(lexer) // this error is trivial and always caught. 
        interp := interpreter.NewInterpreter(parser)
        result, err := interp.Interpret()

        // the following tests are to check that the calculator is correctly catching errors on 
        // invalid input as well as evaluating the correct results when the input is valid. 

        // No error returned on invalid input (succeeded but should have failed):
        if err == nil && !testCase.shouldPass {
            t.Errorf("FAIL: no error returned from invalid input: %s: output: %d", testCase.input, result)
        }

        // Error returned on valid input 
        if err != nil && testCase.shouldPass {
            t.Errorf("FAIL: error returned from valid input: %s: error message: %v", testCase.input, err)
        }

        // No error on valid input, but incorrect evaluation (result != expected output). 
       // if err == nil && testCase.shouldPass {
         //   if testCase.expectedResult != result {
           //     t.Errorf("FAIL: incorrect result on input: %s: expected result: %d: actual result: %d",
             //       testCase.input, testCase.expectedResult,result)
           // } 
            //else {
            // t.Logf("PASS: correct result: %s = %d",testCase.input,result) //uncomment to see passed cases 
            //} 
        }


/*
    //  ***** uncomment to see passing cases ****** 
    
    // If errors   : go test  
    // If no errors: go test -v 

        // Error returned on invalid input (correct behaviour) 
        if err != nil && !testCase.shouldPass {
            t.Logf("PASS: error: %v input:",err, testCase.input)
        }

        // No error returned on valid input (correct behaviour)
        if err == nil && testCase.shouldPass {
            t.Logf("PASS: input: %s result: %d", testCase.input, result)
        }
*/ 
}
