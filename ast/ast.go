package ast

import (
	"fmt"
	"toy_lang/token"
)

type AstNode int

type Node interface {
	NodeType() AstNode
	String() string
}
type Bool interface {
	isBool()
	NodeType() AstNode
	String() string
}

const (
	//Main
	Program AstNode = iota

	//Vars
	LetStmt
	ReferenceExpr
	VarReassign

	//Exprs
	InfixExpr
	BoolInfix
	PrefixExpr
	EmptyExpr
	ReturnExpr

	//Datatypes
	IntLiteral
	BoolLiteral
	StringLiteral

	//Statements
	IfStmt
	FuncDec
	FuncCall
	CallBuiltin
)

func (n AstNode) String() string {
	switch n {
	case LetStmt:
		return "LET_STMT"
	case InfixExpr:
		return "INFIX_EXPR"
	case IntLiteral:
		return "INTEGER_LITERAL"
	case ReferenceExpr:
		return "REF_EXPR"
	case VarReassign:
		return "VAR_REASSIGN"
	case BoolLiteral:
		return "BOOL_LITERAL"
	case BoolInfix:
		return "BOOL_INFIX"
	case IfStmt:
		return "IF_STMT"
	case PrefixExpr:
		return "PREFIX_EXPR"
	case EmptyExpr:
		return "EMPTY_EXPR" //Parens
	case FuncCall:
		return "FUNC_CALL"
	case FuncDec:
		return "FUNC_DEC"
	case ReturnExpr:
		return "RETURN_EXPR"
	case StringLiteral:
		return "STRING_LITERAL"
	case CallBuiltin:
		return "CALL_BUILTIN"
	default:
		return "ILLEGAL"
	}
}

// Let Expression
type LetStmtNode struct {
	Name  string
	Value Node
}

func (n *LetStmtNode) NodeType() AstNode {
	return LetStmt
}

func (n *LetStmtNode) String() string {
	return fmt.Sprintf("let %v = %v", n.Name, n.Value)
}

// Infix Expression
type InfixExprNode struct {
	Left     Node
	Operator token.TokenType
	Right    Node
}

func (n *InfixExprNode) NodeType() AstNode {
	return InfixExpr
}
func (n *InfixExprNode) String() string {
	return fmt.Sprintf("(%v %v %v)", n.Left, n.Operator, n.Right)

}

// Integer Literal
type IntLiteralNode struct {
	Value int
}

func (n *IntLiteralNode) NodeType() AstNode {
	return IntLiteral
}

func (n *IntLiteralNode) String() string {
	return fmt.Sprintf("INT(%d)", n.Value)
}

// Variable Reference
type ReferenceExprNode struct {
	Name string
}

func (n *ReferenceExprNode) NodeType() AstNode {
	return ReferenceExpr
}
func (n *ReferenceExprNode) String() string {
	return fmt.Sprintf("REFERENCE(%v)", n.Name)
}
func (n *ReferenceExprNode) isBool() {}

type VarReassignNode struct {
	Var    ReferenceExprNode
	NewVal Node
}

func (n *VarReassignNode) NodeType() AstNode {
	return VarReassign
}

func (n *VarReassignNode) String() string {
	return fmt.Sprintf("REASSIGN(%v) = %v", n.Var, n.NewVal)
}

// Program
type ProgramNode struct {
	Statements []Node
}

func (n *ProgramNode) NodeType() AstNode {
	return Program
}
func (n *ProgramNode) String() string {
	str := "Program{\n"
	for _, val := range n.Statements {
		str += (val.String() + ",\n")
	}
	str += "}"
	return str

}

// Bool literal
type BoolLiteralNode struct {
	Value bool
}

func (n *BoolLiteralNode) NodeType() AstNode {
	return BoolLiteral
}
func (n *BoolLiteralNode) String() string {
	if n.Value {
		return "BOOL(true)"
	}
	return "BOOL(false)"
}

// Annoying thing to make go work
func (n *BoolLiteralNode) isBool() {}

type BoolInfixNode struct {
	Left     Node
	Operator token.TokenType
	Right    Node
}

func (n *BoolInfixNode) NodeType() AstNode {
	return BoolInfix
}
func (n *BoolInfixNode) String() string {
	return fmt.Sprintf("(%v %v %v)", n.Left, n.Operator, n.Right)
}
func (n *BoolInfixNode) isBool() {}

type PrefixExprNode struct {
	Value    Node
	Operator token.TokenType
}

func (n *PrefixExprNode) NodeType() AstNode {
	return PrefixExpr
}
func (n *PrefixExprNode) String() string {
	return fmt.Sprintf("%v%v", n.Operator, n.Value)
}
func (n *PrefixExprNode) isBool() {}

type IfStmtNode struct {
	Cond Bool
	Body []Node
	Alt  []Node
}

func (n *IfStmtNode) NodeType() AstNode {
	return IfStmt
}
func (n *IfStmtNode) String() string {
	var str string = fmt.Sprintf("if %v {\n", n.Cond)
	for _, val := range n.Body {
		str += fmt.Sprintf("\t%v\n", val)
	}
	str += "}"
	if len(n.Alt) == 0 {
		return str
	}
	str += " else {\n"
	for _, val := range n.Alt {
		str += fmt.Sprintf("\t%v\n", val)
	}
	str += "}"
	return str
}

type EmptyExprNode struct {
	Child Node
}

func (n *EmptyExprNode) NodeType() AstNode {
	return EmptyExpr
}
func (n *EmptyExprNode) String() string {
	return fmt.Sprintf("(%v)", n.Child)
}

type ReturnExprNode struct {
	Val Node
}

func (n *ReturnExprNode) NodeType() AstNode {
	return ReturnExpr
}
func (n *ReturnExprNode) String() string {
	return fmt.Sprintf("return %v", n.Val)
}

type FuncDecNode struct {
	Name   string
	Params []ReferenceExprNode
	Body   []Node
	Return ReturnExprNode
}

func (n *FuncDecNode) NodeType() AstNode {
	return FuncDec
}
func (n *FuncDecNode) String() string {
	str := fmt.Sprintf("fn %v(", n.Name)
	for i, p := range n.Params {
		if i > 0 {
			str += ", "
		}
		str += p.String()
	}
	str += ") {\n"
	for _, stmt := range n.Body {
		str += fmt.Sprintf("\t%v\n", stmt)
	}
	if n.Return.Val != nil {
		str += fmt.Sprintf("\t%v\n", n.Return.String())
	}
	str += "}\n"
	return str
}

type FuncCallNode struct {
	Name   ReferenceExprNode
	Params []Node
}

func (n *FuncCallNode) NodeType() AstNode {
	return FuncCall
}
func (n *FuncCallNode) String() string {
	return fmt.Sprintf("%v(%+v)", n.Name.Name, n.Params)
}

type StringLiteralNode struct {
	Value string
}

func (n *StringLiteralNode) NodeType() AstNode {
	return StringLiteral
}
func (n *StringLiteralNode) String() string {
	return fmt.Sprintf("STRING(%v)", n.Value)
}

type CallBuiltinNode struct{
	Name string
	Params []Node
}
func (n *CallBuiltinNode)NodeType() AstNode{
	return CallBuiltin
}
func (n *CallBuiltinNode)String()string{
	return fmt.Sprintf("BUILTIN_FN_%v(%+v)\n", n.Name, n.Params);
}