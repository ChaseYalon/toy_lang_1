package evaluator

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/token"
)
type v_map  map[string]ast.Node;

func (v v_map) String() string {
    s := "{"
    for k, node := range v {
        s += k + ": " + fmt.Sprintf("%v", node) + ", "
    }
    s += "}"
    return s
}

type Interpreter struct {
	vars v_map
}

func NewInterpreter() Interpreter {
	return Interpreter{
		vars: make(v_map),
	}
}

func (i *Interpreter) execExpr(node ast.Node)int  {
	switch n := node.(type) {		
	case *ast.IntLiteralNode:
		return n.Value
	case *ast.ReferenceExprNode:
		val := i.vars[n.Name]
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

func (i *Interpreter) execBoolExpr(node ast.Node)bool{
	if node.NodeType() == ast.BoolLiteral{
		return node.(*ast.BoolLiteralNode).Value;
	}
	if node.NodeType() == ast.ReferenceExpr{
		return i.execBoolExpr(i.vars[node.(*ast.ReferenceExprNode).Name]);
	}
	if node.NodeType() == ast.BoolInfix{
		execNode := node.(*ast.BoolInfixNode);
		if execNode.Operator == token.OR{
			return i.execBoolExpr(execNode.Left) || i.execBoolExpr(execNode.Right);
		}
		if execNode.Operator == token.AND{
			return i.execBoolExpr(execNode.Left) && i.execBoolExpr(execNode.Right);
		}
		if execNode.Operator == token.EQUALS{
			return i.executeStmt(execNode.Left) == i.executeStmt(execNode.Right);
		}
		if execNode.Operator == token.LESS_THAN{
			return i.execExpr(execNode.Left) < i.execExpr(execNode.Right);
		}
		if execNode.Operator == token.GREATER_THAN{
			return i.execExpr(execNode.Left) > i.execExpr(execNode.Right);
		}
		if execNode.Operator == token.LESS_THAN_EQT{
			return i.execExpr(execNode.Left) <= i.execExpr(execNode.Right);
		}
		if execNode.Operator == token.GREATER_THAN_EQT{
			return i.execExpr(execNode.Left) >= i.execExpr(execNode.Right);
		}
	}
	if node.NodeType() == ast.PrefixExpr{
		execNode := node.(*ast.PrefixExprNode);
		return !i.execBoolExpr(execNode.Value);
	}
	panic(fmt.Sprintf("[ERROR] Invalid expression: %v of type %v", node, node.NodeType().String()));
}
func (i *Interpreter) executeStmt(stmt ast.Node) ast.Node {
	switch n := stmt.(type) {
	case *ast.VarReassignNode:
		//Tests for int literal node
		if n.NewVal.NodeType() == ast.IntLiteral{
			i.vars[n.Var.Name] = n.NewVal;
			return nil;
		}
		if n.NewVal.NodeType() == ast.BoolLiteral{
			i.vars[n.Var.Name] = n.NewVal;
			return nil;
		}
		if n.NewVal.NodeType() == ast.BoolInfix{
			i.vars[n.Var.Name] = &ast.BoolLiteralNode{Value: i.execBoolExpr(n.NewVal)};
			return nil;
		}

		//Default value: Infix expressions
		newVal := i.execExpr(n.NewVal)
		i.vars[n.Var.Name] = &ast.IntLiteralNode{Value: newVal}
		return nil
	case *ast.LetStmtNode:
		if n.Value.NodeType() == ast.BoolLiteral{
			boolNode, ok := n.Value.(*ast.BoolLiteralNode);
			if !ok{
				panic(fmt.Sprintf("[ERROR] Could not coerce %v to boolean literal", n.Value));
			}
			i.vars[n.Name] = boolNode;
			return nil;
		} else if(n.Value.NodeType() == ast.BoolInfix){
			i.vars[n.Name] = &ast.BoolLiteralNode{Value: i.execBoolExpr(n.Value)};
			return nil;
		}else{
			// Evaluate the initial value
			val := i.execExpr(n.Value)
			i.vars[n.Name] = &ast.IntLiteralNode{Value: val}
			return nil
		}
	case *ast.BoolInfixNode:
		return &ast.BoolLiteralNode{Value: i.execBoolExpr(n)};
	case *ast.InfixExprNode, *ast.IntLiteralNode:
		val := i.execExpr(n)
		return &ast.IntLiteralNode{Value: val}
	case *ast.ReferenceExprNode:
		if i.vars[n.Name].NodeType() == ast.BoolLiteral{
			return i.vars[n.Name];
		}
		//Assume it is an it
		val := i.execExpr(n);
		return &ast.IntLiteralNode{Value: val};
	}
	panic(fmt.Sprintf("[ERROR] Unknown statement type: %v", stmt))
}

func (i *Interpreter) Execute(program ast.ProgramNode, shouldPrint bool)v_map {
	for _, stmt := range program.Statements {
		i.executeStmt(stmt)
	}
	if shouldPrint{

		fmt.Printf("Vars:\n%v\n", i.vars);
	}
	return i.vars;
}