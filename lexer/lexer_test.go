package lexer

import (
	"slices"
	"testing"
	"toy_lang/token"
)

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
func TestLexer(t *testing.T) {
	lex := NewLexer()
	tests := []struct {
		input  string
		output []token.Token
	}{
		{
			input: "let x = 4;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let x = 7; let y = 9 + x; x++;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "7"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "9"),
				*token.NewToken(token.PLUS, "+"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.PLUS_PLUS, "++"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let x = 3; x+=1; x-=3; x*=4; x/=6; x--;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "3"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.COMPOUND_PLUS, "+="),
				*token.NewToken(token.INTEGER, "1"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.COMPOUND_MINUS, "-="),
				*token.NewToken(token.INTEGER, "3"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.COMPOUND_MULTIPLY, "*="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.COMPOUND_DIVIDE, "/="),
				*token.NewToken(token.INTEGER, "6"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.MINUS_MINUS, "--"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		//Boolean literal tests
		{
			input: "let b = true; let y = false;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "b"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		//Boolean Operator tests
		{
			input: "let b = 5 < 6; let a = b && true;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "b"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.LESS_THAN, "<"),
				*token.NewToken(token.INTEGER, "6"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "a"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "b"),
				*token.NewToken(token.AND, "&&"),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let b = true || false; let c = !b;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "b"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.OR, "||"),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "c"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.NOT, "!"),
				*token.NewToken(token.VAR_REF, "b"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "let x = 4 >= 5; let y = 5 <= 9;",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.GREATER_THAN_EQT, ">="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.LESS_THAN_EQT, "<="),
				*token.NewToken(token.INTEGER, "9"),
				*token.NewToken(token.SEMICOLON, ";"),
			},
		},
		{
			input: "if 3==4{let x = 5;}",
			output: []token.Token{
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.INTEGER, "3"),
				*token.NewToken(token.EQUALS, "=="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let x = 9; if x < 5{let y = 4;}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "9"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.LESS_THAN, "<"),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let y = true || false; if y{let x = 5;}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.OR, "||"),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.VAR_REF, "y"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let x = true && !false; if x{if false{let y = 4;}}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.AND, "&&"),
				*token.NewToken(token.NOT, "!"),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let x = false; if !x&&true{let y = !x;}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.NOT, "!"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.AND, "&&"),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.NOT, "!"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let x= 5; if x <=6{let y = 9;}else{let z = 5;}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.LESS_THAN_EQT, "<="),
				*token.NewToken(token.INTEGER, "6"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "9"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
				*token.NewToken(token.ELSE, "else"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "z"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "if false{let y = 9;}else{let z = 3;}",
			output: []token.Token{
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.BOOLEAN, "false"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "9"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
				*token.NewToken(token.ELSE, "else"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "z"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "3"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
		{
			input: "let x = 0; if true {let y = 4; x = y;} else {let y = 5; x = y;}",
			output: []token.Token{
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "0"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.IF, "if"),
				*token.NewToken(token.BOOLEAN, "true"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "4"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "y"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
				*token.NewToken(token.ELSE, "else"),
				*token.NewToken(token.LBRACE, "{"),
				*token.NewToken(token.LET, "let"),
				*token.NewToken(token.VAR_NAME, "y"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.INTEGER, "5"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.VAR_REF, "x"),
				*token.NewToken(token.ASSIGN, "="),
				*token.NewToken(token.VAR_REF, "y"),
				*token.NewToken(token.SEMICOLON, ";"),
				*token.NewToken(token.RBRACE, "}"),
			},
		},
	}
	for _, tt := range tests {
		res := lex.Lex(tt.input)
		if !slices.Equal(res, tt.output) {
			res := lex.Lex(tt.input)
			compareTokens(t, res, tt.output)
		}
	}
}
