package evaluator

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

func (i *Interpreter) assignValue(name string, value ast.Node, local_scope *Scope, isDeclaration bool) {
	if emptyNode, ok := value.(*ast.EmptyExprNode); ok {
		value = emptyNode.Child
	}
	var valNode ast.Node

	// Check if this is a string operation first
	if value.NodeType() == ast.InfixExpr {
		iNode, ok := value.(*ast.InfixExprNode)
		if ok && iNode.Operator == token.PLUS {
			// Check if either operand suggests this is a string operation
			if i.isStringExpression(iNode.Left, local_scope) || i.isStringExpression(iNode.Right, local_scope) {
				valNode = &ast.StringLiteralNode{Value: i.execStringExpr(value, local_scope)}
				goto AfterStringParse
			}
		}
	}

	switch v := value.(type) {
	case *ast.IntLiteralNode:
		valNode = &ast.IntLiteralNode{Value: i.execIntExpr(v, local_scope)}
	case *ast.InfixExprNode:
		// If we get here, it's not a string operation, so it must be int
		valNode = i.execExpr(v, local_scope)
	case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
		valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(v, local_scope)}
	case *ast.FloatLiteralNode:
		valNode = v
	case *ast.ReferenceExprNode:
		refVal, ok := local_scope.getVar(v.Name)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Undefined variable %s of type %v\n", v.Name, v.NodeType()))
		}
		switch refVal := refVal.(type) {
		case *ast.IntLiteralNode, *ast.InfixExprNode:
			valNode = &ast.IntLiteralNode{Value: i.execIntExpr(refVal, local_scope)}
		case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
			valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(refVal, local_scope)}
		case *ast.StringLiteralNode:
			valNode = &ast.StringLiteralNode{Value: i.execStringExpr(refVal, local_scope)}
		default:
			panic(fmt.Sprintf("[ERROR] Unknown reference type: %T\n", refVal))
		}
	case *ast.FuncCallNode:
		result := i.execFuncCall(v, local_scope)
		if result == nil {
			panic(fmt.Sprintf("[ERROR] Function %s did not return a value", v.Name.Name))
		}
		switch r := result.(type) {
		case *ast.IntLiteralNode, *ast.InfixExprNode:
			valNode = &ast.IntLiteralNode{Value: i.execIntExpr(r, local_scope)}
		case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
			valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(r, local_scope)}
		case *ast.StringLiteralNode:
			valNode = &ast.StringLiteralNode{Value: i.execStringExpr(r, local_scope)}
		case *ast.FloatLiteralNode:
			valNode = &ast.FloatLiteralNode{Value: i.execFloatExpr(r, local_scope)}
		default:
			panic(fmt.Sprintf("[ERROR] Unsupported return type from function: %T", r))
		}
	case *ast.StringLiteralNode:
		valNode = &ast.StringLiteralNode{Value: i.execStringExpr(v, local_scope)}
	case *ast.ArrLiteralNode:
		arrNode := value.(*ast.ArrLiteralNode)
		elems := make(map[ast.Node]ast.Node)
		for key, val := range arrNode.Elems {
			elems[key] = i.execExpr(val, local_scope)
		}
		valNode = &ast.ArrLiteralNode{Elems: elems}
	case *ast.ArrRefNode:
		refNode := value.(*ast.ArrRefNode)
		arr, found := local_scope.getVar(refNode.Arr.Name)
		if !found {
			panic(fmt.Sprintf("[ERROR] Cannot reference undefined variable, got %v\n", name))
		}
		arrMap, ok := arr.(*ast.ArrLiteralNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Variable %v is not an array, it is a %v\n", name, arr.NodeType()))
		}

		// Find the value by comparing string representations of keys
		var val ast.Node
		idxStr := i.execExpr(refNode.Idx, local_scope).String()
		for key, value := range arrMap.Elems {
			if key.String() == idxStr {
				val = value
				break
			}
		}

		if val == nil {
			panic(fmt.Sprintf("[ERROR] Value %v not found in arr %+v\n", refNode.Idx, arrMap))
		}
		valNode = val
	default:
		panic(fmt.Sprintf("[ERROR] Unknown value type: %v, type: %v\n", value, value.NodeType()))
	}
AfterStringParse:
	if valNode == nil {
		panic(fmt.Sprintf("[ERROR] Variable is undefined, %v\n", name))
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
		if n.Value.NodeType() == ast.LetStmt {
			lNode, ok := n.Value.(*ast.LetStmtNode)
			if !ok {
				panic(fmt.Sprintf("[ERROR] WTF happen with this let statement, got %v\n", node))
			}
			i.assignValue(n.Name, lNode.Value, local_scope, true)
		} else {
			i.assignValue(n.Name, n.Value, local_scope, true)
		}
	case *ast.VarReassignNode:
		varName := n.Var.Name
		i.assignValue(varName, n.NewVal, local_scope, false)
	default:
		panic(fmt.Sprintf("[ERROR] Unsupported node type: %T", node))
	}
}
