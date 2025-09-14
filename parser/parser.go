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
}


func NewParser() *Parser {
	return &Parser{
		program: ast.ProgramNode{
			Statements: []ast.Node{},
		},
		tokens: []token.Token{},
	}
}

func (p *Parser) generatePrecedenceTable() map[token.TokenType]int {
	return map[token.TokenType]int{
		token.PLUS:             1,
		token.MINUS:            1,
		token.MULTIPLY:         2,
		token.DIVIDE:           2,
		token.BOOLEAN:          100,
		token.INTEGER:          100, // Boolean, int, and var ref should never be "bound to"
		token.VAR_REF:          100,
		token.AND:              3,
		token.OR:               3,
		token.NOT:              4, // Logical operators have lower precedence than arithmetic, not is lowest
		token.LESS_THAN:        3,
		token.LESS_THAN_EQT:    3,
		token.GREATER_THAN:     3,
		token.GREATER_THAN_EQT: 3,
	}
}

func (p *Parser) preProcess(tokens []token.Token) []token.Token {
	// Handles preprocessing for compound operators
	var toReturn []token.Token
	for i, val := range tokens {
		if val.TokType == token.COMPOUND_PLUS {
			// ensure there's a left-hand token available
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.PLUS, "+"))
				continue
			}
		}
		if val.TokType == token.PLUS_PLUS {
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.PLUS, "+"))
				toReturn = append(toReturn, *token.NewToken(token.INTEGER, "1"))
				continue
			}
		}
		if val.TokType == token.COMPOUND_MINUS {
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.MINUS, "-"))
				continue
			}
		}
		if val.TokType == token.MINUS_MINUS {
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.MINUS, "-"))
				toReturn = append(toReturn, *token.NewToken(token.INTEGER, "1"))
				continue
			}
		}
		if val.TokType == token.COMPOUND_MULTIPLY {
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.MULTIPLY, "*"))
				continue
			}
		}
		if val.TokType == token.COMPOUND_DIVIDE {
			if i-1 >= 0 {
				toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
				toReturn = append(toReturn, tokens[i-1])
				toReturn = append(toReturn, *token.NewToken(token.DIVIDE, "/"))
				continue
			}
		}
		toReturn = append(toReturn, val)
	}
	return toReturn
}

func (p *Parser) splitIntoLines(tokens []token.Token) [][]token.Token {
	var toReturn [][]token.Token
	var inIfMode bool = false
	var current []token.Token
	for _, val := range tokens {
		if val.TokType == token.SEMICOLON && !inIfMode {
			toReturn = append(toReturn, current)
			current = []token.Token{}
			continue
		}
		if val.TokType == token.IF {
			inIfMode = true
			current = append(current, val)
			continue
		}
		// TEMPORARY CODE - THIS IS BAD
		if inIfMode && val.TokType == token.RBRACE {
			current = append(current, val)
			toReturn = append(toReturn, current)
			current = []token.Token{}
			inIfMode = false
			continue
		}
		current = append(current, val)
	}
	if len(current) > 0 {
		toReturn = append(toReturn, current)
	}
	return toReturn
}

func (p *Parser) findLowestBp(pt map[token.TokenType]int, tokens []token.Token) (token.Token, int) {
	// Functional infinity to start
	var lowestVal int = 10000000
	var lowestToken token.Token
	var lowestIndex int = -1
	for i, tok := range tokens {
		localVal := pt[tok.TokType]
		if localVal < lowestVal {
			lowestVal = localVal
			lowestToken = tok
			lowestIndex = i
		}
	}
	return lowestToken, lowestIndex
}

func (p *Parser) parseExpression(tokens []token.Token) ast.Node {
	lowestTok, lowestIndex := p.findLowestBp(p.generatePrecedenceTable(), tokens)
	if lowestTok.TokType == token.INTEGER {
		val, err := strconv.Atoi(lowestTok.Literal)
		if err != nil {
			panic(fmt.Sprintf("[ERROR] Could not convert from string to int literal, trying to convert %v", lowestTok))
		}
		return &ast.IntLiteralNode{Value: val}
	}
	if lowestTok.TokType == token.BOOLEAN {
		val, err := strconv.ParseBool(lowestTok.Literal)
		if err != nil {
			panic(fmt.Sprintf("[ERROR] Could not convert from string to bool literal, trying to convert %v", lowestTok))
		}
		return &ast.BoolLiteralNode{Value: val}
	}
	if lowestTok.TokType == token.VAR_REF {
		return &ast.ReferenceExprNode{Name: lowestTok.Literal}
	}
	if lowestIndex == -1 {
		panic(fmt.Sprintf("[ERROR] Could not find lowest binding power operator, input was %+v\n", tokens))
	}
	switch lowestTok.TokType {
	case token.AND, token.OR, token.LESS_THAN, token.LESS_THAN_EQT, token.GREATER_THAN, token.GREATER_THAN_EQT:
		left := p.parseExpression(tokens[:lowestIndex])
		right := p.parseExpression(tokens[lowestIndex+1:])
		return &ast.BoolInfixNode{
			Left:     left,
			Operator: lowestTok.TokType,
			Right:    right,
		}
	case token.NOT:
		return &ast.PrefixExprNode{
			Value:    p.parseExpression(tokens[lowestIndex+1:]),
			Operator: token.NOT,
		}
	default:
		left := p.parseExpression(tokens[:lowestIndex])
		right := p.parseExpression(tokens[lowestIndex+1:])
		return &ast.InfixExprNode{
			Left:     left,
			Operator: lowestTok.TokType,
			Right:    right,
		}
	}
}

func (p *Parser) parseLetStmt(toks []token.Token) *ast.LetStmtNode {
	// Do Sematic checks
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
	// Sematic checks
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
	// Do sematic analysis
	if toks[0].TokType != token.IF {
		panic(fmt.Sprintf("[ERROR] Expected \"IF\" got %v\n", toks[0]))
	}

	// Calculate condition
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
		return &ast.IfStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.BoolInfix {
		boolCond, ok := cond.(*ast.BoolInfixNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond))
		}
		return &ast.IfStmtNode{
			Cond: boolCond,
			Body: parsedStmts,
		}
	}
	if cond.NodeType() == ast.ReferenceExpr{
		refExpr, ok := cond.(*ast.ReferenceExprNode);
		if !ok{
			panic(fmt.Sprintf("[ERROR] Could not figure out conditional, got %v\n", cond));
		}
		return &ast.IfStmtNode{
			Cond: &ast.BoolInfixNode{
				Left: refExpr,
				Operator: token.OR,
				Right: &ast.BoolLiteralNode{Value: false},
			},
			Body: parsedStmts,
		}
	}

	panic(fmt.Sprintf("[ERROR] Could not parse if statement, tokens are %+v, cond types are %v\n", toks, cond.NodeType()))
}
func (p *Parser) parseElseStmt(toks []token.Token) {
    // Semantic checks
    if toks[0].TokType != token.ELSE {
        panic(fmt.Sprintf("[ERROR] Expected else, got %v\n", toks[0]))
    }
    
    // Extract the body between braces (skip ELSE and LBRACE, exclude RBRACE)
    bodyTokens := toks[2:]
    
    body := p.splitIntoLines(bodyTokens)
    
    var stmts []ast.Node
    for i, val := range body {
        if len(val) > 0 { // Only parse non-empty token slices
            stmt := p.parseStmt(val)
            if stmt != nil {
                stmts = append(stmts, stmt)
            }
        } else {
            fmt.Printf("Skipping empty slice at index %d\n", i) // Debug line
        }
    }
    
    // Find the last if statement to attach this else to
    if len(p.program.Statements) == 0 || p.program.Statements[len(p.program.Statements)-1].NodeType() != ast.IfStmt {
        panic("[ERROR] Could not find if to attach else to")
    }
    
    lastIf, ok := p.program.Statements[len(p.program.Statements)-1].(*ast.IfStmtNode)
    if !ok {
        panic("[ERROR] Could not find if to attach else to")
    }
    
    lastIf.Alt = stmts
    p.program.Statements[len(p.program.Statements)-1] = lastIf
}

func (p *Parser) parseStmt(line []token.Token) ast.Node {
	if len(line) == 0 {
		return nil
	}

	firstTok := line[0]

	// Avoid indexing line[1] if the line length is 1
	var secondTok token.Token
	if len(line) > 1 {
		secondTok = line[1]
	}

	// Skip lone braces or semicolons that ended up as their own "line"
	if len(line) == 1 {
		if firstTok.TokType == token.RBRACE || firstTok.TokType == token.LBRACE || firstTok.TokType == token.SEMICOLON {
			return nil
		}
		// For single-token lines, let parseExpression handle literals and var refs.
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
	if firstTok.TokType == token.ELSE{
		p.parseElseStmt(line);
		return nil;
	}
	// If it is not a let statement or a reassign statement assume it is an expression
	return p.parseExpression(line)
}

func (p *Parser) Parse(tokens []token.Token) ast.ProgramNode {
	var tokGroups [][]token.Token = p.splitIntoLines(p.preProcess(tokens))
	/*
		At the start of each tok group there should be one of the following
			Let -> Making New variable
			Var_Ref -> Referencing old variable
			Var_ref, assign -> Reassigning old variable
	*/
	for _, line := range tokGroups {
		n := p.parseStmt(line)
		if n != nil {
			p.program.Statements = append(p.program.Statements, n)
		}
	}
	return p.program
}
