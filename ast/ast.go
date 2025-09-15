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
	Program AstNode = iota
	LetStmt
	InfixExpr
	IntLiteral
	ReferenceExpr
	VarReassign
	BoolLiteral
	BoolInfix
	PrefixExpr
	IfStmt
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
		return "VAR_REFERENCE"
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
	return fmt.Sprintf("%v = %v", n.Name, n.Value)
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
		return "true"
	}
	return "false"
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
