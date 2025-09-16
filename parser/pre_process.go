package parser

import (
	"toy_lang/token"
)

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
	//Second pass replace let y = -5; with let y = 0 - 5;;
	for i, val := range toReturn {
		if val.TokType == token.MINUS &&
			i > 0 &&
			(toReturn[i-1].TokType != token.INTEGER &&
				toReturn[i-1].TokType != token.RPAREN &&
				toReturn[i-1].TokType != token.VAR_REF) {

			// insert 0 before the minus
			toReturn = append(toReturn[:i],
				append([]token.Token{*token.NewToken(token.INTEGER, "0"), val}, toReturn[i+1:]...)...)
		}

	}
	return toReturn
}
