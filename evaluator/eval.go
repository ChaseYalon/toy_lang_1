package evaluator

import (
	"fmt"
	"strconv"
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


type Interpreter struct {
	MainScope Scope
}

func NewInterpreter() Interpreter {
	ms := Scope {
		Vars: make(v_map),
		Funcs: make(f_map),
	}
	ms.Funcs["print"] = ast.FuncDecNode{
		Name: "print",
		Params: []ast.ReferenceExprNode{ast.ReferenceExprNode{Name: "input"}},
		Body: []ast.Node{&ast.CallBuiltinNode{
			Name: "print",
			Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
		}},
	}
	ms.Funcs["println"] = ast.FuncDecNode{
		Name: "println",
		Params: []ast.ReferenceExprNode{ast.ReferenceExprNode{Name: "input"}},
		Body: []ast.Node{&ast.CallBuiltinNode{
			Name: "println",
			Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
		}},
	}
	return Interpreter{
		MainScope: ms,
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

func (i *Interpreter) execExpr(node ast.Node, local_scope *Scope) ast.Node {
	if node.NodeType() == ast.IntLiteral {
		return &ast.IntLiteralNode{Value: i.execIntExpr(node, local_scope)}
	}
	if node.NodeType() == ast.InfixExpr {
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
	if node.NodeType() == ast.EmptyExpr {
		emptExpr, ok := node.(*ast.EmptyExprNode)
		if !ok {
			panic(fmt.Sprintf("[ERROR] Could not parse empty node, got %v\n", node))
		}
		return i.execExpr(emptExpr.Child, local_scope)
	}
	if node.NodeType() == ast.ReferenceExpr{
		refNode, ok := node.(*ast.ReferenceExprNode);
		if !ok{
			panic(fmt.Sprintf("[ERROR] Could not convert %v to a variable reference\n", node));
		}
		local_var, found := local_scope.getVar(refNode.Name);
		if !found{
			panic(fmt.Sprintf("[ERROR] Variable \"%v\" not found in current scope\n", local_var));
		}
		return local_var;
	}
	panic(fmt.Sprintf("[ERROR] Could not figure out what to parse, got %v of type %v\n", node, node.NodeType()))
}
func (i *Interpreter) callBuiltin(inode ast.Node, local_scope *Scope) {
    node, ok := inode.(*ast.CallBuiltinNode)
    if !ok {
        panic(fmt.Sprintf("[ERROR] Invalid function call, got %v\n", inode))
    }
    if node.Name == "print" {
        if len(node.Params) != 1 {
            panic(fmt.Sprintf("[ERROR] Print must be called with 1 argument, got %v\n", node))
        }
        expr := i.execExpr(node.Params[0], local_scope)
        sExpr, ok := expr.(*ast.StringLiteralNode)
        if !ok {
            panic(fmt.Sprintf("[ERROR] Could not parse %v into a string\n", expr))
        }

        // Unescape the string so \n, \t, etc. work
        unescaped, err := strconv.Unquote(`"` + sExpr.Value + `"`)
        if err != nil {
            panic(fmt.Sprintf("[ERROR] Could not unescape string %v: %v", sExpr.Value, err))
        }

        fmt.Print(unescaped)
    }
	if node.Name == "println" {
        if len(node.Params) != 1 {
            panic(fmt.Sprintf("[ERROR] Print must be called with 1 argument, got %v\n", node))
        }
        expr := i.execExpr(node.Params[0], local_scope)
        sExpr, ok := expr.(*ast.StringLiteralNode)
        if !ok {
            panic(fmt.Sprintf("[ERROR] Could not parse %v into a string\n", expr))
        }

        // Unescape the string so \n, \t, etc. work
        unescaped, err := strconv.Unquote(`"` + sExpr.Value + `"`)
        if err != nil {
            panic(fmt.Sprintf("[ERROR] Could not unescape string %v: %v", sExpr.Value, err))
        }

        fmt.Print(unescaped + "\n")
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
	case ast.CallBuiltin:
		i.callBuiltin(node, local_scope);
	case ast.EmptyExpr:
		emptExpr, ok := node.(*ast.EmptyExprNode);
		if !ok{
			panic(fmt.Sprintf("WTF happend here, got %v\n", node));
		}
		i.executeStmt(emptExpr.Child, local_scope);
	default:
		panic(fmt.Sprintf("[ERROR] Unknown statement type: %v, of type %v\n", node, node.NodeType()))
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
