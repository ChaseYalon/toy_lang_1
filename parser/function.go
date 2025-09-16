package parser;

import (
	"toy_lang/token"
	"toy_lang/ast"
	"fmt"
)

func (p *Parser) parseFuncDecStmt(toks []token.Token) ast.FuncDecNode{
	if toks[0].TokType != token.FN{
		panic(fmt.Sprintf("[ERROR] Must use fn to declare function, got %v\n", toks[0]));
	}
	
}

func (p *Parser) parseFuncCallStmt(toks []token.Token) ast.FuncCallNode{

}