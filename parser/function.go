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
	if len(toks) == 0 {
		panic("[ERROR] Empty tokens to parseFuncCallStmt")
	}
	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Could not figure out function name, got %v\n", toks[0]))
	}
	if len(toks) < 2 || toks[1].TokType != token.LPAREN {
		panic(fmt.Sprintf("[ERROR] Function name must be followed by \"(\", got %v\n", toks[1]))
	}

	funcName := toks[0].Literal

	// match closing RPAREN
	depth := 1
	j := 2
	for j < len(toks) && depth > 0 {
		if toks[j].TokType == token.LPAREN {
			depth++
		} else if toks[j].TokType == token.RPAREN {
			depth--
			if depth == 0 {
				break
			}
		}
		j++
	}
	if depth != 0 || j >= len(toks) {
		panic("[ERROR] Mismatched parentheses in function call")
	}

	argTokens := toks[2:j]

	// split args by top-level commas
	var args [][]token.Token
	curr := []token.Token{}
	depth = 0
	for _, tk := range argTokens {
		if tk.TokType == token.LPAREN {
			depth++
		} else if tk.TokType == token.RPAREN {
			depth--
		}
		if tk.TokType == token.COMMA && depth == 0 {
			if len(curr) > 0 {
				args = append(args, curr)
				curr = []token.Token{}
			}
			continue
		}
		curr = append(curr, tk)
	}
	if len(curr) > 0 {
		args = append(args, curr)
	}

	var params []ast.Node
	for _, group := range args {
		if len(group) == 0 {
			continue
		}
		val := p.parseExpression(group)
		params = append(params, val)
	}

	return &ast.FuncCallNode{
		Name:   ast.ReferenceExprNode{Name: funcName},
		Params: params,
	}
}
