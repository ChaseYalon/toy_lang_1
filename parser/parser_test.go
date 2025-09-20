package parser

import (
	"fmt"
	"testing"
	"toy_lang/ast"
	"toy_lang/lexer"
	"toy_lang/token"
)

// deepCompare compares two AST nodes recursively
func deepCompare(got, want ast.Node) bool {
	if got == nil || want == nil {
		return got == want
	}

	// Handle EmptyExprNode specially
	if g, ok := got.(*ast.EmptyExprNode); ok {

		return deepCompare(g.Child, want)
	}
	if w, ok := want.(*ast.EmptyExprNode); ok {
		return deepCompare(got, w.Child)
	}
	if got.NodeType() == ast.LetStmt && want.NodeType() == ast.LetStmt {
		g, ok := got.(*ast.LetStmtNode)
		w, ok2 := got.(*ast.LetStmtNode)
		if !ok && !ok2 {
			panic("[TESTS FAILED TO RUN]")
		}
		namesEq := g.Name == w.Name
		valueEq := deepCompare(g.Value, w.Value)
		return namesEq && valueEq
	}

	// If types differ, they’re not equal
	if got.NodeType() != want.NodeType() {
		return false
	}

	// Compare string forms if node types match
	return got.String() == want.String()
}

// compareNodes compares slices of AST nodes and prints mismatches
func compareNodes(t *testing.T, got, want []ast.Node, tt ttype) {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Blue = "\033[34m"
	var Yellow = "\033[33m"

	var stderr string

	if len(got) != len(want) {
		stderr += fmt.Sprintf("Length mismatch: got %d, want %d\n", len(got), len(want))
	}

	minLen := len(got)
	if len(want) < minLen {
		minLen = len(want)
	}

	for i := 0; i < minLen; i++ {
		if !deepCompare(got[i], want[i]) {
			stderr += fmt.Sprintf("Mismatch at index %d:\n Got:  %v\n Want: %v\n",
				i, got[i], want[i])
		}
	}

	// Check extras
	if len(got) > len(want) {
		for i := len(want); i < len(got); i++ {
			stderr += fmt.Sprintf("Extra element in got at index %d: %v\n", i, got[i])
		}
	} else if len(want) > len(got) {
		for i := len(got); i < len(want); i++ {
			stderr += fmt.Sprintf("Missing element in got at index %d: %v\n", i, want[i])
		}
	}

	if stderr != "" {
		errString := Red + fmt.Sprintf("[FAILURE] Test number %d has failed", tt.id) + Reset +
			fmt.Sprintf("\n____________\nInput: %v\n____________\nERROR\n %v\n", tt.input, stderr) +
			Blue + fmt.Sprintf("Full output\n%+v\n\n\n", got) +
			Yellow + fmt.Sprintf("Correct output \n%+v\n\n\n", want) + Reset
		t.Error(errString)
	} else {
		passString := Green + fmt.Sprintf("[PASS] Test number %d has passed", tt.id) + Reset
		fmt.Println(passString)
	}
}

func compareTokens(t *testing.T, got, want []token.Token) {
	if len(got) != len(want) {
		t.Errorf("Length mismatch: got %d, want %d", len(got), len(want))
	}

	minLen := len(got)
	if len(want) < minLen {
		minLen = len(want)
	}

	for i := 0; i < minLen; i++ {
		if got[i] != want[i] {
			t.Errorf("Mismatch at index %d: got %+v, want %+v", i, got[i], want[i])
		}
	}

	if len(got) > len(want) {
		for i := len(want); i < len(got); i++ {
			t.Errorf("Extra element in got at index %d: %+v", i, got[i])
		}
	} else if len(want) > len(got) {
		for i := len(got); i < len(want); i++ {
			t.Errorf("Missing element in got at index %d: %+v", i, want[i])
		}
	}
}

func TestPreParser(t *testing.T) {
	lex := lexer.NewLexer()
	parse := NewParser()

	tests := []struct {
		input  string
		output []token.Token
	}{
		{
			input: "let x = 5; x+=1;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.PLUS, "+"),
				*token.NewToken(token.INTEGER, "1"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},

		{
			input: "let a = 2; let b = 3; a+=b;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "a"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "2"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "b"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "3"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "a"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "a"),
				*token.NewToken(token.PLUS, "+"),
				*token.NewToken(token.VAR_REF, "b"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},

		{
			input: "let x = 10; x-=x;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "10"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.MINUS, "-"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let x = 1; x*=2; x/=x;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "1"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.MULTIPLY, "*"),
				*token.NewToken(token.INTEGER, "2"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.DIVIDE, "/"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let y = -4;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "0"),
				*token.NewToken(token.MINUS, "-"),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let y = 5 + (-4 + 1);",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.PLUS, "+"),
				*token.NewToken(token.LPAREN, "("),
				*token.NewToken(token.INTEGER, "0"),
				*token.NewToken(token.MINUS, "-"),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.PLUS, "+"),
				*token.NewToken(token.INTEGER, "1"),
				*token.NewToken(token.RPAREN, ")"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
	}

	for _, tt := range tests {
		res := parse.preProcess(lex.Lex(tt.input))
		compareTokens(t, res, tt.output)
	}
}

type ttype struct {
	input  string
	output ast.ProgramNode
	id     int
}

func TestParser(t *testing.T) {

	tests := []ttype{
		{
			input: "let x = 4;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.IntLiteralNode{
							Value: 4,
						},
					},
				},
			},
			id: 1,
		},
		{
			input: "let x = 4 + 5;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.InfixExprNode{
							Left: &ast.IntLiteralNode{
								Value: 4,
							},
							Operator: token.PLUS,
							Right: &ast.IntLiteralNode{
								Value: 5,
							},
						},
					},
				},
			},
			id: 2,
		},
		{
			input: "let x = 9; x=x+3;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.IntLiteralNode{
							Value: 9,
						},
					},
					&ast.VarReassignNode{
						Var: ast.ReferenceExprNode{
							Name: "x",
						},
						NewVal: &ast.InfixExprNode{
							Left: &ast.ReferenceExprNode{
								Name: "x",
							},
							Operator: token.PLUS,
							Right: &ast.IntLiteralNode{
								Value: 3,
							},
						},
					},
				},
			},
			id: 3,
		},
		{
			input: "let x = 9; x = x + 3; let y = x / 4;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.IntLiteralNode{
							Value: 9,
						},
					},
					&ast.VarReassignNode{
						Var: ast.ReferenceExprNode{
							Name: "x",
						},
						NewVal: &ast.InfixExprNode{
							Right: &ast.IntLiteralNode{
								Value: 3,
							},
							Operator: token.PLUS,
							Left: &ast.ReferenceExprNode{
								Name: "x",
							},
						},
					},
					&ast.LetStmtNode{
						Name: "y",
						Value: &ast.InfixExprNode{
							Left: &ast.ReferenceExprNode{
								Name: "x",
							},
							Operator: token.DIVIDE,
							Right: &ast.IntLiteralNode{
								Value: 4,
							},
						},
					},
				},
			},
			id: 4,
		},
		{
			input: "let x = true; let y = 4;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.BoolLiteralNode{
							Value: true,
						},
					},
					&ast.LetStmtNode{
						Name: "y",
						Value: &ast.IntLiteralNode{
							Value: 4,
						},
					},
				},
			},
			id: 5,
		},
		{
			input: "let x = true || false;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.BoolInfixNode{
							Left: &ast.BoolLiteralNode{
								Value: true,
							},
							Operator: token.OR,
							Right: &ast.BoolLiteralNode{
								Value: false,
							},
						},
					},
				},
			},
			id: 6,
		},
		{
			input: "let x = true; let y = !x && !true;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.BoolLiteralNode{
							Value: true,
						},
					},
					&ast.LetStmtNode{
						Name: "y",
						Value: &ast.BoolInfixNode{
							Left: &ast.PrefixExprNode{
								Operator: token.NOT,
								Value:    &ast.ReferenceExprNode{Name: "x"},
							},
							Operator: token.AND,
							Right: &ast.PrefixExprNode{
								Operator: token.NOT,
								Value:    &ast.BoolLiteralNode{Value: true},
							},
						},
					},
				},
			},
			id: 7,
		},
		{
			input: "let x = 5 < 6;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.BoolInfixNode{
							Left:     &ast.IntLiteralNode{Value: 5},
							Operator: token.LESS_THAN,
							Right:    &ast.IntLiteralNode{Value: 6},
						},
					},
				},
			},
			id: 8,
		},
		{
			input: "let x = 5 >= 6; let y = !x || true;",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.BoolInfixNode{
							Left:     &ast.IntLiteralNode{Value: 5},
							Operator: token.GREATER_THAN_EQT,
							Right:    &ast.IntLiteralNode{Value: 6},
						},
					},
					&ast.LetStmtNode{
						Name: "y",
						Value: &ast.BoolInfixNode{
							Left: &ast.PrefixExprNode{
								Operator: token.NOT,
								Value:    &ast.ReferenceExprNode{Name: "x"},
							},
							Operator: token.OR,
							Right:    &ast.BoolLiteralNode{Value: true},
						},
					},
				},
			},
			id: 9,
		},
		{
			input: "if true{let y = 4;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 4},
							},
						},
					},
				},
			},
			id: 10,
		},
		{
			input: "let y = 4; if y < 9{let z = 5;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "y",
						Value: &ast.IntLiteralNode{Value: 4},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left:     &ast.ReferenceExprNode{Name: "y"},
							Operator: token.LESS_THAN,
							Right:    &ast.IntLiteralNode{Value: 9},
						},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "z",
								Value: &ast.IntLiteralNode{Value: 5},
							},
						},
					},
				},
			},
			id: 11,
		},
		{
			input: "let x = 1; if true { if false { x += 5; } }",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.IntLiteralNode{
							Value: 1,
						},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.IfStmtNode{
								Cond: &ast.BoolLiteralNode{Value: false},
								Body: []ast.Node{
									&ast.VarReassignNode{
										Var: ast.ReferenceExprNode{Name: "x"},
										NewVal: &ast.InfixExprNode{
											Left:     &ast.ReferenceExprNode{Name: "x"},
											Operator: token.PLUS,
											Right:    &ast.IntLiteralNode{Value: 5},
										},
									},
								},
							},
						},
					},
				},
			},
			id: 12,
		},
		{
			input: "let x = false; if !x&&true{let y = !x;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "x",
						Value: &ast.BoolLiteralNode{Value: false},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left: &ast.PrefixExprNode{
								Operator: token.NOT,
								Value:    &ast.ReferenceExprNode{Name: "x"},
							},
							Operator: token.AND,
							Right:    &ast.BoolLiteralNode{Value: true},
						},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name: "y",
								Value: &ast.PrefixExprNode{
									Operator: token.NOT,
									Value:    &ast.ReferenceExprNode{Name: "x"},
								},
							},
						},
					},
				},
			},
			id: 13,
		},
		{
			input: "let x = 9; if x < 10{let y = 4;} else {let y = 5;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "x",
						Value: &ast.IntLiteralNode{Value: 9},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left:     &ast.ReferenceExprNode{Name: "x"},
							Operator: token.LESS_THAN,
							Right:    &ast.IntLiteralNode{Value: 10},
						},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 4},
							},
						},
						Alt: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 5},
							},
						},
					},
				},
			},
			id: 14,
		},
		{
			input: "if true { let x = 1; } else { let x = 3 >= 4; }",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "x",
								Value: &ast.IntLiteralNode{Value: 1},
							},
						},
						Alt: []ast.Node{
							&ast.LetStmtNode{
								Name: "x",
								Value: &ast.BoolInfixNode{
									Left:     &ast.IntLiteralNode{Value: 3},
									Operator: token.GREATER_THAN_EQT,
									Right:    &ast.IntLiteralNode{Value: 4},
								},
							},
						},
					},
				},
			},
			id: 15,
		},
		{
			input: "let v = true || false; if v{v = false;} else {v = true;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "v",
						Value: &ast.BoolInfixNode{
							Left:     &ast.BoolLiteralNode{Value: true},
							Operator: token.OR,
							Right:    &ast.BoolLiteralNode{Value: false},
						},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left:     &ast.ReferenceExprNode{Name: "v"},
							Operator: token.OR,
							Right:    &ast.BoolLiteralNode{Value: false},
						},
						Body: []ast.Node{
							&ast.VarReassignNode{
								Var:    ast.ReferenceExprNode{Name: "v"},
								NewVal: &ast.BoolLiteralNode{Value: false},
							},
						},
						Alt: []ast.Node{
							&ast.VarReassignNode{
								Var:    ast.ReferenceExprNode{Name: "v"},
								NewVal: &ast.BoolLiteralNode{Value: true},
							},
						},
					},
				},
			},
			id: 16,
		},
		{
			input: "let x = 9;if true{let y = 4; x = y;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "x",
						Value: &ast.IntLiteralNode{Value: 9},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 4},
							},
							&ast.VarReassignNode{
								Var:    ast.ReferenceExprNode{Name: "x"},
								NewVal: &ast.ReferenceExprNode{Name: "y"},
							},
						},
					},
				},
			},
			id: 17,
		},
		{
			input: "let x = 0; if true {let y = 4; x = y;} else {let y = 5; x = y;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "x",
						Value: &ast.IntLiteralNode{Value: 0},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 4},
							},
							&ast.VarReassignNode{
								Var:    ast.ReferenceExprNode{Name: "x"},
								NewVal: &ast.ReferenceExprNode{Name: "y"},
							},
						},
						Alt: []ast.Node{
							&ast.LetStmtNode{
								Name:  "y",
								Value: &ast.IntLiteralNode{Value: 5},
							},
							&ast.VarReassignNode{
								Var:    ast.ReferenceExprNode{Name: "x"},
								NewVal: &ast.ReferenceExprNode{Name: "y"},
							},
						},
					},
				},
			},
			id: 18,
		},
		{
			input: "let x = 4 * (8 + 3); let x = 9 / (x + 1);",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.InfixExprNode{
							Left:     &ast.IntLiteralNode{Value: 4},
							Operator: token.MULTIPLY,
							Right: &ast.EmptyExprNode{
								Child: &ast.InfixExprNode{
									Left:     &ast.IntLiteralNode{Value: 8},
									Operator: token.PLUS,
									Right:    &ast.IntLiteralNode{Value: 3},
								},
							},
						},
					},
					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.InfixExprNode{
							Left:     &ast.IntLiteralNode{Value: 9},
							Operator: token.DIVIDE,
							Right: &ast.EmptyExprNode{
								Child: &ast.InfixExprNode{
									Left:     &ast.ReferenceExprNode{Name: "x"},
									Operator: token.PLUS,
									Right:    &ast.IntLiteralNode{Value: 1},
								},
							},
						},
					},
				},
			},
			id: 19,
		},
		{
			input: "if true{if !false{let y = 4;}}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.IfStmtNode{
						Cond: &ast.BoolLiteralNode{Value: true},
						Body: []ast.Node{
							&ast.IfStmtNode{
								Cond: &ast.PrefixExprNode{
									Operator: token.NOT,
									Value:    &ast.BoolLiteralNode{Value: false},
								},
								Body: []ast.Node{
									&ast.LetStmtNode{
										Name:  "y",
										Value: &ast.IntLiteralNode{Value: 4},
									},
								},
							},
						},
					},
				},
			},
			id: 20,
		},
		{
			input: "fn a(b){return b + 3;} let c = a(4);",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.FuncDecNode{
						Name:   "a",
						Params: []ast.ReferenceExprNode{ast.ReferenceExprNode{Name: "b"}},
						Body:   []ast.Node{},
						Return: ast.ReturnExprNode{
							Val: &ast.InfixExprNode{
								Left:     &ast.ReferenceExprNode{Name: "b"},
								Operator: token.PLUS,
								Right:    &ast.IntLiteralNode{Value: 3},
							},
						},
					},
					&ast.LetStmtNode{
						Name: "c",
						Value: &ast.FuncCallNode{
							Name: ast.ReferenceExprNode{Name: "a"},
							Params: []ast.Node{
								&ast.IntLiteralNode{Value: 4},
							},
						},
					},
				},
			},
			id: 21,
		},
		{
			input: "fn a(b){return b - 2;} fn c(b){return b + 2;} let d = a(2) + c(2);",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.FuncDecNode{
						Name:   "a",
						Params: []ast.ReferenceExprNode{ast.ReferenceExprNode{Name: "b"}},
						Body:   []ast.Node{},
						Return: ast.ReturnExprNode{
							Val: &ast.InfixExprNode{
								Left:     &ast.ReferenceExprNode{Name: "b"},
								Operator: token.MINUS,
								Right:    &ast.IntLiteralNode{Value: 2},
							},
						},
					},
					&ast.FuncDecNode{
						Name:   "c",
						Params: []ast.ReferenceExprNode{ast.ReferenceExprNode{Name: "b"}},
						Body:   []ast.Node{},
						Return: ast.ReturnExprNode{
							Val: &ast.InfixExprNode{
								Left:     &ast.ReferenceExprNode{Name: "b"},
								Operator: token.PLUS,
								Right:    &ast.IntLiteralNode{Value: 2},
							},
						},
					},
					&ast.LetStmtNode{
						Name: "d",
						Value: &ast.InfixExprNode{
							Left: &ast.FuncCallNode{
								Name:   ast.ReferenceExprNode{Name: "a"},
								Params: []ast.Node{&ast.IntLiteralNode{Value: 2}},
							},
							Operator: token.PLUS,
							Right: &ast.FuncCallNode{
								Name:   ast.ReferenceExprNode{Name: "c"},
								Params: []ast.Node{&ast.IntLiteralNode{Value: 2}},
							},
						},
					},
				},
			},
			id: 22,
		},
		{
			input: "fn add(a, b){return a + b;} let c = add(add(1, 1), 2);",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.FuncDecNode{
						Name: "add",
						Params: []ast.ReferenceExprNode{
							ast.ReferenceExprNode{Name: "a"},
							ast.ReferenceExprNode{Name: "b"},
						},
						Body: []ast.Node{
							&ast.ReturnExprNode{
								Val: &ast.InfixExprNode{
									Left:     &ast.ReferenceExprNode{Name: "a"},
									Operator: token.PLUS,
									Right:    &ast.ReferenceExprNode{Name: "b"},
								},
							},
						},
					},
					&ast.LetStmtNode{
						Name: "c",
						Value: &ast.FuncCallNode{
							Name: ast.ReferenceExprNode{Name: "add"},
							Params: []ast.Node{
								&ast.FuncCallNode{
									Name: ast.ReferenceExprNode{Name: "add"},
									Params: []ast.Node{
										&ast.IntLiteralNode{Value: 1},
										&ast.IntLiteralNode{Value: 1},
									},
								},
								&ast.IntLiteralNode{Value: 2},
							},
						},
					},
				},
			},
			id: 23,
		},
		{
			input: `let x = "hi";`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name:  "x",
						Value: &ast.StringLiteralNode{Value: "hi"},
					},
				},
			},
			id: 24,
		},
		{
			input: `let x = "hello" + "world";`,
			output: ast.ProgramNode{
				Statements: []ast.Node{

					&ast.LetStmtNode{
						Name: "x",
						Value: &ast.InfixExprNode{
							Left:     &ast.StringLiteralNode{Value: "hello"},
							Operator: token.PLUS,
							Right:    &ast.StringLiteralNode{Value: "world"},
						},
					},
				},
			},
			id: 25,
		},
		{
			input: `let hello = "hello " + "world";`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "hello",
						Value: &ast.InfixExprNode{
							Left:     &ast.StringLiteralNode{Value: "hello "},
							Operator: token.PLUS,
							Right:    &ast.StringLiteralNode{Value: "world"},
						},
					},
				},
			},
			id: 26,
		},
		{
			input: `fn concat(a, b){return a + b;} let hello = concat("hello ", "world");`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.FuncDecNode{
						Name: "concat",
						Params: []ast.ReferenceExprNode{
							ast.ReferenceExprNode{Name: "a"},
							ast.ReferenceExprNode{Name: "b"},
						},
						Body: []ast.Node{
							&ast.ReturnExprNode{
								Val: &ast.InfixExprNode{
									Left:     &ast.ReferenceExprNode{Name: "a"},
									Operator: token.PLUS,
									Right:    &ast.ReferenceExprNode{Name: "b"},
								},
							},
						},
					},
					&ast.LetStmtNode{
						Name: "hello",
						Value: &ast.FuncCallNode{
							Name: ast.ReferenceExprNode{Name: "concat"},
							Params: []ast.Node{
								&ast.StringLiteralNode{Value: "hello "},
								&ast.StringLiteralNode{Value: "world"},
							},
						},
					},
				},
			},
			id: 27,
		},
		{
			input: `if "a" < "b"{let h = "hi";}`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left:     &ast.StringLiteralNode{Value: "a"},
							Operator: token.LESS_THAN,
							Right:    &ast.StringLiteralNode{Value: "b"},
						},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "h",
								Value: &ast.StringLiteralNode{Value: "hi"},
							},
						},
					},
				},
			},
			id: 28,
		},
		{
			input: `if "hello" == "world"{let val = true;} else {let val = false;}`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left:     &ast.StringLiteralNode{Value: "hello"},
							Operator: token.EQUALS,
							Right:    &ast.StringLiteralNode{Value: "world"},
						},
						Body: []ast.Node{
							&ast.LetStmtNode{
								Name:  "val",
								Value: &ast.BoolLiteralNode{Value: true},
							},
						},
						Alt: []ast.Node{
							&ast.LetStmtNode{
								Name:  "val",
								Value: &ast.BoolLiteralNode{Value: false},
							},
						},
					},
				},
			},
			id: 29,
		},
		{
			input: `print("hello world");`,
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.FuncCallNode{
						Name:   ast.ReferenceExprNode{Name: "print"},
						Params: []ast.Node{&ast.StringLiteralNode{Value: "hello world"}},
					},
				},
			},
			id: 30,
		},
	}

	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := NewParser()
		toks := lex.Lex(tt.input)
		prog := parse.Parse(toks)

		compareNodes(t, prog.Statements, tt.output.Statements, tt)
	}
}
