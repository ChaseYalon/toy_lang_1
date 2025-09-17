package parser;

import (
	"toy_lang/token"
	"toy_lang/ast"
	"fmt"
)

func (p *Parser) parseFuncDecStmt(toks []token.Token) *ast.FuncDecNode{
	if toks[0].TokType != token.FN{
		panic(fmt.Sprintf("[ERROR] Must use fn to declare function, got %v\n", toks[0]));
	}
	if toks[1].TokType != token.FUNC_NAME{
		panic(fmt.Sprintf("[ERROR] Could not figure out function name, got %v\n", toks[1]));
	}
	if toks[2].TokType != token.LPAREN{
		panic(fmt.Sprintf("[ERROR] Could Function name must be followed by \"(\", got %v\n", toks[2]));
	}
	var params []token.Token;
	bodyStartAt := -1;
	for i, val := range toks[3:]{
		if val.TokType == token.RPAREN{
			bodyStartAt = i + 4; //Account for "{" and first 3
			break;
		}
		if val.TokType == token.COMMA{
			continue;
		}
		params = append(params, val);
	}
	var astPrams []ast.ReferenceExprNode;
	for _, val := range params{
		astPrams = append(astPrams, ast.ReferenceExprNode{Name: val.Literal});
	}
	var body []ast.Node;
	splitToks := p.splitIntoLines(toks[bodyStartAt : len(toks) - 2]);
	for _, val := range splitToks{
		body = append(body, p.parseStmt(val));
	}
	protoFuncNode := ast.FuncDecNode{
		Name: toks[1].Literal,
		Params: astPrams,
		Body: body,
		
	}
	return &protoFuncNode
}
func (p *Parser) parseReturnExpr(toks []token.Token) *ast.ReturnExprNode{
	if toks[0].TokType != token.RETURN{
		panic(fmt.Sprintf("[ERROR] Return statement must start with return, got %v\n", toks[0]));
	}
	protoReturnExpr := &ast.ReturnExprNode{
		Val: p.parseStmt(toks[1:]),
	}
	return protoReturnExpr;
}
func (p *Parser) parseFuncCallStmt(toks []token.Token) ast.FuncCallNode{

}