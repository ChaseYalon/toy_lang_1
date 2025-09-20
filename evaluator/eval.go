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

type Interpreter struct {
	MainScope Scope
	reader    *bufio.Reader
}

func NewInterpreter() Interpreter {
	ms := Scope{
		Vars:  make(v_map),
		Funcs: make(f_map),
	}
	i := Interpreter{
		MainScope: ms,
		reader:    bufio.NewReader(os.Stdin),
	}

	// Builtin: print
	i.MainScope.Funcs["print"] = ast.FuncDecNode{
		Name:   "print",
		Params: []ast.ReferenceExprNode{{Name: "input"}},
		Body: []ast.Node{
			&ast.CallBuiltinNode{
				Name:   "print",
				Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
			},
		},
	}

	// Builtin: println
	i.MainScope.Funcs["println"] = ast.FuncDecNode{
		Name:   "println",
		Params: []ast.ReferenceExprNode{{Name: "input"}},
		Body: []ast.Node{
			&ast.CallBuiltinNode{
				Name:   "println",
				Params: []ast.Node{&ast.ReferenceExprNode{Name: "input"}},
			},
		},
	}

	// Builtin: input
	ms.Funcs["input"] = ast.FuncDecNode{
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

	return i
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

func (i *Interpreter) callBuiltin(inode ast.Node, local_scope *Scope) ast.StringLiteralNode {
	node := inode.(*ast.CallBuiltinNode)

	switch node.Name {
	case "print", "println":
		if len(node.Params) != 1 {
			panic(fmt.Sprintf("[ERROR] Builtin %s must be called with 1 argument, got %v", node.Name, node))
		}
		val := i.execExpr(node.Params[0], local_scope)

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

		if node.Name == "print" {
			fmt.Print(output)
		} else {
			fmt.Println(output)
		}
		return ast.StringLiteralNode{Value: ""}

	case "input":
		// print prompt first
		promptNode := &ast.CallBuiltinNode{Name: "print", Params: []ast.Node{node.Params[0]}}
		i.callBuiltin(promptNode, local_scope)

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("[ERROR] Could not read input: %v", err))
		}
		text = strings.TrimSuffix(text, "\r\n")
		text = strings.TrimSuffix(text, "\n")
		return ast.StringLiteralNode{Value: text}

	}

	panic(fmt.Sprintf("[ERROR] Unknown builtin function %v", node))
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
		i.callBuiltin(node, local_scope)
	case ast.EmptyExpr:
		emptExpr, ok := node.(*ast.EmptyExprNode)
		if !ok {
			panic(fmt.Sprintf("WTF happend here, got %v\n", node))
		}
		i.executeStmt(emptExpr.Child, local_scope)
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
