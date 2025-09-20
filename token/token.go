package token

type TokenType int

const (
	//Misc
	ILLEGAL TokenType = iota

	//Arithmetic
	PLUS
	MINUS
	MULTIPLY
	DIVIDE

	//Syntax
	SEMICOLON
	ASSIGN
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	COMMA

	//Keywords
	VAR_REF
	VAR_NAME
	LET
	IF
	ELSE
	FN
	RETURN
	FUNC_NAME

	//Compound operators
	COMPOUND_PLUS
	COMPOUND_MINUS
	COMPOUND_MULTIPLY
	COMPOUND_DIVIDE
	PLUS_PLUS
	MINUS_MINUS

	//Datatypes
	INTEGER
	BOOLEAN
	STRING

	//Boolean operators
	LESS_THAN
	LESS_THAN_EQT
	GREATER_THAN
	GREATER_THAN_EQT
	EQUALS
	NOT_EQUAL
	AND
	OR
	NOT

	//Placeholder
	EMPTY
)

func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case INTEGER:
		return "INTEGER"
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULTIPLY:
		return "*"
	case DIVIDE:
		return "/"
	case LET:
		return "LET"
	case ASSIGN:
		return "ASSIGN"
	case VAR_NAME:
		return "VAR_NAME"
	case SEMICOLON:
		return "SEMICOLON"
	case VAR_REF:
		return "VAR_REF"
	case COMPOUND_PLUS:
		return "COMPOUND_PLUS"
	case COMPOUND_MINUS:
		return "COMPOUND_MINUS"
	case COMPOUND_MULTIPLY:
		return "COMPOUND_MULTIPLY"
	case COMPOUND_DIVIDE:
		return "COMPOUND_DIVIDE"
	case PLUS_PLUS:
		return "PLUS_PLUS"
	case MINUS_MINUS:
		return "MINUS_MINUS"
	case BOOLEAN:
		return "BOOLEAN"
	case LESS_THAN:
		return "<"
	case LESS_THAN_EQT:
		return "<="
	case GREATER_THAN:
		return ">"
	case GREATER_THAN_EQT:
		return ">="
	case EQUALS:
		return "=="
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case IF:
		return "IF"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case ELSE:
		return "ELSE"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case EMPTY:
		return "EMPTY"
	case FN:
		return "FN"
	case FUNC_NAME:
		return "FUNC_NAME"
	case RETURN:
		return "RETURN"
	case COMMA:
		return "COMMA"
	case STRING:
		return "STRING"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	TokType TokenType
	Literal string
}

func NewToken(tokType TokenType, literal string) *Token {
	return &Token{TokType: tokType, Literal: literal}
}
