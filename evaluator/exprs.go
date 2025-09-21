package evaluator

import (
	"fmt"
	"math"
	"toy_lang/ast"
	"toy_lang/token"
)

func intPow(x, y int) int {
	if y < 0 {
		panic("negative exponent not supported for integers")
	}
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}
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
		case token.MODULO:
			return i.execIntExpr(node.Left, local_scope) % i.execIntExpr(node.Right, local_scope)
		case token.EXPONENT:
			return intPow(i.execIntExpr(node.Left, local_scope), i.execIntExpr(node.Right, local_scope))
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

func (i *Interpreter) execFloatExpr(inode ast.Node, local_scope *Scope) float64{
	var node ast.Node = inode

	if emptyNode, ok := inode.(*ast.EmptyExprNode); ok {
		node = emptyNode.Child
	}

	switch node := node.(type) {
	case *ast.FloatLiteralNode:
		return node.Value
	case *ast.IntLiteralNode:
		intNode, ok := inode.(*ast.IntLiteralNode);
		if !ok{
			panic(fmt.Sprintf("[ERROR] Could not convert to int, got %v\n", inode));
		}
		return float64(intNode.Value);
	case *ast.ReferenceExprNode:
		val, ok := local_scope.getVar(node.Name)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s", node.Name))
		}
		intNode, ok := val.(*ast.FloatLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected int, got %v", val))
		}
		return intNode.Value
	case *ast.FuncCallNode:
		result := i.execFuncCall(node, local_scope)
		if result == nil {
			panic(fmt.Sprintf("[ERROR] Function %s did not return a value", node.Name.Name))
		}
		intNode, ok := result.(*ast.FloatLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected int return from function, got %v", result))
		}
		return intNode.Value
	case *ast.InfixExprNode:
		switch node.Operator {
		case token.PLUS:
			return i.execFloatExpr(node.Left, local_scope) + i.execFloatExpr(node.Right, local_scope)
		case token.MINUS:
			return i.execFloatExpr(node.Left, local_scope) - i.execFloatExpr(node.Right, local_scope)
		case token.MULTIPLY:
			return i.execFloatExpr(node.Left, local_scope) * i.execFloatExpr(node.Right, local_scope)
		case token.DIVIDE:
			return i.execFloatExpr(node.Left, local_scope) / i.execFloatExpr(node.Right, local_scope)
		case token.MODULO:
			return math.Mod(i.execFloatExpr(node.Left, local_scope), i.execFloatExpr(node.Right, local_scope));
		case token.EXPONENT:
			return math.Pow(i.execFloatExpr(node.Left, local_scope), i.execFloatExpr(node.Right, local_scope))
		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown float expression: %v", node))
}

func (i *Interpreter) isFloatExpr(node ast.InfixExprNode, local_scope *Scope) bool {
	leftIsFloat := false
	rightIsFloat := false

	// direct float literals
	if node.Left.NodeType() == ast.FloatLiteral {
		leftIsFloat = true
	}
	if node.Right.NodeType() == ast.FloatLiteral {
		rightIsFloat = true
	}

	// variable references
	if node.Left.NodeType() == ast.ReferenceExpr {
		if refExpr, ok := node.Left.(*ast.ReferenceExprNode); ok {
			if lVar, found := local_scope.getVar(refExpr.Name); found {
				if lVar.NodeType() == ast.FloatLiteral {
					leftIsFloat = true
				}
			}
		}
	}
	if node.Right.NodeType() == ast.ReferenceExpr {
		if refExpr, ok := node.Right.(*ast.ReferenceExprNode); ok {
			if rVar, found := local_scope.getVar(refExpr.Name); found {
				if rVar.NodeType() == ast.FloatLiteral {
					rightIsFloat = true
				}
			}
		}
	}

	// recurse only when the child is another infix expression
	if !rightIsFloat {
		if rightNode, ok := node.Right.(*ast.InfixExprNode); ok {
			rightIsFloat = i.isFloatExpr(*rightNode, local_scope)
		}
	}
	if !leftIsFloat {
		if leftNode, ok := node.Left.(*ast.InfixExprNode); ok {
			leftIsFloat = i.isFloatExpr(*leftNode, local_scope)
		}
	}

	// expression is float if either side is float
	return leftIsFloat || rightIsFloat
}


func (i *Interpreter) execExpr(node ast.Node, local_scope *Scope) ast.Node {
	if node.NodeType() == ast.IntLiteral {
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.FloatLiteral{
		return &ast.FloatLiteralNode{Value: i.execFloatExpr(node, local_scope)};
	}
	if node.NodeType() == ast.InfixExpr {
		infixNode, ok := node.(*ast.InfixExprNode)
		if i.isStringExpression(infixNode.Left, local_scope) || i.isStringExpression(infixNode.Right, local_scope) && (ok && infixNode.Operator == token.PLUS ) {
			return &ast.StringLiteralNode{Value: i.execStringExpr(node, local_scope)}
		}
		if i.isFloatExpr(*infixNode, local_scope){
			return &ast.FloatLiteralNode{Value: i.execFloatExpr(infixNode, local_scope)};
		}
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.CallBuiltin {
		res := i.callBuiltin(node, local_scope)
		if res.NodeType() == ast.BoolLiteral {

			bRes, ok := res.(*ast.BoolLiteralNode)
			if !ok {
				panic(fmt.Sprintf("[ERROR] Expected bool return, got %v\n", res))
			}
			return bRes
		}
		if res.NodeType() == ast.IntLiteral {
			ires, ok := res.(*ast.IntLiteralNode)
			if !ok {
				panic(fmt.Sprintf("[ERROR] Expected int return, got %v\n", res))
			}
			return ires
		}
		if res.NodeType() == ast.StringLiteral {
			sRes, ok := res.(*ast.StringLiteralNode)
			if !ok {
				panic(fmt.Sprintf("[ERROR] Expected string return, got %v\n", res))
			}
			return sRes
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
	if node.NodeType() == ast.CallBuiltin {
		res := i.callBuiltin(node, local_scope)

		return res
	}

	panic(fmt.Sprintf("[ERROR] Could not figure out what to evaluate, got %v of type %v\n", node, node.NodeType()))
}
