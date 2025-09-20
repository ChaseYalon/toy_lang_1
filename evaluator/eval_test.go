package evaluator

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"toy_lang/ast"
	"toy_lang/lexer"
	"toy_lang/parser"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func compareVMap(t *testing.T, got map[string]ast.Node, want map[string]ast.Node, tt tEvalRes) {
	Reset := "\033[0m"
	Red := "\033[31m"
	Green := "\033[32m"
	Blue := "\033[34m"
	Yellow := "\033[33m"
	stderr := ""
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
	if stderr != "" {
		errorString := Red + fmt.Sprintf("[FAILURE] Test number %d has failed", tt.id) + Reset +
			"\n____________\n" + tt.input + "\n____________\n" +
			"ERROR\n" + stderr + "\n\n\n" +
			Blue + fmt.Sprintf("Full output\n %+v\n\n\n%v Correct output\n%+v", got, Yellow, want) + Reset
		t.Error(errorString)
	} else {
		fmt.Printf("%v[PASS] Test number %d has passed%v\n", Green, tt.id, Reset)
	}
}

type tEvalRes struct {
	input     string
	output    map[string]ast.Node
	want_str  string
	enter_str string
	id        int
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

		{
			input: `let x = "hi";`,
			output: map[string]ast.Node{
				"x": &ast.StringLiteralNode{Value: "hi"},
			},
			id: 20,
		},

		{
			input: `let x = "hello" + "world";`,
			output: map[string]ast.Node{
				"x": &ast.StringLiteralNode{Value: "helloworld"},
			},
			id: 21,
		},

		{
			input: `fn outStr(a, b){return a + b;} let h = outStr("hello", "world");`,
			output: map[string]ast.Node{
				"h": &ast.StringLiteralNode{Value: "helloworld"},
			},
			id: 22,
		},

		{
			input: `fn concat(a, b){return a + b;} let hello = concat("hello ", "world");`,
			output: map[string]ast.Node{
				"hello": &ast.StringLiteralNode{Value: "hello world"},
			},
			id: 23,
		},

		{
			input: `fn hello(){return "hello";}let x = ""; if hello() == "hello"{x = "equals";} else {x = "not equals";}`,
			output: map[string]ast.Node{
				"x": &ast.StringLiteralNode{Value: "equals"},
			},
			id: 24,
		},

		{
			input:    `print("hello world");`,
			want_str: "hello world",
			id:       25,
		},

		{
			input:    `println("hello " + 2);`,
			want_str: "hello 2\n",
			id:       26,
		},

		{
			input:    `println(true);`,
			want_str: "true\n",
			id:       27,
		},

		{
			input:     `let x = input("Enter your name: ");`,
			enter_str: "Chase",
			output: map[string]ast.Node{
				"x": &ast.StringLiteralNode{Value: "Chase"},
			},
			id: 28,
		},

		{
			input: `let st = str(1); print(st + 2);`,
			output: map[string]ast.Node{
				"st": &ast.StringLiteralNode{Value: "1"},
			},
			want_str: "12",
			id:       29,
		},

		{
			input: `let i = int("42"); let d = i * 6;`,
			output: map[string]ast.Node{
				"i": &ast.IntLiteralNode{Value: 42},
				"d": &ast.IntLiteralNode{Value: 42 * 6}, //To lazy to open a calculator
			},
			id: 30,
		},

		{
			input: `let x = !true || bool("false");`,
			output: map[string]ast.Node{
				"x": &ast.BoolLiteralNode{Value: false},
			},
			id: 31,
		},

		{
			input: "let x = 0; while x < 10{x++;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 10},
			},
			id: 32,
		},

		{
			input: "let x = 0; while x < 10{break;}",
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 0},
			},
			id: 33,
		},

		{
			input: `println("bye");let x = 0; while x < 10{x++; continue; print("hi");}`,
			output: map[string]ast.Node{
				"x": &ast.IntLiteralNode{Value: 10},
			},
			want_str: "bye\n",
			id:       34,
		},

		{
			input: "let x = 3.1415",
			output: map[string]ast.Node{
				"x": &ast.FloatLiteralNode{Value: 3.1415},
			},
		},
	}

	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := parser.NewParser()
		exec := NewInterpreter()
		program := parse.Parse(lex.Lex(tt.input))

		if tt.enter_str != "" {
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			w.WriteString(tt.enter_str + "\n")
			w.Close()
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()
		}

		if tt.output != nil {
			compareVMap(t, exec.Execute(program, false).Vars, tt.output, tt)
		}
		if tt.want_str != "" {
			out := captureOutput(func() {
				exec.Execute(program, false)
			})
			if out != tt.want_str {
				t.Errorf("[FAILURE] Test number %d has failed\nGot: %q\nWant: %q\n", tt.id, out, tt.want_str)
			} else {
				fmt.Printf("\033[32m[PASS] Test number %d has passed\033[0m\n", tt.id)
			}
		}
	}
}
