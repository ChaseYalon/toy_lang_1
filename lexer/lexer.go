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
}

func NewLexer() *Lexer {
	return &Lexer{
		chars:      []rune{},
		currNum:    []rune{},
		currString: []rune{},
		pos:        0,
		tokens:     []token.Token{},
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
		if l.tokens[len(l.tokens)-1].TokType == token.LET {
			l.tokens = append(l.tokens, *token.NewToken(token.VAR_NAME, string(l.currString)))
			l.currString = []rune{}
			return
		}
		l.tokens = append(l.tokens, *token.NewToken(token.VAR_REF, string(l.currString)))
		l.currString = []rune{}
		return
	}
}

func (l *Lexer) eat() {
	l.pos++
}

func (l *Lexer) Lex(s string) []token.Token {
	l.chars = []rune(s)
	l.pos = 0
	l.tokens = []token.Token{}
	l.currNum = []rune{}
	l.currString = []rune{}

	for l.pos < len(l.chars) {
		ch := l.getChar()

		if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' {
			l.flushInt()
			l.flushStr()
			l.eat()
			continue;
		}

		if (ch == 'l' || ch == 'L') &&
			(l.peek(1) == 'e' || l.peek(1) == 'E') &&
			(l.peek(2) == 't' || l.peek(2) == 'T') {
			l.flushStr()
			l.flushInt()

			l.tokens = append(l.tokens, *token.NewToken(token.LET,
				string([]rune{ch, l.peek(1), l.peek(2)})))
			l.pos += 3
			continue
		}
		if ch == 't' && l.peek(1) == 'r' && l.peek(2) == 'u' && l.peek(3) == 'e'{
			l.flushStr();
			l.flushInt();
			l.tokens = append(l.tokens, *token.NewToken(token.BOOLEAN, "true"));
			l.pos += 4;
			continue;
		}
		if ch == 'f' && l.peek(1) == 'a' && l.peek(2) == 'l' && l.peek(3) == 's' && l.peek(4) == 'e'{
			l.flushStr();
			l.flushInt();
			l.tokens = append(l.tokens, *token.NewToken(token.BOOLEAN, "false"));
			l.pos += 5;
			continue;
		}
		if ch == ';' {
			l.flushInt()
			l.flushStr()
			l.tokens = append(l.tokens, *token.NewToken(token.SEMICOLON, ";"))
			l.eat()
			continue
		}
		switch {
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
					if l.peek(1) == '='{
						l.tokens = append(l.tokens, *token.NewToken(token.EQUALS, "=="));
						l.eat();
					} else {

						l.tokens = append(l.tokens, *token.NewToken(token.ASSIGN, "="))
					}
				case ch == '>':
					l.flushInt();
					l.flushStr();
					if l.peek(1) == '='{
						l.tokens = append(l.tokens, *token.NewToken(token.GREATER_THAN_EQT, ">="))
						l.eat();
					} else {
						l.tokens = append(l.tokens, *token.NewToken(token.GREATER_THAN, ">"));
					}
				case ch == '<':
					l.flushInt();
					l.flushStr();
					if l.peek(1) == '='{
						l.tokens = append(l.tokens, *token.NewToken(token.LESS_THAN_EQT, "<="));
						l.eat();
					} else {
						l.tokens = append(l.tokens, *token.NewToken(token.LESS_THAN, "<"));
					}
				case ch == '&' && l.peek(1) == '&':
					l.flushInt();
					l.flushStr();
					l.tokens = append(l.tokens, *token.NewToken(token.AND, "&&"));
					l.eat();
				case ch == '|' && l.peek(1) == '|':
					l.flushInt();
					l.flushStr();
					l.tokens = append(l.tokens, *token.NewToken(token.OR, "||"));
					l.eat();
				case ch == '!':
					l.flushInt();
					l.flushStr();
					l.tokens = append(l.tokens, *token.NewToken(token.NOT, "!"));
					
				case unicode.IsLetter(ch):
					l.flushInt()
					l.currString = append(l.currString, ch)
				case '0' <= ch && ch <= '9':
					l.flushStr()
					l.currNum = append(l.currNum, ch)
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
