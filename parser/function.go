package parser

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

func (p *Parser) parseFuncDecStmt(toks []token.Token) *ast.FuncDecNode {
	if toks[0].TokType != token.FN {
		panic(fmt.Sprintf("[ERROR] Must use fn to declare function, got %v\n", toks[0]))
	}
	if toks[1].TokType != token.FUNC_NAME {
		panic(fmt.Sprintf("[ERROR] Could not figure out function name, got %v\n", toks[1]))
	}
	if toks[2].TokType != token.LPAREN {
		panic(fmt.Sprintf("[ERROR] Function name must be followed by \"(\", got %v\n", toks[2]))
	}

	// ----------- Collect params -----------
	params := []token.Token{}
	i := 3
	for i < len(toks) && toks[i].TokType != token.RPAREN {
		if toks[i].TokType != token.COMMA {
			params = append(params, toks[i])
		}
		i++
	}
	if i >= len(toks) || toks[i].TokType != token.RPAREN {
		panic("[ERROR] Unmatched parenthesis in function declaration")
	}

	var astParams []ast.ReferenceExprNode
	for _, val := range params {
		astParams = append(astParams, ast.ReferenceExprNode{Name: val.Literal})
	}

	// ----------- Collect body -----------
	if toks[i+1].TokType != token.LBRACE {
		panic(fmt.Sprintf("[ERROR] Expected { after params, got %v", toks[i+1]))
	}

	// find matching }
	depth := 1
	j := i + 2
	for j < len(toks) && depth > 0 {
		if toks[j].TokType == token.LBRACE {
			depth++
		} else if toks[j].TokType == token.RBRACE {
			depth--
		}
		j++
	}
	if depth != 0 {
		panic("[ERROR] Mismatched braces in function body")
	}

	bodyTokens := toks[i+2 : j-1]
	var body []ast.Node
	for _, line := range p.splitIntoLines(bodyTokens) {
		body = append(body, p.parseStmt(line))
	}

	return &ast.FuncDecNode{
		Name:   toks[1].Literal,
		Params: astParams,
		Body:   body,
	}
}

func (p *Parser) parseReturnExpr(toks []token.Token) *ast.ReturnExprNode {
	if toks[0].TokType != token.RETURN {
		panic(fmt.Sprintf("[ERROR] Return statement must start with return, got %v\n", toks[0]))
	}
	return &ast.ReturnExprNode{
		Val: p.parseStmt(toks[1:]),
	}
}

func (p *Parser) parseFuncCallStmt(toks []token.Token) *ast.FuncCallNode {
	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Could not figure out function name, got %v\n", toks[0]))
	}
	if toks[1].TokType != token.LPAREN {
		panic(fmt.Sprintf("[ERROR] Function name must be followed by \"(\", got %v\n", toks[1]))
	}

	funcName := toks[0].Literal

	// find matching RPAREN
	depth := 1
	j := 2
	for j < len(toks) && depth > 0 {
		if toks[j].TokType == token.LPAREN {
			depth++
		} else if toks[j].TokType == token.RPAREN {
			depth--
		}
		j++
	}
	if depth != 0 {
		panic("[ERROR] Mismatched parentheses in function call")
	}

	argTokens := toks[2 : j-1]

	// split args by comma
	var args [][]token.Token
	curr := []token.Token{}
	for _, tok := range argTokens {
		if tok.TokType == token.COMMA {
			if len(curr) > 0 {
				args = append(args, curr)
				curr = []token.Token{}
			}
		} else {
			curr = append(curr, tok)
		}
	}
	if len(curr) > 0 {
		args = append(args, curr)
	}

	// lookup function declaration in program
	var funcDecl *ast.FuncDecNode
	for _, stmt := range p.program.Statements {
		if f, ok := stmt.(*ast.FuncDecNode); ok && f.Name == funcName {
			funcDecl = f
			break
		}
	}

	var params []ast.Node
	for i, group := range args {
		if len(group) > 0 {
			expr := p.parseExpression(group)
			paramName := fmt.Sprintf("arg%d", i) // fallback name
			if funcDecl != nil && i < len(funcDecl.Params) {
				paramName = funcDecl.Params[i].Name
			}
			params = append(params, &ast.LetStmtNode{
				Name:  paramName,
				Value: expr,
			})
		}
	}

	return &ast.FuncCallNode{
		Name:   ast.ReferenceExprNode{Name: funcName},
		Params: params,
	}
}
