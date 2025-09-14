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

func (s *Scope) getVar(name string) (ast.Node, bool) {
	if val, ok := s.Vars[name]; ok {
		return val, true
	}
	if s.Parent != nil {
		return s.Parent.getVar(name)
	}
	return nil, false
}

func (s *Scope) setVar(name string, value ast.Node) {
	if _, ok := s.Vars[name]; ok {
		s.Vars[name] = value
		return
	}
	if s.Parent != nil {
		if _, ok := s.Parent.getVar(name); ok {
			s.Parent.setVar(name, value)
			return
		}
	}
	// Default: declare in current scope
	s.Vars[name] = value
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

func (i *Interpreter) newScope(parent *Scope) Scope {
	return Scope{
		Vars:   make(v_map),
		Parent: parent,
	}
}

func (i *Interpreter) execExpr(node ast.Node) int {
	switch n := node.(type) {
	case *ast.IntLiteralNode:
		return n.Value
	case *ast.ReferenceExprNode:
		val, _ := i.MainScope.getVar(n.Name)
		return i.execExpr(val)
	case *ast.InfixExprNode:
		left := i.execExpr(n.Left)
		right := i.execExpr(n.Right)
		switch n.Operator {
		case token.PLUS:
			return left + right
		case token.MINUS:
			return left - right
		case token.MULTIPLY:
			return left * right
		case token.DIVIDE:
			return left / right
		default:
			panic(fmt.Sprintf("[ERROR] Unknown operator: %v", n.Operator))
		}
	}
	panic(fmt.Sprintf("[ERROR] Invalid expression: %v of type %v", node, node.NodeType().String()))
}

func (i *Interpreter) execBoolExpr(node ast.Node, parentScope *Scope, currentScope *Scope) bool {
	switch node.NodeType() {
	case ast.BoolLiteral:
		return node.(*ast.BoolLiteralNode).Value
	case ast.ReferenceExpr:
		val, _ := currentScope.getVar(node.(*ast.ReferenceExprNode).Name)
		return i.execBoolExpr(val, parentScope, currentScope)
	case ast.BoolInfix:
		execNode := node.(*ast.BoolInfixNode)
		switch execNode.Operator {
		case token.OR:
			return i.execBoolExpr(execNode.Left, parentScope, currentScope) ||
				i.execBoolExpr(execNode.Right, parentScope, currentScope)
		case token.AND:
			return i.execBoolExpr(execNode.Left, parentScope, currentScope) &&
				i.execBoolExpr(execNode.Right, parentScope, currentScope)
		case token.EQUALS:
			left := i.executeStmt(execNode.Left, parentScope, currentScope)
			right := i.executeStmt(execNode.Right, parentScope, currentScope)
			switch l := left.(type) {
			case *ast.IntLiteralNode:
				if r, ok := right.(*ast.IntLiteralNode); ok {
					return l.Value == r.Value
				}
			case *ast.BoolLiteralNode:
				if r, ok := right.(*ast.BoolLiteralNode); ok {
					return l.Value == r.Value
				}
			}
			panic("[ERROR] Invalid equality comparison")
		case token.LESS_THAN:
			return i.execExpr(execNode.Left) < i.execExpr(execNode.Right)
		case token.GREATER_THAN:
			return i.execExpr(execNode.Left) > i.execExpr(execNode.Right)
		case token.LESS_THAN_EQT:
			return i.execExpr(execNode.Left) <= i.execExpr(execNode.Right)
		case token.GREATER_THAN_EQT:
			return i.execExpr(execNode.Left) >= i.execExpr(execNode.Right)
		}
	case ast.PrefixExpr:
		execNode := node.(*ast.PrefixExprNode)
		return !i.execBoolExpr(execNode.Value, parentScope, currentScope)
	}
	panic(fmt.Sprintf("[ERROR] Invalid boolean expression: %v of type %v", node, node.NodeType().String()))
}

func (i *Interpreter) executeStmt(stmt ast.Node, parentScope *Scope, currentScope *Scope) ast.Node {
	switch n := stmt.(type) {
	case *ast.VarReassignNode:
		if n.NewVal.NodeType() == ast.IntLiteral || n.NewVal.NodeType() == ast.BoolLiteral {
			currentScope.setVar(n.Var.Name, n.NewVal)
			return nil
		}
		if n.NewVal.NodeType() == ast.BoolInfix {
			currentScope.setVar(n.Var.Name,
				&ast.BoolLiteralNode{Value: i.execBoolExpr(n.NewVal, parentScope, currentScope)})
			return nil
		}
		newVal := i.execExpr(n.NewVal)
		currentScope.setVar(n.Var.Name, &ast.IntLiteralNode{Value: newVal})
		return nil

	case *ast.LetStmtNode:
		switch n.Value.NodeType() {
		case ast.BoolLiteral:
			boolNode, ok := n.Value.(*ast.BoolLiteralNode)
			if !ok {
				panic(fmt.Sprintf("[ERROR] Could not coerce %v to boolean literal", n.Value))
			}
			currentScope.setVar(n.Name, boolNode)
			return nil
		case ast.BoolInfix, ast.PrefixExpr:
			currentScope.setVar(n.Name,
				&ast.BoolLiteralNode{Value: i.execBoolExpr(n.Value, parentScope, currentScope)})
			return nil
		default:
			val := i.execExpr(n.Value)
			currentScope.setVar(n.Name, &ast.IntLiteralNode{Value: val})
			return nil
		}

	case *ast.BoolInfixNode, *ast.BoolLiteralNode:
		return &ast.BoolLiteralNode{Value: i.execBoolExpr(n, parentScope, currentScope)}

	case *ast.InfixExprNode, *ast.IntLiteralNode:
		val := i.execExpr(n)
		return &ast.IntLiteralNode{Value: val}

	case *ast.ReferenceExprNode:
		bvar, _ := currentScope.getVar(n.Name)
		if bvar.NodeType() == ast.BoolLiteral {
			return bvar
		}
		val := i.execExpr(n)
		return &ast.IntLiteralNode{Value: val}

	case *ast.IfStmtNode:
		localScope := i.newScope(currentScope)
		boolCond := i.execBoolExpr(n.Cond, parentScope, currentScope)
		if boolCond {
			for _, val := range n.Body {
				i.executeStmt(val, &localScope, &localScope)
			}
		} else if len(n.Alt) != 0 {
			for _, val := range n.Alt {
				i.executeStmt(val, &localScope, &localScope)
			}
		}
		return nil
	}

	panic(fmt.Sprintf("[ERROR] Unknown statement type: %v\n", stmt))
}

func (i *Interpreter) Execute(program ast.ProgramNode, shouldPrint bool) v_map {
	for _, stmt := range program.Statements {
		i.executeStmt(stmt, &i.MainScope, &i.MainScope)
	}
	if shouldPrint {
		fmt.Printf("Vars:\n%v\n", i.MainScope.Vars)
	}
	return i.MainScope.Vars
}
