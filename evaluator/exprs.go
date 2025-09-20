package evaluator

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

func (i *Interpreter) execIntExpr(inode ast.Node, local_scope *Scope) int {
	var node ast.Node = inode

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
	case *ast.FuncCallNode:
		result := i.execFuncCall(node, local_scope)
		if result == nil {
			panic(fmt.Sprintf("[ERROR] Function %s did not return a value", node.Name.Name))
		}
		intNode, ok := result.(*ast.IntLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected int return from function, got %v", result))
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
	case *ast.FuncCallNode:
		result := i.execFuncCall(node, local_scope)
		if result == nil {
			panic(fmt.Sprintf("[ERROR] Function %s did not return a value", node.Name.Name))
		}
		boolNode, ok := result.(*ast.BoolLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected bool return from function, got %v", result))
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
		case token.EQUALS:
			leftVal := i.execExpr(node.Left, local_scope)
			rightVal := i.execExpr(node.Right, local_scope)

			switch l := leftVal.(type) {
			case *ast.IntLiteralNode:
				r, ok := rightVal.(*ast.IntLiteralNode)
				if !ok {
					return false
				}
				return l.Value == r.Value
			case *ast.BoolLiteralNode:
				r, ok := rightVal.(*ast.BoolLiteralNode)
				if !ok {
					return false
				}
				return l.Value == r.Value
			case *ast.StringLiteralNode:
				r, ok := rightVal.(*ast.StringLiteralNode)
				if !ok {
					return false
				}
				return l.Value == r.Value
			default:
				return false
			}

		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown bool expression: %v", node))
}

func (i *Interpreter) execExpr(node ast.Node, local_scope *Scope) ast.Node {
	if node.NodeType() == ast.IntLiteral {
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.InfixExpr {
		infixNode, ok := node.(*ast.InfixExprNode)
		if ok && infixNode.Operator == token.PLUS {
			if i.isStringExpression(infixNode.Left, local_scope) || i.isStringExpression(infixNode.Right, local_scope) {
				return &ast.StringLiteralNode{Value: i.execStringExpr(node, local_scope)}
			}
			return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
		}
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.CallBuiltin {
		res := i.callBuiltin(node, local_scope)
		if res.NodeType() == ast.BoolLiteral{

			bres, ok := res.(*ast.BoolLiteralNode);
			if !ok{
				panic(fmt.Sprintf("[ERROR] Expected bool return, got %v\n", res));
			}
			return bres 
		}
		if res.NodeType() == ast.IntLiteral{
			ires, ok := res.(*ast.IntLiteralNode);
			if !ok{
				panic(fmt.Sprintf("[ERROR] Expected int return, got %v\n", res));
			}
			return ires 
		}
		if res.NodeType() == ast.StringLiteral{
			sres, ok := res.(*ast.StringLiteralNode);
			if !ok{
				panic(fmt.Sprintf("[ERROR] Expected string return, got %v\n", res));
			}
			return sres 
		}
	}

	if node.NodeType() == ast.BoolLiteral || node.NodeType() == ast.BoolInfix || node.NodeType() == ast.PrefixExpr {
		return &ast.BoolLiteralNode{Value: i.execBoolExpr(node, local_scope)}
	}
	if node.NodeType() == ast.FuncCall {
		return i.execFuncCall(node, local_scope)
	}
	if node.NodeType() == ast.StringLiteral {
		return &ast.StringLiteralNode{Value: i.execStringExpr(node, local_scope)}
	}
	if node.NodeType() == ast.EmptyExpr {
		emptExpr, ok := node.(*ast.EmptyExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not parse empty node, got %v\n", node))
		}
		return i.execExpr(emptExpr.Child, local_scope)
	}
	if node.NodeType() == ast.ReferenceExpr {
		refNode, ok := node.(*ast.ReferenceExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert %v to a variable reference\n", node))
		}
		local_var, found := local_scope.getVar(refNode.Name)
		if !found {
			panic(fmt.Sprintf("[ERROR] Variable \"%v\" not found in current scope\n", local_var))
		}
		return local_var
	}
	if node.NodeType() == ast.CallBuiltin{
		res := i.callBuiltin(node, local_scope);

		return res;
	}
	panic(fmt.Sprintf("[ERROR] Could not figure out what to parse, got %v of type %v\n", node, node.NodeType()))
}
