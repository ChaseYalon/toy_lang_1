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
		token.BOOLEAN: 0, 
		token.INTEGER: 0, //Boolean and int should never be "bound to"
	}
}

func (p *Parser) parseVarAssign(toks []token.Token) ast.Node {
	if toks[0].TokType != token.LET {
		panic(fmt.Sprintf("[ERROR] Expected LET received %s", toks[0]))
	}
	if toks[1].TokType != token.VAR_NAME {
		panic(fmt.Sprintf("[ERROR] Expected VAR_NAME received %s", toks[1]))
	}
	if toks[2].TokType != token.ASSIGN {
		panic(fmt.Sprintf("[ERROR] Expected '=' received %s", toks[2]))
	}
	value := p.parseExpression(toks[3:])
	return &ast.LetExprNode{
		Name:  toks[1].Literal,
		Value: value,
	}
}

func (p *Parser) parseVarRef(toks []token.Token) *ast.ReferenceExprNode {

	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Expected VAR_REF received %s", toks[0].TokType))
	}
	return &ast.ReferenceExprNode{
		Name: toks[0].Literal,
	}
}

func (p *Parser) parseVarReassign(toks []token.Token) *ast.VarReassignNode {
	//Format should be <VAR_REF> = <VALUE>
	if toks[0].TokType != token.VAR_REF {
		panic(fmt.Sprintf("[ERROR] Could not find what variable to reassign, got: %v", toks[0]))
	}
	if toks[1].TokType != token.ASSIGN {
		panic("[ERROR] Reassigning variables requires an equal sign between the name and the value")
	}
	referencedVar := ast.ReferenceExprNode{
		Name: toks[0].Literal,
	}
	//Case 1: Reassigning to Int Literal
	if toks[2].TokType == token.INTEGER {
		newValInt, err := strconv.Atoi(toks[2].Literal)
		if err != nil {
			panic(fmt.Sprintf("[ERROR] Could not convert variable reassign value to int literal, got %v", toks[2].Literal))
		}
		return &ast.VarReassignNode{
			Var: referencedVar,
			NewVal: &ast.IntLiteralNode{
				Value: newValInt,
			},
		}
	}
	//Case 2: Reassigning to boolean expression
	if toks[2].TokType == token.BOOLEAN{
		var boolVal bool;
		if toks[2].Literal == "true"{
			boolVal = true;
		} else {
			boolVal = false;
		}
		return &ast.VarReassignNode{
			Var: referencedVar,
			NewVal: &ast.BoolLiteralNode{
				Value: boolVal,
			},
		}
	}
	//Case 3: Reassigning to an Infix Expression, Default case
	exprToks := toks[2:]
	newValExpr := p.parseExpression(exprToks)
	return &ast.VarReassignNode{
		Var:    referencedVar,
		NewVal: newValExpr,
	}
}

func (p *Parser) findLowestBp(toks []token.Token) (token.Token, int) {
	p_table := p.generatePrecedenceTable()
	lowestBp := 100000
	lowestBpIdx := -1
	var lowestOp token.Token
	for i, val := range toks {
		if val.TokType == token.PLUS || val.TokType == token.MINUS ||
			val.TokType == token.MULTIPLY || val.TokType == token.DIVIDE {
			if p_table[val.TokType] <= lowestBp {
				lowestBp = p_table[val.TokType]
				lowestOp = val
				lowestBpIdx = i
			}
		}
	}
	if lowestBpIdx == -1 {
		panic("[ERROR] No operator found in token slice")
	}
	return lowestOp, lowestBpIdx
}

func (p *Parser) parseExpression(toks []token.Token) ast.Node {
	if len(toks) == 0 {
		panic("[ERROR] parseExpression received empty token slice")
	}
	if len(toks) == 1 {
		switch toks[0].TokType {
		case token.INTEGER:
			val, err := strconv.Atoi(toks[0].Literal)
			if err != nil {
				panic(fmt.Sprintf("[ERROR] Failed to convert %s to int", toks[0].Literal))
			}
			return &ast.IntLiteralNode{Value: val}
		case token.BOOLEAN:
			var boolVal bool;
			if toks[0].Literal == "true"{
				boolVal = true;
			} else {
				boolVal = false;
			}
			return &ast.BoolLiteralNode{Value: boolVal};
		case token.VAR_REF:
			return p.parseVarRef(toks)

		}
		
	}

	if len(toks) == 3 {
		op := toks[1]
		var left, right ast.Node

		if toks[0].TokType == token.INTEGER {
			val, _ := strconv.Atoi(toks[0].Literal)
			left = &ast.IntLiteralNode{Value: val}
		} else {
			left = p.parseVarRef([]token.Token{toks[0]})
		}

		if toks[2].TokType == token.INTEGER {
			val, _ := strconv.Atoi(toks[2].Literal)
			right = &ast.IntLiteralNode{Value: val}
		} else {
			right = p.parseVarRef([]token.Token{toks[2]})
		}

		return &ast.InfixExprNode{
			Left:     left,
			Operator: op.TokType,
			Right:    right,
		}
	}

	// recursive for longer expressions
	op, idx := p.findLowestBp(toks)
	left := p.parseExpression(toks[:idx])
	right := p.parseExpression(toks[idx+1:])
	return &ast.InfixExprNode{
		Left:     left,
		Operator: op.TokType,
		Right:    right,
	}
}

func splitBySemiColon(toks []token.Token) [][]token.Token {
	var result [][]token.Token
	var current []token.Token
	for _, t := range toks {
		if t.TokType == token.SEMICOLON {
			if len(current) > 0 {
				result = append(result, current)
			}
			current = []token.Token{}
		} else {
			current = append(current, t)
		}
	}
	if len(current) > 0 {
		result = append(result, current)
	}
	return result
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
func (p *Parser) Parse(tokens []token.Token) ast.ProgramNode {
	u_in := p.preProcess(tokens)
	lines := splitBySemiColon(u_in)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		switch line[0].TokType {
		case token.LET:
			p.program.Statements = append(p.program.Statements, p.parseVarAssign(line))
		case token.VAR_REF:
			if len(line) > 1 && line[1].TokType == token.ASSIGN {
				p.program.Statements = append(p.program.Statements, p.parseVarReassign(line))
			} else {
				p.program.Statements = append(p.program.Statements, p.parseExpression(line))
			}
		case token.INTEGER:
			p.program.Statements = append(p.program.Statements, p.parseExpression(line))
		default:
			panic(fmt.Sprintf("[ERROR] Unknown token: %v", line[0]))
		}
	}
	return p.program
}
