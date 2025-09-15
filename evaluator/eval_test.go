package evaluator

import (
	"testing"
	"toy_lang/ast"
	"toy_lang/lexer"
	"toy_lang/parser"
)

func compareVMap(t *testing.T, got map[string]ast.Node, want map[string]ast.Node) {
	if len(got) != len(want) {
		t.Errorf("[FAIL] Wanted %d variables, Got %d variables", len(want), len(got))
	}
	for key, wantVal := range want {
		gotVal, ok := got[key]
		if !ok {
			t.Errorf("[FAIL] Missing variable %s in got map", key)
			continue
		}
		if gotVal.String() != wantVal.String() {
			t.Errorf("[FAIL] Wanted %v = %v, Got %v = %v", key, wantVal, key, gotVal)
		}
	}
}

func TestEvaluator(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]ast.Node
	}{
		{
			input: "let x = 0",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 0},
			},
		},
		{
			input: "let x = 9; x=x+3; let y = x/4",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 12},
				"y": &ast.IntLiteralNode{Value: 3},
			},
		},
		{
			input: "let x = true; let y = false;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = true; x = false",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = true; let y = false; let z = x || y;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
				"z": &ast.BoolLiteralNode{Value: true},
			},
		},
		{
			input: "let x = true; let y = x && false;",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = true; let y = !!x && true",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: true},
				"y": &ast.BoolLiteralNode{Value: true},
			},
		},
		{
			input: "let x = 9;if true{let y = 4; x = y;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},
			},
		},
		{
			input: "let x = 5; if 5 < 6{let y = true;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 5},
			},
		},
		{
			input: "let y = false; if y && true{let x = 5;}",
			output: map[string]ast.Node{
				"y": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = false; if !x&&true{let y = !x;}",
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = 0; if true {let y = 4; x = y;} else {let y = 5; x = y;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},
			},
		},
		{
			input: "let y = 9; if y < 10{y = 8;} else {y = 11;}",
			output: map[string]ast.Node{
				"y": &ast.IntLiteralNode{Value: 8},
			},
		},
		{
			input: "let v = true || false; if v{v = false;} else {v = true;}",
			output: map[string]ast.Node{
				"v": &ast.BoolLiteralNode{Value: false},
			},
		},
		{
			input: "let x = 4 * (4 + 2);",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 24},
			},
		},
		{
			input: "let x = 4 * (4 + 2); let y = 9 / (x + 1);",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 24},
				"y": &ast.IntLiteralNode{Value: 0},
			},
		},
		{
			input: "let x = 0;if true{if !false{let y = 4; x=y;}}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 4},	
			},
		},
	}

	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := parser.NewParser()
		exec := NewInterpreter()
		program := parse.Parse(lex.Lex(tt.input))
		compareVMap(t, exec.Execute(program, false).Vars, tt.output)
	}
}
