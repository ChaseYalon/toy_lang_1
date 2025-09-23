package parser

import (
	"toy_lang/token"
)

func (p *Parser) generatePrecedenceTable() map[token.TokenType]int {
	return map[token.TokenType]int{
		token.PLUS:             1,
		token.MINUS:            1,
		token.MULTIPLY:         2,
		token.DIVIDE:           2,
		token.MODULO:           2,
		token.EXPONENT:         3,
		token.BOOLEAN:          100,
		token.INTEGER:          100, // Boolean, int, string, and var ref should never be "bound to"
		token.STRING:           100,
		token.VAR_REF:          100,
		token.FLOAT:            100,
		token.AND:              4,
		token.OR:               4,
		token.NOT:              50, // Logical operators have lower precedence than arithmetic, not is lowest
		token.LESS_THAN:        4,
		token.LESS_THAN_EQT:    4,
		token.GREATER_THAN:     4,
		token.GREATER_THAN_EQT: 4,
		token.EQUALS:           4,
	}
}

func (p *Parser) splitIntoLines(tokens []token.Token) [][]token.Token {
	var lines [][]token.Token
	var current []token.Token
	inBlock := 0
	parenDepth := 0

	for _, tok := range tokens {
		current = append(current, tok)

		switch tok.TokType {
		case token.LBRACE:
			inBlock++
		case token.RBRACE:
			inBlock--
			if inBlock == 0 {
				lines = append(lines, current)
				current = []token.Token{}
			}
		case token.LPAREN:
			parenDepth++
		case token.RPAREN:
			parenDepth--
		case token.SEMICOLON:
			if inBlock == 0 && parenDepth == 0 {
				lines = append(lines, current[:len(current)-1])
				current = []token.Token{}
			}
		}
	}

	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

func (p *Parser) findLowestBp(pt map[token.TokenType]int, tokens []token.Token) (token.Token, int) {
	lowestVal := 10000000
	var lowestTok token.Token
	lowestIndex := -1
	depth := 0

	for i, tok := range tokens {
		switch tok.TokType {
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
		default:
			if depth == 0 {
				if tok.TokType != token.NOT{

					if val, ok := pt[tok.TokType]; ok && (val <= lowestVal ) {
						lowestVal = val
						lowestTok = tok
						lowestIndex = i
					}
				}else {
					if val, ok := pt[tok.TokType]; ok && (val < lowestVal ) {
						lowestVal = val
						lowestTok = tok
						lowestIndex = i
					}
				}
			}
		}
	}
	return lowestTok, lowestIndex
}

//Random helper function because go takes minimalism to far

func includesItem(arr []token.Token, tok token.Token) (bool, int) {
	for i, val := range arr {
		if val == tok {
			return true, i
		}
	}
	return false, -1
}

func includesAny(arr []token.Token, checkFor []token.TokenType) (bool, int) {
	for i, val := range arr {
		for _, val2 := range checkFor {
			if val2 == val.TokType {
				return true, i
			}
		}
	}
	return false, -1
}
