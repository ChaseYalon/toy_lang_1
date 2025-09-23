package compiler

import (
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
	want := wantI.output

	for i, val := range want {
		if !reflect.DeepEqual(val, got[i]) {
			t.Errorf("[ERROR] Test number %v has failed, \nwanted %v, \ngot %v\n", wantI.id, val, got[i])
		}
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
				&bytecode.INFIX_INT_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
			},
			id: 2,
		},
		{
			input: "4 + 2 - 2;",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 2},
				&bytecode.INFIX_INT_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
				&bytecode.LOAD_INT_INS{Address: 3, Value: 2},
				&bytecode.INFIX_INT_INS{Left_addr: 2, Right_addr: 3, Save_to_addr: 4, Operation: 2},
			},
			id: 3,
		},
		{
			input: "4 + 3 * 2",
			output: []bytecode.Instruction{
				&bytecode.LOAD_INT_INS{Address: 0, Value: 4},
				&bytecode.LOAD_INT_INS{Address: 1, Value: 3},
				&bytecode.LOAD_INT_INS{Address: 2, Value: 2},
				&bytecode.INFIX_INT_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 3},
				&bytecode.INFIX_INT_INS{Left_addr: 0, Right_addr: 3, Save_to_addr: 4, Operation: 1},
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
				&bytecode.INFIX_INT_INS{Left_addr: 0, Right_addr: 1, Save_to_addr: 2, Operation: 1},
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
				&bytecode.INFIX_INT_INS{Left_addr: 1, Right_addr: 2, Save_to_addr: 3, Operation: 3},
			},
			id: 8,
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
