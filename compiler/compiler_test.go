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
			stderr += fmt.Sprintf("%v[ERROR] Test number %v has failed, \nwanted %v, \ngot %v%v\n", Red, wantI.id, val, got[i], Reset)
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
				&bytecode.JMP_IF_FALSE_INS{CondAddr: 0, TargetAddr: 5 },
				&bytecode.LOAD_INT_INS{Address: 1, Value: 2},
				&bytecode.DECLARE_VAR_INS{Addr: 1, Name: "x"},
				&bytecode.JMP_INS{InstNum: 7},
				&bytecode.LOAD_INT_INS{Value: 4, Address: 2},
				&bytecode.DECLARE_VAR_INS{Name: "x", Addr: 2},
			},
			id: 12,
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
