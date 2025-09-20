package lexer

import (
	"toy_lang/token"
	"unicode"
)

type Lexer struct {
	currNum    []rune
	currString []rune
	chars      []rune
	pos        int
	tokens     []token.Token
	inString   bool
}

func NewLexer() *Lexer {
	return &Lexer{
		chars:      []rune{},
		currNum:    []rune{},
		currString: []rune{},
		pos:        0,
		tokens:     []token.Token{},
		inString:   false,
	}
}

func (l *Lexer) getChar() rune {
	if l.pos >= len(l.chars) {
		return 0
	}
	return l.chars[l.pos]
}

func (l *Lexer) peek(n int) rune {
	if l.pos+n >= len(l.chars) {
		return 0
	}
	return l.chars[l.pos+n]
}

func (l *Lexer) flushInt() {
	if len(l.currNum) != 0 {
		l.tokens = append(l.tokens, *token.NewToken(token.INTEGER, string(l.currNum)))
		l.currNum = []rune{}
	}
}

func (l *Lexer) flushStr() {
	if len(l.currString) != 0 {
		if len(l.tokens) > 0 {
			if l.tokens[len(l.tokens)-1].TokType == token.LET {
				l.tokens = append(l.tokens, *token.NewToken(token.VAR_NAME, string(l.currString)))
				l.currString = []rune{}
				return
			}
			if l.tokens[len(l.tokens)-1].TokType == token.FN {
				l.tokens = append(l.tokens, *token.NewToken(token.FUNC_NAME, string(l.currString)))
				l.currString = []rune{}
				return
			}
		}
		l.tokens = append(l.tokens, *token.NewToken(token.VAR_REF, string(l.currString)))
		l.currString = []rune{}
		return
	}
}


func (l *Lexer) eat() {
	l.pos++
}
func (l *Lexer) parseKeyword(word string, tok token.Token) bool {
	for i, val := range []rune(word) {
		if val != l.peek(i) {
			return false
		}
	}
	l.flushStr()
	l.flushInt()
	l.tokens = append(l.tokens, tok)
	l.pos += len([]rune(word))
	return true
}

func (l *Lexer) Lex(s string) []token.Token {
	l.chars = []rune(s)
	l.pos = 0
	l.tokens = []token.Token{}
	l.currNum = []rune{}
	l.currString = []rune{}

	for l.pos < len(l.chars) {
		ch := l.getChar()
		if l.inString && ch != '"' {
			l.currString = append(l.currString, ch)
			l.eat()
			continue
		}

		if (ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r') && !l.inString {
			l.flushInt()
			l.flushStr()
			l.eat()
			continue
		}
		if l.parseKeyword("let", *token.NewToken(token.LET, "let")) {
			continue
		}
		if l.parseKeyword("true", *token.NewToken(token.BOOLEAN, "true")) {
			continue
		}
		if l.parseKeyword("false", *token.NewToken(token.BOOLEAN, "false")) {
			continue
		}
		if l.parseKeyword("if", *token.NewToken(token.IF, "if")) {
			continue
		}
		if l.parseKeyword("else", *token.NewToken(token.ELSE, "else")) {
			continue
		}
		if l.parseKeyword("fn", *token.NewToken(token.FN, "fn")) {
			continue
		}
		if l.parseKeyword("return", *token.NewToken(token.RETURN, "return")) {
			continue
		}

		switch {
		case ch == ';':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.SEMICOLON, ";"))
			l.eat()
			continue
		case ch == ',':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.COMMA, ","))
			l.eat()
			continue
		case ch == '+':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.COMPOUND_PLUS, "+="))
				l.eat()
			} else if l.peek(1) == '+' {
				l.tokens = append(l.tokens, *token.NewToken(token.PLUS_PLUS, "++"))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.PLUS, "+"))
			}
		case ch == '-':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.COMPOUND_MINUS, "-="))
				l.eat()
			} else if l.peek(1) == '-' {
				l.tokens = append(l.tokens, *token.NewToken(token.MINUS_MINUS, "--"))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.MINUS, "-"))
			}
		case ch == '*':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.COMPOUND_MULTIPLY, "*="))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.MULTIPLY, "*"))
			}
		case ch == '/':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.COMPOUND_DIVIDE, "/="))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.DIVIDE, "/"))
			}
		case ch == '=':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.EQUALS, "=="))
				l.eat()
			} else {

				l.tokens = append(l.tokens, *token.NewToken(token.ASSIGN, "="))
			}
		case ch == '"':
			if l.inString {
				// closing quote
				l.tokens = append(l.tokens, *token.NewToken(token.STRING, string(l.currString)))
				l.currString = []rune{}
				l.inString = false
			} else {
				// opening quote
				l.inString = true
			}
			l.eat()
			continue
		case ch == '>':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.GREATER_THAN_EQT, ">="))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.GREATER_THAN, ">"))
			}
		case ch == '<':
			l.flushInt()
			l.flushStr()
			if l.peek(1) == '=' {
				l.tokens = append(l.tokens, *token.NewToken(token.LESS_THAN_EQT, "<="))
				l.eat()
			} else {
				l.tokens = append(l.tokens, *token.NewToken(token.LESS_THAN, "<"))
			}
		case ch == '&' && l.peek(1) == '&':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.AND, "&&"))
			l.eat()
		case ch == '|' && l.peek(1) == '|':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.OR, "||"))
			l.eat()
		case ch == '!':
			if l.peek(1) == '='{
				l.flushInt();
				l.flushStr();
				l.tokens = append(l.tokens, *token.NewToken(token.NOT_EQUAL, "!="))
				l.eat();
			} else {

				l.flushInt()
				l.flushStr()
				l.tokens = append(l.tokens, *token.NewToken(token.NOT, "!"))
			}
		case ch == '{':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.LBRACE, "{"))
		case ch == '}':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.RBRACE, "}"))
		case ch == '(':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.LPAREN, "("))
		case ch == ')':
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.RPAREN, ")"))
		case unicode.IsLetter(ch) || (len(l.currString) > 0 && unicode.IsDigit(ch)):
			l.flushInt()
			l.currString = append(l.currString, ch)
			l.eat()
			continue
		case unicode.IsDigit(ch):
			l.flushStr()
			l.currNum = append(l.currNum, ch)
			l.eat()
			continue

		default:
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.ILLEGAL, string(ch)))
		}

		l.eat()
	}

	l.flushInt()
	l.flushStr()

	return l.tokens
}
