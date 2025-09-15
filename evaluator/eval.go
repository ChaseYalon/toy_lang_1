package evaluator

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

type v_map map[string]ast.Node

func (v v_map) String() string {
	s := "{"
	for k, node := range v {
		s += k + ": " + fmt.Sprintf("%v", node) + ", "
	}
	s += "}"
	return s
}

type Scope struct {
	Vars   map[string]ast.Node
	Parent *Scope
}

// Look up a variable recursively
func (s *Scope) getVar(name string) (ast.Node, bool) {
	if val, ok := s.Vars[name]; ok {
		return val, true
	}
	if s.Parent != nil {
		return s.Parent.getVar(name)
	}
	return nil, false
}

func (s *Scope) declareVar(name string, val ast.Node) {
	s.Vars[name] = val
}

func (s *Scope) assignVar(name string, val ast.Node) bool {
	if val.NodeType() == ast.IntLiteral || val.NodeType() == ast.BoolLiteral {
		if _, ok := s.Vars[name]; ok {
			s.Vars[name] = val
			return true
		}
		if s.Parent != nil {
			return s.Parent.assignVar(name, val)
		}
		return false
	}
	panic(fmt.Sprintf("[ERROR] Tried to assign non primitive value to variable, got %v\n", val))
}

func (s *Scope) newChild() *Scope {
	return &Scope{
		Vars:   make(v_map),
		Parent: s,
	}
}

func (s *Scope) String() string {
	return fmt.Sprintf("Vars: %+v, Parent: %v\n", s.Vars, s.Parent)
}

type Interpreter struct {
	MainScope Scope
}

func NewInterpreter() Interpreter {
	return Interpreter{
		MainScope: Scope{
			Vars: make(v_map),
		},
	}
}

func (i *Interpreter) execIntExpr(inode ast.Node, local_scope *Scope) int {
	var node ast.Node = inode
	
	// Handle EmptyExprNode by unwrapping it
	if emptyNode, ok := inode.(*ast.EmptyExprNode); ok {
		node = emptyNode.Child
	}
	
	switch node := node.(type) {
	case *ast.IntLiteralNode:
		return node.Value
	case *ast.ReferenceExprNode:
		val, ok := local_scope.getVar(node.Name)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s", node.Name))
		}
		intNode, ok := val.(*ast.IntLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected int, got %v", val))
		}
		return intNode.Value
	case *ast.InfixExprNode:
		switch node.Operator {
		case token.PLUS:
			return i.execIntExpr(node.Left, local_scope) + i.execIntExpr(node.Right, local_scope)
		case token.MINUS:
			return i.execIntExpr(node.Left, local_scope) - i.execIntExpr(node.Right, local_scope)
		case token.MULTIPLY:
			return i.execIntExpr(node.Left, local_scope) * i.execIntExpr(node.Right, local_scope)
		case token.DIVIDE:
			return i.execIntExpr(node.Left, local_scope) / i.execIntExpr(node.Right, local_scope)
		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown int expression: %v", node))
}

func (i *Interpreter) execBoolExpr(inode ast.Node, local_scope *Scope) bool {
	var node ast.Node = inode
	
	// Handle EmptyExprNode by unwrapping it
	if emptyNode, ok := inode.(*ast.EmptyExprNode); ok {
		node = emptyNode.Child
	}
	
	switch node := node.(type) {
	case *ast.BoolLiteralNode:
		return node.Value
	case *ast.PrefixExprNode:
		return !i.execBoolExpr(node.Value, local_scope)
	case *ast.ReferenceExprNode:
		val, ok := local_scope.getVar(node.Name)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s", node.Name))
		}
		boolNode, ok := val.(*ast.BoolLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected bool, got %v", val))
		}
		return boolNode.Value
	case *ast.BoolInfixNode:
		switch node.Operator {
		case token.AND:
			return i.execBoolExpr(node.Left, local_scope) && i.execBoolExpr(node.Right, local_scope)
		case token.OR:
			return i.execBoolExpr(node.Left, local_scope) || i.execBoolExpr(node.Right, local_scope)
		case token.LESS_THAN:
			return i.execIntExpr(node.Left, local_scope) < i.execIntExpr(node.Right, local_scope)
		case token.LESS_THAN_EQT:
			return i.execIntExpr(node.Left, local_scope) <= i.execIntExpr(node.Right, local_scope)
		case token.GREATER_THAN:
			return i.execIntExpr(node.Left, local_scope) > i.execIntExpr(node.Right, local_scope)
		case token.GREATER_THAN_EQT:
			return i.execIntExpr(node.Left, local_scope) >= i.execIntExpr(node.Right, local_scope)
		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown bool expression: %v", node))
}

func (i *Interpreter) assignValue(name string, value ast.Node, local_scope *Scope, isDeclaration bool) {

	var valNode ast.Node

	switch v := value.(type) {
	case *ast.IntLiteralNode, *ast.InfixExprNode:
		valNode = &ast.IntLiteralNode{Value: i.execIntExpr(v, local_scope)}
	case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
		valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(v, local_scope)}
	case *ast.ReferenceExprNode:
		// Look up the referenced variable
		refVal, ok := local_scope.getVar(v.Name)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s of type %v\n", v.Name, v.NodeType()))
		}
		// Evaluate based on the actual type of the referenced node
		switch refVal := refVal.(type) {
		case *ast.IntLiteralNode, *ast.InfixExprNode:
			valNode = &ast.IntLiteralNode{Value: i.execIntExpr(refVal, local_scope)}
		case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
			valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(refVal, local_scope)}
		default:
			panic(fmt.Sprintf("[ERROR] Unknown reference type: %T\n", refVal))
		}

	default:
		panic(fmt.Sprintf("[ERROR] Unknown value type: %v\n", value))
	}
	if valNode == nil {
		panic(fmt.Sprintf("[ERROR] Variable is undefined, %v\n", name))
	} else {
	}
	if isDeclaration {
		local_scope.declareVar(name, valNode)
	} else {
		if !local_scope.assignVar(name, valNode) {
			panic(fmt.Sprintf("[ERROR] Reassigning undefined variable %s", name))
		}
	}
}

func (i *Interpreter) changeVarVal(node ast.Node, local_scope *Scope) {
	switch n := node.(type) {
	case *ast.LetStmtNode:
		i.assignValue(n.Name, n.Value, local_scope, true)
	case *ast.VarReassignNode:
		varName := n.Var.Name
		i.assignValue(varName, n.NewVal, local_scope, false)
	default:
		panic(fmt.Sprintf("[ERROR] Unsupported node type: %T", node))
	}
}
func (i *Interpreter) execIfStmt(node ast.Node, local_scope *Scope) {
	ifStmt := node.(*ast.IfStmtNode)
	cond := i.execBoolExpr(ifStmt.Cond, local_scope)
	newScope := local_scope.newChild()
	if cond {
		fmt.Println("Processing body")
		for _, stmt := range ifStmt.Body {
			i.executeStmt(stmt, newScope)
		}
	} else {
		fmt.Println("Processing alt")
		for _, stmt := range ifStmt.Alt {
			i.executeStmt(stmt, newScope)
		}
	}
}

func (i *Interpreter) executeStmt(node ast.Node, local_scope *Scope) {
	switch node.NodeType() {
	case ast.IntLiteral, ast.InfixExpr:
		i.execIntExpr(node, local_scope)
	case ast.BoolLiteral, ast.BoolInfix, ast.PrefixExpr:
		i.execBoolExpr(node, local_scope)
	case ast.LetStmt, ast.VarReassign:
		i.changeVarVal(node, local_scope)
	case ast.IfStmt:
		i.execIfStmt(node, local_scope)
	default:
		panic(fmt.Sprintf("[ERROR] Unknown statement type: %v", node))
	}
}

func (i *Interpreter) Execute(program ast.ProgramNode, should_print bool) Scope {
	for _, stmt := range program.Statements {
		i.executeStmt(stmt, &i.MainScope)
	}
	if should_print {
		fmt.Printf("%v\n", i.MainScope);
	}
	return i.MainScope
}
