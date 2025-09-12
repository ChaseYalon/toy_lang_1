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
		token.PLUS:     1,
		token.MINUS:    1,
		token.MULTIPLY: 2,
		token.DIVIDE:   2,
		token.BOOLEAN:  100,
		token.INTEGER:  100, //Boolean, int, and var ref should never be "bound to"
		token.VAR_REF:  100,
		token.AND:      3,
		token.OR:       3,
		token.NOT:      4, //Logical operators have lower precedence than arithmetic, not is lowest

	}
}

func (p *Parser) preProcess(tokens []token.Token) []token.Token {
	//Handles preprocessing for compound operators
	var toReturn []token.Token
	for i, val := range tokens {
		if val.TokType == token.COMPOUND_PLUS {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.PLUS, "+"))
			continue
		}
		if val.TokType == token.PLUS_PLUS {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.PLUS, "+"))
			toReturn = append(toReturn, *token.NewToken(token.INTEGER, "1"))
			continue
		}
		if val.TokType == token.COMPOUND_MINUS {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.MINUS, "-"))
			continue
		}
		if val.TokType == token.MINUS_MINUS {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.MINUS, "-"))
			toReturn = append(toReturn, *token.NewToken(token.INTEGER, "1"))
			continue
		}
		if val.TokType == token.COMPOUND_MULTIPLY {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.MULTIPLY, "*"))
			continue
		}
		if val.TokType == token.COMPOUND_DIVIDE {
			toReturn = append(toReturn, *token.NewToken(token.ASSIGN, "="))
			toReturn = append(toReturn, tokens[i-1])
			toReturn = append(toReturn, *token.NewToken(token.DIVIDE, "/"))
			continue
		}
		toReturn = append(toReturn, val)
	}
	return toReturn
}
func (p *Parser) splitBySemicolon(tokens []token.Token) [][]token.Token {
	var toReturn [][]token.Token
	var current []token.Token
	for _, val := range tokens {
		if val.TokType == token.SEMICOLON {
			toReturn = append(toReturn, current)
			current = []token.Token{}
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
	//Functional infinity to start
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
	case token.AND, token.OR:
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
	//Do Sematic checks
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
	//Sematic checks
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
func (p *Parser) Parse(tokens []token.Token) ast.ProgramNode {
	var tokGroups [][]token.Token = p.splitBySemicolon(p.preProcess(tokens))
	/*
		At the start of each tok group there should be one of the following
			Let -> Making New variable
			Var_Ref -> Referencing old variable
			Var_ref, assign -> Reassigning old variable
	*/
	for _, line := range tokGroups {
		firstTok := line[0]
		secondTok := line[1]
		if firstTok.TokType == token.LET {
			p.program.Statements = append(p.program.Statements, p.parseLetStmt(line))
			continue
		}
		if firstTok.TokType == token.VAR_REF && secondTok.TokType == token.ASSIGN {
			p.program.Statements = append(p.program.Statements, p.parseVarReassign(line))
			continue
		}
		//If it is not a let statement or a reassign statement assume it an expression
		p.program.Statements = append(p.program.Statements, p.parseExpression(line))
	}
	return p.program
}
