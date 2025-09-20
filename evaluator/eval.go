package evaluator

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)

type v_map map[string]ast.Node
type f_map map[string]ast.FuncDecNode

func (v v_map) String() string {
	s := "{"
	for k, node := range v {
		s += k + ": " + fmt.Sprintf("%v", node) + ", "
	}
	s += "}"
	return s
}

type Scope struct {
	Vars   v_map
	Funcs  f_map
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

func (s *Scope) getFunc(name string) (ast.FuncDecNode, bool) {
	if val, ok := s.Funcs[name]; ok {
		return val, true
	}
	if s.Parent != nil {
		return s.Parent.getFunc(name)
	}
	return ast.FuncDecNode{}, false
}

func (s *Scope) declareFunc(f ast.FuncDecNode) {
	s.Funcs[f.Name] = f
}

func (s *Scope) declareVar(name string, val ast.Node) {
	s.Vars[name] = val
}

func (s *Scope) assignVar(name string, val ast.Node) bool {
	// Fixed: Added StringLiteral to the condition
	if val.NodeType() == ast.IntLiteral || val.NodeType() == ast.BoolLiteral || val.NodeType() == ast.StringLiteral {
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
		Funcs:  make(f_map),
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
			Vars:  make(v_map),
			Funcs: make(f_map),
		},
	}
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
		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown bool expression: %v", node))
}

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
		valNode = &ast.IntLiteralNode{Value: i.execIntExpr(v, local_scope)}
	case *ast.BoolLiteralNode, *ast.BoolInfixNode, *ast.PrefixExprNode:
		valNode = &ast.BoolLiteralNode{Value: i.execBoolExpr(v, local_scope)}
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
		default:
			panic(fmt.Sprintf("[ERROR] Unsupported return type from function: %T", r))
		}
	case *ast.StringLiteralNode:
		valNode = &ast.StringLiteralNode{Value: i.execStringExpr(v, local_scope)}
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

func (i *Interpreter) execIfStmt(node ast.Node, local_scope *Scope) {
	ifStmt := node.(*ast.IfStmtNode)
	cond := i.execBoolExpr(ifStmt.Cond, local_scope)
	newScope := local_scope.newChild()
	if cond {
		for _, stmt := range ifStmt.Body {
			i.executeStmt(stmt, newScope)
		}
	} else {
		for _, stmt := range ifStmt.Alt {
			i.executeStmt(stmt, newScope)
		}
	}
}

func (i *Interpreter) execFuncCall(node ast.Node, local_scope *Scope) ast.Node {
	fCall := node.(*ast.FuncCallNode)

	f, found := local_scope.getFunc(fCall.Name.Name)
	if !found {
		panic(fmt.Sprintf("[ERROR] Could not find function %s\n", fCall.Name.Name))
	}

	callScope := local_scope.newChild()

	if len(f.Params) != len(fCall.Params) {
		panic(fmt.Sprintf("[ERROR] Function %s must be called with exactly %d params, got %d\n",
			f.Name, len(f.Params), len(fCall.Params)))
	}

	for j, param := range f.Params {
		letStmt := &ast.LetStmtNode{Name: param.Name, Value: fCall.Params[j]}
		i.changeVarVal(letStmt, callScope)
	}

	for _, stmt := range f.Body {
		if stmt.NodeType() == ast.ReturnExpr {
			ret := stmt.(*ast.ReturnExprNode)
			return i.execExpr(ret.Val, callScope)
		}
		i.executeStmt(stmt, callScope)
	}

	return nil
}

func (i *Interpreter) execStringExpr(node ast.Node, local_scope *Scope) string {
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
	panic(fmt.Sprintf("[ERROR] Type unsupported for string operations, got %v\n", node.NodeType()))
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
		// For function calls, we need to execute and check the return type
		// This is a bit more complex, but for now we'll assume it could be a string
		// A more robust solution would involve type checking or trying to execute
		return true
	}
	return false
}

func (i *Interpreter) execExpr(node ast.Node, local_scope *Scope) ast.Node {
	if node.NodeType() == ast.IntLiteral {
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.InfixExpr {
		// Check if this is a string operation by examining the operands
		infixNode, ok := node.(*ast.InfixExprNode)
		if ok && infixNode.Operator == token.PLUS {
			// Check if either operand is a string
			if i.isStringExpression(infixNode.Left, local_scope) || i.isStringExpression(infixNode.Right, local_scope) {
				return &ast.StringLiteralNode{Value: i.execStringExpr(node, local_scope)}
			}
		}
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
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
	panic(fmt.Sprintf("[ERROR] Could not figure out what to parse, got %v of type %v\n", node, node.NodeType()))
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
	case ast.FuncDec:
		fNode, ok := node.(*ast.FuncDecNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not find function %v\n", node))
		}
		local_scope.declareFunc(*fNode)
	case ast.FuncCall:
		i.execFuncCall(node, local_scope)
	case ast.StringLiteral:
		i.execStringExpr(node, local_scope)
	default:
		panic(fmt.Sprintf("[ERROR] Unknown statement type: %v", node))
	}
}

func (i *Interpreter) Execute(program ast.ProgramNode, should_print bool) Scope {
	for _, stmt := range program.Statements {
		i.executeStmt(stmt, &i.MainScope)
	}
	if should_print {
		fmt.Printf("Main scope: %v\n", i.MainScope)
	}
	return i.MainScope
}
