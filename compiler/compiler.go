package compiler

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/bytecode"
	"toy_lang/token"
)

type Compiler struct {
	ins         []bytecode.Instruction
	currBestAdr int
}

func NewCompiler() *Compiler {
	return &Compiler{
		ins:         []bytecode.Instruction{},
		currBestAdr: 0,
	}
}

func (c *Compiler) emit(ins bytecode.Instruction) {
	c.ins = append(c.ins, ins)
}
func (c *Compiler) compileExpr(node ast.Node) int {
	if node.NodeType() == ast.IntLiteral {
		intNode := node.(*ast.IntLiteralNode)

		toRet := bytecode.LOAD_INT_INS{
			Address: c.currBestAdr,
			Value:   intNode.Value,
		}
		c.currBestAdr++
		c.emit(&toRet)
		return c.currBestAdr - 1
	}
	if node.NodeType() == ast.InfixExpr {
		infixNode := node.(*ast.InfixExprNode)

		leftAddr := c.compileExpr(infixNode.Left)
		rightAddr := c.compileExpr(infixNode.Right)

		opInstr := -1
		switch infixNode.Operator {
		case token.PLUS:
			opInstr = 1
		case token.MINUS:
			opInstr = 2
		case token.MULTIPLY:
			opInstr = 3
		case token.DIVIDE:
			opInstr = 4
		default:
			panic(fmt.Sprintf("[ERROR] %v is not a valid infix operator", infixNode.Operator))
		}
		toRet := bytecode.INFIX_INT_INS{
			Left_addr:    leftAddr,
			Right_addr:   rightAddr,
			Save_to_addr: c.currBestAdr,
			Operation:    opInstr,
		}
		c.currBestAdr++
		c.emit(&toRet)
		return c.currBestAdr - 1
	}
	if node.NodeType() == ast.ReferenceExpr {
		refExpr := node.(*ast.ReferenceExprNode)
		toRet := bytecode.REF_VAR_INS{
			Name:   refExpr.Name,
			SaveTo: c.currBestAdr,
		}
		c.currBestAdr++
		c.emit(&toRet)
		return c.currBestAdr - 1
	}
	panic(fmt.Sprintf("[ERROR] Got unknown type of %v\n", node.NodeType()))
}

func (c *Compiler) compileStmt(node ast.Node) {
	if node.NodeType() == ast.InfixExpr || node.NodeType() == ast.IntLiteral {
		c.compileExpr(node)
		return

	}

	if node.NodeType() == ast.LetStmt {
		letStmt := node.(*ast.LetStmtNode)
		valAddr := c.compileExpr(letStmt.Value)
		toEmit := bytecode.DECLARE_VAR_INS{
			Name: letStmt.Name,
			Addr: valAddr,
		}
		c.emit(&toEmit)
		return
	}
	if node.NodeType() == ast.VarReassign {
		letStmt := node.(*ast.VarReassignNode)
		valAddr := c.compileExpr(letStmt.NewVal)
		toEmit := bytecode.DECLARE_VAR_INS{
			Name: letStmt.Var.Name,
			Addr: valAddr,
		}
		c.emit(&toEmit)
		return
	}
	panic(fmt.Sprintf("[ERROR] Could not compile ast node of val %v, type %v\n", node, node.NodeType()))
}

func (c *Compiler) Compile(input ast.ProgramNode) []bytecode.Instruction {
	for _, val := range input.Statements {
		c.compileStmt(val)
	}
	return c.ins
}
