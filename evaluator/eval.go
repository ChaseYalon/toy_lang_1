package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"toy_lang/ast"
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

// Wrapper to propagate return values
type ReturnValue struct {
	Val ast.Node
}

type Interpreter struct {
	MainScope Scope
	reader    *bufio.Reader
}

func NewInterpreter() Interpreter {
	builtinScope := Scope{
		Vars:  make(v_map),
		Funcs: make(f_map),
	}
	ms := builtinScope.newChild()

	// Builtin functions
	builtinScope.Funcs["print"] = ast.FuncDecNode{
		Name:   "print",
		Params: []ast.ReferenceExprNode{{Name: "input"}},
		Body: []ast.Node{
			&ast.CallBuiltinNode{
				Name:   "print",
				Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
			},
		},
	}
	builtinScope.Funcs["println"] = ast.FuncDecNode{
		Name:   "println",
		Params: []ast.ReferenceExprNode{{Name: "input"}},
		Body: []ast.Node{
			&ast.CallBuiltinNode{
				Name:   "println",
				Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
			},
		},
	}
	builtinScope.Funcs["input"] = ast.FuncDecNode{
		Name:   "input",
		Params: []ast.ReferenceExprNode{{Name: "prompt"}},
		Body: []ast.Node{
			&ast.ReturnExprNode{
				Val: &ast.CallBuiltinNode{
					Name:   "input",
					Params: []ast.Node{&ast.ReferenceExprNode{Name: "prompt"}},
				},
			},
		},
	}
	builtinScope.Funcs["int"] = ast.FuncDecNode{
		Name:   "int",
		Params: []ast.ReferenceExprNode{{Name: "convertToInt"}},
		Body: []ast.Node{
			&ast.ReturnExprNode{
				Val: &ast.CallBuiltinNode{
					Name:   "int",
					Params: []ast.Node{&ast.ReferenceExprNode{Name: "convertToInt"}},
				},
			},
		},
	}
	builtinScope.Funcs["bool"] = ast.FuncDecNode{
		Name:   "bool",
		Params: []ast.ReferenceExprNode{{Name: "convertToBool"}},
		Body: []ast.Node{
			&ast.ReturnExprNode{
				Val: &ast.CallBuiltinNode{
					Name:   "bool",
					Params: []ast.Node{&ast.ReferenceExprNode{Name: "convertToBool"}},
				},
			},
		},
	}
	builtinScope.Funcs["str"] = ast.FuncDecNode{
		Name:   "str",
		Params: []ast.ReferenceExprNode{{Name: "convertToStr"}},
		Body: []ast.Node{
			&ast.ReturnExprNode{
				Val: &ast.CallBuiltinNode{
					Name:   "str",
					Params: []ast.Node{&ast.ReferenceExprNode{Name: "convertToStr"}},
				},
			},
		},
	}

	return Interpreter{
		MainScope: *ms,
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (i *Interpreter) executeStmt(node ast.Node, local_scope *Scope) any {
	switch node.NodeType() {
	case ast.LetStmt, ast.VarReassign:
		i.changeVarVal(node, local_scope)
	case ast.IfStmt:
		return i.execIfStmt(node, local_scope)
	case ast.WhileStmt:
		return i.execWhileStmt(node, local_scope)
	case ast.FuncDec:
		local_scope.declareFunc(*node.(*ast.FuncDecNode))
	case ast.FuncCall:
		return i.execFuncCall(node, local_scope)
	case ast.CallBuiltin:
		return i.callBuiltin(node, local_scope)
	case ast.ReturnExpr: 
		returnNode := node.(*ast.ReturnExprNode)
		returnVal := i.execExpr(returnNode.Val, local_scope)
		return ReturnValue{Val: returnVal}
	case ast.EmptyExpr:
		child := node.(*ast.EmptyExprNode).Child
		return i.executeStmt(child, local_scope)
	default:
		// For expressions used as statements
		switch node.(type) {
		case *ast.IntLiteralNode, *ast.FloatLiteralNode, *ast.BoolLiteralNode, *ast.StringLiteralNode:
			i.execExpr(node, local_scope)
		case ast.Node:
			i.execExpr(node, local_scope)
		default:
			panic(fmt.Sprintf("[ERROR] Unknown statement type: %v, of type %v\n", node, node.NodeType()))
		}
	}
	return nil
}

func (i *Interpreter) execIfStmt(node ast.Node, local_scope *Scope) interface{} {
	ifStmt := node.(*ast.IfStmtNode)
	cond := i.execBoolExpr(ifStmt.Cond, local_scope)
	newScope := local_scope.newChild()

	if cond {
		for _, stmt := range ifStmt.Body {
			if ret := i.executeStmt(stmt, newScope); ret != nil {
				if r, ok := ret.(ReturnValue); ok {
					return r
				}
			}
		}
	} else {
		for _, stmt := range ifStmt.Alt {
			if ret := i.executeStmt(stmt, newScope); ret != nil {
				if r, ok := ret.(ReturnValue); ok {
					return r
				}
			}
		}
	}
	return nil
}

func (i *Interpreter) execWhileStmt(node ast.Node, local_scope *Scope) interface{} {
	whileStmt := node.(*ast.WhileStmtNode)

	for i.execBoolExpr(whileStmt.Cond, local_scope) {
		bodyScope := local_scope.newChild()
		for _, stmt := range whileStmt.Body {
			if stmt.NodeType() == ast.BreakSmt{
				goto EndOfOuter
			}
			if stmt.NodeType() == ast.ContinueStmt{
				goto EndOfInner
			}
			if ret := i.executeStmt(stmt, bodyScope); ret != nil {
				if r, ok := ret.(ReturnValue); ok {
					return r
				}
			}
		}
		// Apply body scope changes back
		for k, v := range bodyScope.Vars {
			local_scope.Vars[k] = v
		}
		EndOfInner:
	}
	EndOfOuter:
	return nil
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

	// Assign parameters
	for j, param := range f.Params {
		letStmt := &ast.LetStmtNode{Name: param.Name, Value: fCall.Params[j]}
		i.changeVarVal(letStmt, callScope)
	}

	// Execute function body
	for _, stmt := range f.Body {
		if ret := i.executeStmt(stmt, callScope); ret != nil {
			if r, ok := ret.(ReturnValue); ok {
				return r.Val
			}
		}
	}
	return nil
}

func (i *Interpreter) callBuiltin(node ast.Node, local_scope *Scope) ast.Node {
	inode := node.(*ast.CallBuiltinNode)

	switch inode.Name {
	case "print", "println":
		if len(inode.Params) != 1 {
			panic(fmt.Sprintf("[ERROR] Builtin %s must be called with 1 argument, got %v", inode.Name, inode))
		}
		val := i.execExpr(inode.Params[0], local_scope)
		var output string
		switch v := val.(type) {
		case *ast.StringLiteralNode:
			output = v.Value
		case *ast.IntLiteralNode:
			output = strconv.Itoa(v.Value)
		case *ast.BoolLiteralNode:
			output = strconv.FormatBool(v.Value)
		default:
			panic(fmt.Sprintf("[ERROR] Cannot print value of type %v", val))
		}
		if inode.Name == "print" {
			fmt.Print(output)
		} else {
			fmt.Println(output)
		}
		return &ast.StringLiteralNode{Value: ""}
	case "input":
		promptNode := &ast.CallBuiltinNode{Name: "print", Params: []ast.Node{inode.Params[0]}}
		i.callBuiltin(promptNode, local_scope)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("[ERROR] Could not read input: %v", err))
		}
		text = strings.TrimSuffix(text, "\r\n")
		text = strings.TrimSuffix(text, "\n")
		return &ast.StringLiteralNode{Value: text}
	case "str":
		toConv := i.execExpr(inode.Params[0], local_scope)
		switch t := toConv.(type) {
		case *ast.StringLiteralNode:
			return t
		case *ast.BoolLiteralNode:
			return &ast.StringLiteralNode{Value: strconv.FormatBool(t.Value)}
		case *ast.IntLiteralNode:
			return &ast.StringLiteralNode{Value: strconv.Itoa(t.Value)}
		default:
			panic(fmt.Sprintf("[ERROR] Cannot convert type %v to string", toConv.NodeType()))
		}
	case "int":
		toConv := i.execExpr(inode.Params[0], local_scope)
		switch t := toConv.(type) {
		case *ast.IntLiteralNode:
			return t
		case *ast.BoolLiteralNode:
			v := 0
			if t.Value {
				v = 1
			}
			return &ast.IntLiteralNode{Value: v}
		case *ast.StringLiteralNode:
			val, err := strconv.Atoi(t.Value)
			if err != nil {
				panic(fmt.Sprintf("[ERROR] Cannot convert string to int: %v", err))
			}
			return &ast.IntLiteralNode{Value: val}
		default:
			panic(fmt.Sprintf("[ERROR] Cannot convert type %v to int", toConv.NodeType()))
		}
	case "bool":
		toConv := i.execExpr(inode.Params[0], local_scope)
		switch t := toConv.(type) {
		case *ast.BoolLiteralNode:
			return t
		case *ast.IntLiteralNode:
			return &ast.BoolLiteralNode{Value: t.Value > 0}
		case *ast.StringLiteralNode:
			if t.Value == "" || t.Value == "false" {
				return &ast.BoolLiteralNode{Value: false}
			}
			return &ast.BoolLiteralNode{Value: true}
		default:
			panic(fmt.Sprintf("[ERROR] Cannot convert type %v to bool", toConv.NodeType()))
		}
	}
	panic(fmt.Sprintf("[ERROR] Unknown builtin function %v", inode.Name))
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
