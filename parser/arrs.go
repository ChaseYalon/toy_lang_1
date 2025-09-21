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

	arrElems := make(map[string]ast.Node)
	for i, val := range valNodes {
		temp := ast.IntLiteralNode{Value: i};
		arrElems[temp.String()] = val
	}

	return &ast.ArrLiteralNode{
		Elems: arrElems,
	}
}
func (p *Parser) parseArrRef(toks []token.Token) ast.Node {
	if len(toks) < 3 || toks[0].TokType != token.VAR_REF || toks[1].TokType != token.LBRACK {
		panic(fmt.Sprintf("[ERROR] Invalid array reference, got %v\n", toks))
	}

	// check if this is an assignment like arr[2] = 4
	if includes, assignIdx := includesItem(toks, *token.NewToken(token.ASSIGN, "=")); includes {
		if toks[assignIdx-1].TokType != token.RBRACK {
			panic(fmt.Sprintf("[ERROR] Array index must end with ']', got %v\n", toks))
		}
		name := p.parseVarReference(toks[0])
		// index tokens: between [ and ]
		idxTokens := toks[2 : assignIdx-1]
		idxNode := p.parseExpression(idxTokens)
		newValNode := p.parseExpression(toks[assignIdx+1:])
		return &ast.ArrReassignNode{
			Arr:    *name,
			Idx:    idxNode,
			NewVal: newValNode,
		}
	}

	// normal array reference like arr[2]
	if toks[len(toks)-1].TokType != token.RBRACK {
		panic(fmt.Sprintf("[ERROR] Array index must end with ']', got %v\n", toks))
	}
	name := p.parseVarReference(toks[0])
	idxNode := p.parseExpression(toks[2 : len(toks)-1])
	return &ast.ArrRefNode{
		Arr: *name,
		Idx: idxNode,
	}
}
