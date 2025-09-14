package parser

import (
	"testing"
	"toy_lang/ast"
	"toy_lang/lexer"
	"toy_lang/token"
)

// compareNodes compares slices of AST nodes and prints mismatches
func compareNodes(t *testing.T, got, want []ast.Node) {
	if len(got) != len(want) {
		t.Errorf("Length mismatch: got %d, want %d", len(got), len(want))
	}

	minLen := len(got)
	if len(want) < minLen {
		minLen = len(want)
	}

	for i := 0; i < minLen; i++ {
		if got[i].String() != want[i].String() {
			t.Errorf("Mismatch at index %d:\n Got:  %v\n Want: %v", i, got[i], want[i])
		}
	}

	if len(got) > len(want) {
		for i := len(want); i < len(got); i++ {
			t.Errorf("Extra element in got at index %d: %v", i, got[i])
		}
	} else if len(want) > len(got) {
		for i := len(got); i < len(want); i++ {
			t.Errorf("Missing element in got at index %d: %v", i, want[i])
		}
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
	}

	for _, tt := range tests {
		res := parse.preProcess(lex.Lex(tt.input))
		compareTokens(t, res, tt.output)
	}
}

func TestParser(t *testing.T) {

	tests := []struct {
		input  string
		output ast.ProgramNode
	}{
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
									Left: &ast.IntLiteralNode{Value: 3},
									Operator: token.GREATER_THAN_EQT,
									Right: &ast.IntLiteralNode{Value: 4},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "let v = true || false; if v{v = false;} else {v = true;}",
			output: ast.ProgramNode{
				Statements: []ast.Node{
					&ast.LetStmtNode{
						Name: "v",
						Value: &ast.BoolInfixNode{
							Left: &ast.BoolLiteralNode{Value: true},
							Operator: token.OR,
							Right: &ast.BoolLiteralNode{Value: false},
						},
					},
					&ast.IfStmtNode{
						Cond: &ast.BoolInfixNode{
							Left: &ast.ReferenceExprNode{Name: "v"},
							Operator: token.OR,
							Right: &ast.BoolLiteralNode{Value: false},
						},
						Body: []ast.Node{
							&ast.VarReassignNode{
								Var: ast.ReferenceExprNode{Name: "v"},
								NewVal: &ast.BoolLiteralNode{Value: false},
							},
						},
						Alt: []ast.Node{
							&ast.VarReassignNode{
								Var: ast.ReferenceExprNode{Name: "v"},
								NewVal: &ast.BoolLiteralNode{Value: true},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		lex := lexer.NewLexer()
		parse := NewParser()
		toks := lex.Lex(tt.input)
		prog := parse.Parse(toks)

		compareNodes(t, prog.Statements, tt.output.Statements)
	}
}
