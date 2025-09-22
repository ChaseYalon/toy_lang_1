package parser

import (
	"fmt"
	"strconv"
	"toy_lang/ast"
	"toy_lang/token"
)

type Parser struct {
	program ast.ProgramNode
	tokens  []token.Token
	ifStack []*ast.IfStmtNode
}

func NewParser() *Parser {
	return &Parser{
		program: ast.ProgramNode{
			Statements: []ast.Node{},
		},
		tokens: []token.Token{},
		ifStack: []*ast.IfStmtNode{},
	}
}

func (p *Parser) parseExpression(tokens []token.Token) ast.Node {
	if len(tokens) == 0 {
		panic("[ERROR] Empty expression")
	}

	var newTokens []token.Token
	var subNodes []*ast.EmptyExprNode

	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		if tok.TokType == token.VAR_REF && i+1 < len(tokens) && tokens[i+1].TokType == token.LPAREN {
			depth := 1
			j := i + 2
			for j < len(tokens) && depth > 0 {
				if tokens[j].TokType == token.LPAREN {
					depth++
				} else if tokens[j].TokType == token.RPAREN {
					depth--
				}
				j++
			}
			if depth != 0 {
				panic("[ERROR] Mismatched parentheses in function call")
			}

			funcCall := p.parseFuncCallStmt(tokens[i:j])
			emptyNode := &ast.EmptyExprNode{Child: funcCall}

			newTokens = append(newTokens, *token.NewToken(token.EMPTY, ""))
			subNodes = append(subNodes, emptyNode)

			i = j
		} else if tok.TokType == token.LPAREN {
			depth := 1
			j := i + 1
			for j < len(tokens) && depth > 0 {
				if tokens[j].TokType == token.LPAREN {
					depth++
				} else if tokens[j].TokType == token.RPAREN {
					depth--
				}
				j++
			}
			if depth != 0 {
				panic("[ERROR] Mismatched parentheses")
			}

			sub := p.parseExpression(tokens[i+1 : j-1])
			emptyNode := &ast.EmptyExprNode{Child: sub}

			newTokens = append(newTokens, *token.NewToken(token.EMPTY, ""))
			subNodes = append(subNodes, emptyNode)

			i = j
		} else {
			newTokens = append(newTokens, tok)
			i++
		}
	}

	if tokens[0].TokType == token.LBRACK {
		depth := 1
		var toks []token.Token
		toks = append(toks, tokens[0])
		for i := 1; i < len(tokens); i++ {
			toks = append(toks, tokens[i])
			if tokens[i].TokType == token.LBRACK {
				depth++
			} else if tokens[i].TokType == token.RBRACK {
				depth--
				if depth == 0 {
					if i == len(tokens)-1 {
						return p.parseArr(toks)
					}
					break
				}
			}
		}
		if depth != 0 {
			panic(fmt.Sprintf("[ERROR] Could not find right brace corresponding to left brace, got %v\n", tokens))
		}
	}

	contains, _ := includesItem(tokens, *token.NewToken(token.ASSIGN, "="))
	if contains {
		return p.parseSubExpression(tokens, subNodes)
	}

	return p.parseSubExpression(newTokens, subNodes)
}

func (p *Parser) parseSubExpression(tokens []token.Token, subNodes []*ast.EmptyExprNode) ast.Node {
	if len(tokens) == 0 {
		panic("[ERROR] Empty expression")
	}

	if len(tokens) >= 4 && tokens[0].TokType == token.VAR_REF && tokens[1].TokType == token.LBRACK {
		if includes, _ := includesItem(tokens, *token.NewToken(token.ASSIGN, "=")); includes {
			return p.parseArrRef(tokens)
		}
	}

	if len(tokens) == 1 {
		tok := tokens[0]
		switch tok.TokType {
		case token.INTEGER:
			val, _ := strconv.Atoi(tok.Literal)
			return &ast.IntLiteralNode{Value: val}
		case token.BOOLEAN:
			val, _ := strconv.ParseBool(tok.Literal)
			return &ast.BoolLiteralNode{Value: val}
		case token.STRING:
			return &ast.StringLiteralNode{Value: tok.Literal}
		case token.FLOAT:
			val, err := strconv.ParseFloat(tok.Literal, 64)
			if err != nil {
				panic(fmt.Sprintf("[ERROR] Could not convert to floating point, got %v\n", val))
			}
			return &ast.FloatLiteralNode{Value: val}
		case token.VAR_REF:
			if len(tokens) != 1 {
				if tokens[1].TokType == token.LPAREN {
					return p.parseFuncCallStmt(tokens)
				}
			}
			return &ast.ReferenceExprNode{Name: tok.Literal}
		case token.EMPTY:
			if len(subNodes) == 0 || subNodes[0] == nil {
				panic("[ERROR] EMPTY token without corresponding subnode")
			}
			return subNodes[0]
		default:
			panic(fmt.Sprintf("[ERROR] Unexpected single token: %+v", tok))
		}
	}
	hasOperator, _ := includesAny(tokens, []token.TokenType{
		token.PLUS,
		token.MINUS,
		token.MULTIPLY,
		token.DIVIDE,
		token.MODULO,
		token.EXPONENT,
		token.LESS_THAN,
		token.LESS_THAN_EQT,
		token.GREATER_THAN,
		token.GREATER_THAN_EQT,
		token.EQUALS,
		token.NOT_EQUAL,
		token.AND,
		token.OR,
		token.NOT,
	})
	if tokens[0].TokType == token.VAR_REF && tokens[1].TokType == token.LBRACK && !hasOperator {
		arr := p.parseArrRef(tokens)
		return arr
	}
	lowestTok, lowestIndex := p.findLowestBp(p.generatePrecedenceTable(), tokens)
	if lowestIndex == -1 {
		var leftSubNodes []*ast.EmptyExprNode
		if len(subNodes) > 0 {
			leftSubNodes = subNodes[:1]
		}
		return p.parseSubExpression(tokens[:1], leftSubNodes)
	}

	sliceSubNodes := func(start, end int) []*ast.EmptyExprNode {
		if len(subNodes) == 0 {
			return []*ast.EmptyExprNode{}
		}

		emptyCount := 0
		for i := start; i < end && i < len(tokens); i++ {
			if tokens[i].TokType == token.EMPTY {
				emptyCount++
			}
		}

		if emptyCount == 0 {
			return []*ast.EmptyExprNode{}
		}

		subNodesStart := 0
		for i := 0; i < start && i < len(tokens); i++ {
			if tokens[i].TokType == token.EMPTY {
				subNodesStart++
			}
		}

		subNodesEnd := subNodesStart + emptyCount
		if subNodesEnd > len(subNodes) {
			subNodesEnd = len(subNodes)
		}

		return subNodes[subNodesStart:subNodesEnd]
	}

	switch lowestTok.TokType {
	case token.AND, token.OR,
		token.LESS_THAN, token.LESS_THAN_EQT,
		token.GREATER_THAN, token.GREATER_THAN_EQT,
		token.EQUALS:
		leftSubNodes := sliceSubNodes(0, lowestIndex)
		rightSubNodes := sliceSubNodes(lowestIndex+1, len(tokens))

		left := p.parseSubExpression(tokens[:lowestIndex], leftSubNodes)
		right := p.parseSubExpression(tokens[lowestIndex+1:], rightSubNodes)
		return &ast.BoolInfixNode{
			Left:     left,
			Operator: lowestTok.TokType,
			Right:    right,
		}
	case token.NOT:
		rightSubNodes := sliceSubNodes(lowestIndex+1, len(tokens))
		right := p.parseSubExpression(tokens[lowestIndex+1:], rightSubNodes)
		return &ast.PrefixExprNode{
			Value:    right,
			Operator: token.NOT,
		}
	default:
		leftSubNodes := sliceSubNodes(0, lowestIndex)
		rightSubNodes := sliceSubNodes(lowestIndex+1, len(tokens))

		left := p.parseSubExpression(tokens[:lowestIndex], leftSubNodes)
		right := p.parseSubExpression(tokens[lowestIndex+1:], rightSubNodes)
		return &ast.InfixExprNode{
			Left:     left,
			Operator: lowestTok.TokType,
			Right:    right,
		}
	}
}

func (p *Parser) parseLetStmt(toks []token.Token) *ast.LetStmtNode {
	if toks[0].TokType != token.LET {
		panic(fmt.Sprintf("[ERROR] Let statement is required to initialize variable, got %v\n", toks[0]))
	}
	if toks[1].TokType != token.VAR_NAME {
		panic(fmt.Sprintf("[ERROR] Could not figure out what to name variable, got %v\n", toks[1]))
	}
	if toks[2].TokType != token.ASSIGN {
		panic(fmt.Sprintf("[ERROR] Assignment operator is required to assign value to variable, got %v\n", toks[2]))
	}
	name := toks[1].Literal
	val := p.parseExpression(toks[3:])
	return &ast.LetStmtNode{
		Name:  name,
		Value: val,
	}
}

func (p *Parser) parseVarReference(tok token.Token) *ast.ReferenceExprNode {
	if tok.TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Expected name of variable, got %v\n", tok))
	}
	return &ast.ReferenceExprNode{
		Name: tok.Literal,
	}
}

func (p *Parser) parseVarReassign(toks []token.Token) *ast.VarReassignNode {
	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Expected var name, got %v\n", toks[0]))
	}
	if toks[1].TokType != token.ASSIGN {
		panic(fmt.Sprintf("[ERROR] Expected equals sign, got %v\n", toks[0]))
	}
	name := p.parseVarReference(toks[0])
	value := p.parseExpression(toks[2:])
	return &ast.VarReassignNode{
		Var:    *name,
		NewVal: value,
	}
}

func (p *Parser) parseIfStmt(toks []token.Token) *ast.IfStmtNode {
	if toks[0].TokType != token.IF {
		panic(fmt.Sprintf("[ERROR] Expected \"IF\" got %v\n", toks[0]))
	}

	var condToks []token.Token
	condLen := 0
	for i, val := range toks[1:] {
		if val.TokType != token.LBRACE {
			condToks = append(condToks, val)
		} else {
			condLen = i
			break
		}
	}
	cond := p.parseExpression(condToks)
	body := p.splitIntoLines(toks[condLen+2 : len(toks)-1])
	var parsedStmts []ast.Node
	for _, val := range body {
		n := p.parseStmt(val)
		if n != nil {
			parsedStmts = append(parsedStmts, n)
		}
	}
	var toReturn *ast.IfStmtNode
	if cond.NodeType() == ast.BoolLiteral {
		boolCond, ok := cond.(*ast.BoolLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %+v\n", cond))
		}
		toReturn = &ast.IfStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.BoolInfix {
		boolCond, ok := cond.(*ast.BoolInfixNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		toReturn = &ast.IfStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.ReferenceExpr {
		refExpr, ok := cond.(*ast.ReferenceExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		toReturn = &ast.IfStmtNode{
			Cond: &ast.BoolInfixNode{
				Left:     refExpr,
				Operator: token.OR,
				Right:    &ast.BoolLiteralNode{Value: false},
			},
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.PrefixExpr {
		prefixExpr, ok := cond.(*ast.PrefixExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		toReturn = &ast.IfStmtNode{
			Cond: prefixExpr,
			Body: parsedStmts,
		}
	}
	if toReturn == nil {
		panic(fmt.Sprintf("[ERROR] Could not parse if statement, tokens are %+v, cond types are %v\n", toks, cond.NodeType()))
	}
	p.ifStack = append(p.ifStack, toReturn)
	return toReturn
}

func (p *Parser) parseWhileStmt(toks []token.Token) *ast.WhileStmtNode {
	if toks[0].TokType != token.WHILE {
		panic(fmt.Sprintf("[ERROR] Expected \"WHILE\" got %v\n", toks[0]))
	}
	var condToks []token.Token
	condLen := 0
	for i, val := range toks[1:] {
		if val.TokType != token.LBRACE {
			condToks = append(condToks, val)
		} else {
			condLen = i
			break
		}
	}
	cond := p.parseExpression(condToks)
	body := p.splitIntoLines(toks[condLen+2 : len(toks)-1])
	var parsedStmts []ast.Node
	for _, val := range body {
		n := p.parseStmt(val)
		if n != nil {
			parsedStmts = append(parsedStmts, n)
		}
	}
	if cond.NodeType() == ast.BoolLiteral {
		boolCond, ok := cond.(*ast.BoolLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %+v\n", cond))
		}
		return &ast.WhileStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.BoolInfix {
		boolCond, ok := cond.(*ast.BoolInfixNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		return &ast.WhileStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.ReferenceExpr {
		refExpr, ok := cond.(*ast.ReferenceExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		return &ast.WhileStmtNode{
			Cond: &ast.BoolInfixNode{
				Left:     refExpr,
				Operator: token.OR,
				Right:    &ast.BoolLiteralNode{Value: false},
			},
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.PrefixExpr {
		prefixExpr, ok := cond.(*ast.PrefixExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		return &ast.WhileStmtNode{
			Cond: prefixExpr,
			Body: parsedStmts,
		}
	}
	panic(fmt.Sprintf("[ERROR] Could not parse while statement, tokens are %+v\n, cond types are %v\n", toks, cond.NodeType()))
}

func (p *Parser) parseElseStmt(toks []token.Token) {
	if toks[0].TokType != token.ELSE {
		panic(fmt.Sprintf("[ERROR] Expected else, got %v\n", toks[0]))
	}

	bodyTokens := toks[2:]
	body := p.splitIntoLines(bodyTokens)

	var stmts []ast.Node
	for _, val := range body {
		if len(val) > 0 {
			stmt := p.parseStmt(val)
			if stmt != nil {
				stmts = append(stmts, stmt)
			}
		}
	}

	if len(p.ifStack) == 0 {
		panic("[ERROR] Could not find if to attach else to")
	}

	lastIfIndex := len(p.ifStack) - 1
	lastIf := p.ifStack[lastIfIndex]
	lastIf.Alt = stmts
	p.ifStack = p.ifStack[:lastIfIndex]
}

func (p *Parser) parseStmt(line []token.Token) ast.Node {
	if len(line) == 0 {
		return nil
	}

	firstTok := line[0]

	var secondTok token.Token
	if len(line) > 1 {
		secondTok = line[1]
	}

	if len(line) == 1 {
		if firstTok.TokType == token.RBRACE || firstTok.TokType == token.LBRACE || firstTok.TokType == token.SEMICOLON {
			return nil
		}
		if firstTok.TokType == token.CONTINUE {
			return &ast.ContinueStmtNode{}
		}
		if firstTok.TokType == token.BREAK {
			return &ast.BreakStmtNode{}
		}
		return p.parseExpression(line)
	}

	if firstTok.TokType == token.LET {
		return p.parseLetStmt(line)
	}
	if firstTok.TokType == token.VAR_REF && secondTok.TokType == token.ASSIGN {
		return p.parseVarReassign(line)
	}
	if firstTok.TokType == token.IF {
		return p.parseIfStmt(line)
	}
	if firstTok.TokType == token.ELSE {
		p.parseElseStmt(line)
		return nil
	}
	if firstTok.TokType == token.FN {
		return p.parseFuncDecStmt(line)
	}
	if firstTok.TokType == token.RETURN {
		return p.parseReturnExpr(line)
	}
	if firstTok.TokType == token.WHILE {
		return p.parseWhileStmt(line)
	}
	if firstTok.TokType == token.CONTINUE {
		return &ast.ContinueStmtNode{}
	}
	if firstTok.TokType == token.BREAK {
		return &ast.BreakStmtNode{}
	}
	return p.parseExpression(line)
}

func (p *Parser) Parse(tokens []token.Token) ast.ProgramNode {
	var tokGroups [][]token.Token = p.splitIntoLines(p.preProcess(tokens))
	for _, line := range tokGroups {
		n := p.parseStmt(line)
		if n != nil {
			p.program.Statements = append(p.program.Statements, n)
		}
	}
	return p.program
}