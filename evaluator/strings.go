package evaluator

import (
	"fmt"
	"strconv"
	"toy_lang/ast"
	"toy_lang/token"
)

func (i *Interpreter) execStringExpr(node ast.Node, local_scope *Scope) string {
	if node.NodeType() == ast.CallBuiltin {
		callNode, ok := node.(*ast.CallBuiltinNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to CallBuiltinNode, got %v\n", node))
		}
		res := i.callBuiltin(callNode, local_scope)
		sres, ok := res.(*ast.StringLiteralNode);
		if !ok{
			panic(fmt.Sprintf("[ERROR] Could not convert %v to string\n", res));
		}
		return sres.Value
	}

	if node.NodeType() == ast.StringLiteral {
		strNode, ok := node.(*ast.StringLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to string literal, got %v\n", node))
		}
		return strNode.Value
	}
	if node.NodeType() == ast.InfixExpr {
		infNode, ok := node.(*ast.InfixExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to infix expression, got %v\n", node))
		}
		if infNode.Operator != token.PLUS {
			panic(fmt.Sprintf("[ERROR] Only supported operator on strings is plus, got %v", infNode.Operator))
		}
		return i.execStringExpr(infNode.Left, local_scope) + i.execStringExpr(infNode.Right, local_scope)
	}
	if node.NodeType() == ast.ReferenceExpr {
		refNode, ok := node.(*ast.ReferenceExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to reference expression, got %v\n", node))
		}
		val, exists := local_scope.getVar(refNode.Name)
		if !exists {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s", refNode.Name))
		}
		strNode, ok := val.(*ast.StringLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Expected string, got %v", val))
		}
		return strNode.Value
	}
	if node.NodeType() == ast.IntLiteral {
		refNode, ok := node.(*ast.IntLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to reference expression, got %v\n", node))
		}
		str := strconv.Itoa(refNode.Value)
		return str
	}
	if node.NodeType() == ast.BoolLiteral {
		refNode, ok := node.(*ast.BoolLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not convert node to reference expression, got %v\n", node))
		}
		str := strconv.FormatBool(refNode.Value)
		return str
	}
	panic(fmt.Sprintf("[ERROR] Type unsupported for string operations, got %v of value %v\n", node.NodeType(), node))
}

func (i *Interpreter) isStringExpression(node ast.Node, local_scope *Scope) bool {
	switch node.NodeType() {
	case ast.StringLiteral:
		return true
	case ast.ReferenceExpr:
		refNode, ok := node.(*ast.ReferenceExprNode)
		if !ok {
			return false
		}
		val, exists := local_scope.getVar(refNode.Name)
		if !exists {
			return false
		}
		return val.NodeType() == ast.StringLiteral
	case ast.InfixExpr:
		infixNode, ok := node.(*ast.InfixExprNode)
		if !ok {
			return false
		}
		return i.isStringExpression(infixNode.Left, local_scope) || i.isStringExpression(infixNode.Right, local_scope)
	
	case ast.FuncCall:
		//THIS IS TEMPORARY PLACEHOLDER CODE!!!!!!!
		return true
	}
	return false
}
