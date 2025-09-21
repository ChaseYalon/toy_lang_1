package parser

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

func (p *Parser) parseArr(toks []token.Token) ast.Node {
	if len(toks) < 2 {
		panic(fmt.Sprintf("[ERROR] Array must contain at least \"[\" and \"]\" got %+v\n", toks))
	}
	if toks[0].TokType != token.LBRACK {
		panic(fmt.Sprintf("[ERROR] Array literal must start with an opening bracket, got %v\n", toks[0]))
	}
	if toks[len(toks)-1].TokType != token.RBRACK {
		panic(fmt.Sprintf("[ERROR] Array must end with closing bracket, got %v\n", toks[len(toks)-1]))
	}

	var vals [][]token.Token
	var currVal []token.Token
	for _, val := range toks[1 : len(toks)-1] {
		if val.TokType == token.COMMA {
			vals = append(vals, currVal)
			currVal = []token.Token{}
			continue
		}
		currVal = append(currVal, val)
	}
	// append last value
	if len(currVal) > 0 {
		vals = append(vals, currVal)
	}

	var valNodes []ast.Node
	for _, val := range vals {
		valNodes = append(valNodes, p.parseExpression(val))
	}

	arrElems := make(map[ast.Node]ast.Node)
	for i, val := range valNodes {
		arrElems[&ast.IntLiteralNode{Value: i}] = val
	}

	return &ast.ArrLiteralNode{
		Elems: arrElems,
	}
}

func (p *Parser) parseArrRef(toks []token.Token) ast.ArrRefNode {
	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Could not figure out what array to reference, got %v\n", toks))
	}
	if toks[1].TokType != token.LBRACK {
		panic(fmt.Sprintf("[ERROR] Var name must be followed by \"[\" got %v\n", toks))
	}
	if toks[len(toks)-1].TokType != token.RBRACK {
		panic(fmt.Sprintf("[ERROR] Could not figure out where array index ends, it must end with \"]\", got %v\n", toks))
	}

	name := p.parseVarReference(toks[0])
	val := p.parseExpression(toks[2 : len(toks)-1])
	return ast.ArrRefNode{
		Arr: *name,
		Idx: val,
	}
}
