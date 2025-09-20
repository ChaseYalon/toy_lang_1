package evaluator

import (
	"testing"
	"fmt"
	"toy_lang/ast"
	"toy_lang/lexer"
	"toy_lang/parser"
)

func compareVMap(t *testing.T, got map[string]ast.Node, want map[string]ast.Node, tt tEvalRes) {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Blue = "\033[34m"
	var Yellow = "\033[33m"
	var stderr = "";
	if len(got) != len(want) {
		stderr += fmt.Sprintf("[FAIL] Wanted %d variables, Got %d variables", len(want), len(got))
	}
	for key, wantVal := range want {
		gotVal, ok := got[key]
		if !ok {
			stderr += fmt.Sprintf("[FAIL] Missing variable %s in got map", key)
			continue
		}
		if gotVal.String() != wantVal.String() {
			stderr += fmt.Sprintf("[FAIL] Wanted %v = %v, Got %v = %v", key, wantVal, key, gotVal)
		}
	}
	if stderr != ""{
		errorString := Red + fmt.Sprintf("[FAILURE] Test number %d has failed", tt.id) + Reset + "\n____________\n" + tt.input + "\n____________\n" +"ERROR\n" + stderr + "\n\n\n" + Blue + fmt.Sprintf("Full output\n %+v\n\n\n%v Correct output\n%+v", got, Yellow, want) + Reset; 
		t.Error(errorString);
	} else {
		fmt.Printf("%v[PASS] Test number %d has passed%v\n", Green, tt.id, Reset);
	}
}
type tEvalRes struct{
	input string
	output map[string]ast.Node
	id int
}
func TestEvaluator(t *testing.T) {
	tests := []tEvalRes{
		{
			input: "let x = 0",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 0},
			},
			id: 1,
		},
		{
			input: "let x = 9; x=x+3; let y = x/4",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 12},
				"y": &ast.IntLiteralNode{Value: 3},
			},
			id: 2,
		},
		{
			input: "let x = true; let y = false;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
			},
			id: 3,
		},
		{
			input: "let x = true; x = false",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: false},
			},
			id: 4,
		},
		{
			input: "let x = true; let y = false; let z = x || y;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
				"z": &ast.BoolLiteralNode{Value: true},
			},
			id: 5,
		},
		{
			input: "let x = true; let y = x && false;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
			},
			id: 6,
		},
		{
			input: "let x = true; let y = !!x && true",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: true},
			},
			id: 7,
		},
		{
			input: "let x = 9;if true{let y = 4; x = y;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},
			},
			id: 8,
		},
		{
			input: "let x = 5; if 5 < 6{let y = true;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 5},
			},
			id: 9,
		},
		{
			input: "let y = false; if y && true{let x = 5;}",
			output: map[string]ast.Node{
				"y": &ast.BoolLiteralNode{Value: false},
			},
			id: 10,
		},
		{
			input: "let x = false; if !x&&true{let y = !x;}",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: false},
			},
			id: 11,
		},
		{
			input: "let x = 0; if true {let y = 4; x = y;} else {let y = 5; x = y;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},
			},
			id: 12,
		},
		{
			input: "let y = 9; if y < 10{y = 8;} else {y = 11;}",
			output: map[string]ast.Node{
				"y": &ast.IntLiteralNode{Value: 8},
			},
			id: 13,
		},
		{
			input: "let v = true || false; if v{v = false;} else {v = true;}",
			output: map[string]ast.Node{
				"v": &ast.BoolLiteralNode{Value: false},
			},
			id: 14,
		},
		{
			input: "let x = 4 * (4 + 2);",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 24},
			},
			id: 15,
		},
		{
			input: "let x = 4 * (4 + 2); let y = 9 / (x + 1);",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 24},
				"y": &ast.IntLiteralNode{Value: 0},
			},
			id: 16,
		},
		{
			input: "let x = 0;if true{if !false{let y = 4; x=y;}}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},
			},
			id: 17,
		},
		{
			input: "fn add(a, b){return a + b;}; let c = add(2, 3);",
			output: map[string]ast.Node{
				"c": &ast.IntLiteralNode{Value: 5},
			},
			id: 18,
		},
		{
			input: "fn a(b){return b - 2;} fn c(b){return b + 2;} let d = a(2) + c(2);",
			output: map[string]ast.Node{
				"d": &ast.IntLiteralNode{Value: 4},
			},
			id: 19,
		},
	}

	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := parser.NewParser()
		exec := NewInterpreter()
		program := parse.Parse(lex.Lex(tt.input))
		compareVMap(t, exec.Execute(program, false).Vars, tt.output, tt)
	}
}
