package compiler

import (
	"fmt"
	"reflect"
	"testing"
	"toy_lang/bytecode"
	"toy_lang/lexer"
	"toy_lang/parser"
)

type tType struct {
	input  string
	output []bytecode.Instruction
	id     int
}

func compareInputs(t *testing.T, wantI tType, got []bytecode.Instruction) {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	want := wantI.output

	stderr := ""
	for i, val := range want {
		if !reflect.DeepEqual(val, got[i]) {
			stderr += fmt.Sprintf("%v[ERROR] Test number %v has failed,%v \nwanted %v, \ngot %v%v\n", Red, wantI.id, Reset, val, got[i], Reset)
		}
	}
	if stderr != "" {
		t.Error(stderr)
	} else {
		fmt.Printf("%v[PASS] Test number %d has passed%v\n", Green, wantI.id, Reset)
	}

}

func TestCompiler(t *testing.T) {
	tests := []tType{
		{
			input: "4",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
			},
			id: 1,
		},
		{
			input: "4 + 3;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
			},
			id: 2,
		},
		{
			input: "4 + 2 - 2;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 2},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
				&bytecode.LOAD_INT_INS{Address: 3, Value: 2},
				&bytecode.INFIX_INS{Left_addr: 2, Right_addr: 3, Save_to_addr: 4, Operation: 2},
			},
			id: 3,
		},
		{
			input: "4 + 3 * 2",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.LOAD_INT_INS{Address: 2, Value: 2},
				&bytecode.INFIX_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 3},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 3, Save_to_addr: 4, Operation: 1},
			},
			id: 4,
		},
		{
			input: "let x = 4;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 0},
			},
			id: 5,
		},
		{
			input: "let x = 4 + 3;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 2},
			},
			id: 6,
		},
		{
			input: "let x = 5; x = 4;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 5},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 0},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 4},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 1},
			},
			id: 7,
		},
		{
			input: "let x = 5; let y = x * 3;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 5},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 0},
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 1},
				&bytecode.LOAD_INT_INS{Address: 2, Value: 3},
				&bytecode.INFIX_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 3},
			},
			id: 8,
		},
		{
			input: "let x = true;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_BOOL_INS{Address: 0, Value: true},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 0},
			},
			id: 9,
		},
		{
			input: "let x = 5 < 3;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 5},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 5},
			},
			id: 10,
		},
		{
			input: "let x = true || false;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_BOOL_INS{Address: 0, Value: true},
				&bytecode.LOAD_BOOL_INS{Address: 1, Value: false},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 12},
			},
			id: 11,
		},
		{
			input: "if true{let x = 2;} else {let x = 4;}",
			output: []bytecode.Instruction{
				&bytecode.LOAD_BOOL_INS{Address: 0, Value: true},
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 0, TargetAddr: 5},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 2},
				&bytecode.DECLARE_VAR_INS{Addr: 1, Name: "x"},
				&bytecode.JMP_INS{InstNum: 7},
				&bytecode.LOAD_INT_INS{Value: 4, Address: 2},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 2},
			},
			id: 12,
		},
		{
			input: "let x = 4 * (3 + 2);",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.LOAD_INT_INS{Address: 2, Value: 2},
				&bytecode.INFIX_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 1},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 3, Save_to_addr: 4, Operation: 3},
			},
			id: 13,
		},
		{
			input: "fn add(a, b){return a + b;} let c = add(2, 3);",
			output: []bytecode.Instruction{
				&bytecode.FUNC_DEC_START_INS{Name: "add", ParamCount: 2, ParamNames: []string{"a", "b"}},
				&bytecode.REF_VAR_INS{Name: "a", SaveTo: 0},
				&bytecode.REF_VAR_INS{Name: "b", SaveTo: 1},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
				&bytecode.RETURN_INS{Ptr: 2},
				&bytecode.FUNC_DEC_END_INS{},
				&bytecode.LOAD_INT_INS{Address: 3, Value: 2},
				&bytecode.LOAD_INT_INS{Address: 4, Value: 3},
				&bytecode.FUNC_CALL_INS{Params: []int{3, 4}, Name: "add", PutRet: 5},
				&bytecode.DECLARE_VAR_INS{Name: "c", Addr: 5},
			},
			id: 14,
		},
		{
			input: "let x = 5 % 2;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 5},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 2},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 13},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 2},
			},
			id: 15,
		},
		{
			input: `let s = "hello world";`,
			output: []bytecode.Instruction{
				&bytecode.LOAD_STRING_INS{Address: 0, Value: "hello world"},
				&bytecode.DECLARE_VAR_INS{Name: "s", Addr: 0},
			},
			id: 16,
		},
		{
			input: `println("hello world")`,
			output: []bytecode.Instruction{
				&bytecode.LOAD_STRING_INS{Address: 0, Value: "hello world"},
				&bytecode.CALL_BUILTIN_INS{Name: "println", Params: []int{0}, PutRet: 1},
			},
			id: 17,
		},
		{
			input: `while true{println("hello")}`,
			output: []bytecode.Instruction{
				&bytecode.LOAD_BOOL_INS{Address: 0, Value: true},
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 0, TargetAddr: 5},
				&bytecode.LOAD_STRING_INS{Address: 1, Value: "hello"},
				&bytecode.CALL_BUILTIN_INS{Name: "println", Params: []int{1}, PutRet: 2},
				&bytecode.JMP_INS{InstNum: 0},
			},
			id: 18,
		},
		{
			input: "let x = 0; while x < 10{x++;}",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 0},
				&bytecode.DECLARE_VAR_INS{Addr: 0, Name: "x"},
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 1},
				&bytecode.LOAD_INT_INS{Address: 2, Value: 10},
				&bytecode.INFIX_INS{Left_addr: 1, Right_addr: 2, Operation: 5, Save_to_addr: 3},
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 3, TargetAddr: 11},
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 4},
				&bytecode.LOAD_INT_INS{Address: 5, Value: 1},
				&bytecode.INFIX_INS{Left_addr: 4, Right_addr: 5, Operation: 1, Save_to_addr: 6},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 6},
				&bytecode.JMP_INS{InstNum: 1},
			},
			id: 19,
		},
		{
			input: "let x = 0; while x < 100{x++; if x == 7{continue;} if x == 8{break;}}",
			output: []bytecode.Instruction{
				// let x = 0
				&bytecode.LOAD_INT_INS{Address: 0, Value: 0},
				&bytecode.DECLARE_VAR_INS{Addr: 0, Name: "x"},

				// while condition start
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 1},                                     //0
				&bytecode.LOAD_INT_INS{Address: 2, Value: 100},                                  //1
				&bytecode.INFIX_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 5}, //3 // x < 100
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 3, TargetAddr: 21},                         //4

				// x++
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 4},                                     //5
				&bytecode.LOAD_INT_INS{Address: 5, Value: 1},                                    //6
				&bytecode.INFIX_INS{Left_addr: 4, Right_addr: 5, Save_to_addr: 6, Operation: 1}, //7
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 6},                                   //8

				// if x == 7
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 7},                                     //9
				&bytecode.LOAD_INT_INS{Address: 8, Value: 7},                                    //10
				&bytecode.INFIX_INS{Left_addr: 7, Right_addr: 8, Save_to_addr: 9, Operation: 9}, //11  // x == 7
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 9, TargetAddr: 14},                         //12                         // skip continue if false
				&bytecode.JMP_INS{InstNum: 3},                                                   //13                                                   // continue → back to while condition

				// if x == 8
				&bytecode.REF_VAR_INS{Name: "x", SaveTo: 10},                                       //14
				&bytecode.LOAD_INT_INS{Address: 11, Value: 8},                                      //15
				&bytecode.INFIX_INS{Left_addr: 10, Right_addr: 11, Save_to_addr: 12, Operation: 9}, //16 // x == 8
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 12, TargetAddr: 19},                           //17                        // skip break if false
				&bytecode.JMP_INS{InstNum: 21},                                                     //18                                                  // break → exit loop

				// jump back to start of while condition
				&bytecode.JMP_INS{InstNum: 1}, //19
			},
			id: 20,
		},
		{
			input: "let pi = 3.1415965;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_FLOAT_INS{Address: 0, Value: 3.1415965},
				&bytecode.DECLARE_VAR_INS{Name: "pi", Addr: 0},
			},
			id: 21,
		},
		{
			input: "let x = 3.4 * 9.2;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_FLOAT_INS{Address: 0, Value: 3.4},
				&bytecode.LOAD_FLOAT_INS{Address: 1, Value: 9.2},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 3},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 2},
			},
			id: 22,
		},
		{
			input: "let x = 1.0 == 1.0;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_FLOAT_INS{Address: 0, Value: 1.0},
				&bytecode.LOAD_FLOAT_INS{Address: 1, Value: 1.0},
				&bytecode.INFIX_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 9},
				&bytecode.DECLARE_VAR_INS{Addr: 2, Name: "x"},
			},
			id: 23,
		},
	}
	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := parser.NewParser()
		compile := NewCompiler()

		got := compile.Compile(parse.Parse(lex.Lex(tt.input)))
		compareInputs(t, tt, got)
	}
}
