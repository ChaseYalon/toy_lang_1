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

func (c *Compiler) compileExpr(node ast.Node, mem *int) int {
	if node.NodeType() == ast.EmptyExpr {
		op := node.(*ast.EmptyExprNode)
		return c.compileExpr(op.Child, mem)
	}

	if node.NodeType() == ast.IntLiteral {
		intNode := node.(*ast.IntLiteralNode)

		toRet := bytecode.LOAD_INT_INS{
			Address: *mem,
			Value:   intNode.Value,
		}
		*mem = *mem + 1;
		c.emit(&toRet)
		return c.currBestAdr - 1
	}

	if node.NodeType() == ast.InfixExpr {
		infixNode := node.(*ast.InfixExprNode)

		leftAddr := c.compileExpr(infixNode.Left, mem)
		rightAddr := c.compileExpr(infixNode.Right, mem)

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
		toRet := bytecode.INFIX_INS{
			Left_addr:    leftAddr,
			Right_addr:   rightAddr,
			Save_to_addr: *mem,
			Operation:    opInstr,
		}
		*mem = *mem + 1;
		c.emit(&toRet)
		return *mem - 1
	}

	if node.NodeType() == ast.ReferenceExpr {
		refExpr := node.(*ast.ReferenceExprNode)
		toRet := bytecode.REF_VAR_INS{
			Name:   refExpr.Name,
			SaveTo: *mem,
		}
		*mem = *mem + 1;
		c.emit(&toRet)
		return *mem - 1
	}

	if node.NodeType() == ast.BoolLiteral {
		boolNode := node.(*ast.BoolLiteralNode)
		toRet := bytecode.LOAD_BOOL_INS{
			Address: *mem,
			Value:   boolNode.Value,
		}
		*mem = *mem + 1
		c.emit(&toRet)
		return *mem - 1
	}

	if node.NodeType() == ast.BoolInfix {
		infixNode := node.(*ast.BoolInfixNode)

		leftAddr := c.compileExpr(infixNode.Left, mem)
		rightAddr := c.compileExpr(infixNode.Right, mem)

		opInstr := -1
		switch infixNode.Operator {
		case token.LESS_THAN:
			opInstr = 5
		case token.LESS_THAN_EQT:
			opInstr = 6
		case token.GREATER_THAN:
			opInstr = 7
		case token.GREATER_THAN_EQT:
			opInstr = 8
		case token.EQUALS:
			opInstr = 9
		case token.NOT_EQUAL:
			opInstr = 10
		case token.AND:
			opInstr = 11
		case token.OR:
			opInstr = 12
		default:
			panic(fmt.Sprintf("[ERROR] %v is not a valid infix operator", infixNode.Operator))
		}
		toRet := bytecode.INFIX_INS{
			Left_addr:    leftAddr,
			Right_addr:   rightAddr,
			Save_to_addr: *mem,
			Operation:    opInstr,
		}
		*mem = *mem + 1
		c.emit(&toRet)
		return *mem - 1
	}

	if node.NodeType() == ast.ReturnExpr {
		retNode := node.(*ast.ReturnExprNode)
		addr := c.compileExpr(retNode.Val, mem)
		c.emit(&bytecode.RETURN_INS{Ptr: addr})
		return *mem - 1;
	}

	if node.NodeType() == ast.FuncCall {
		fCallNode := node.(*ast.FuncCallNode)

		var addrs []int = []int{}
		for _, val := range fCallNode.Params {
			addrs = append(addrs, c.compileExpr(val, mem))
		}
		c.emit(&bytecode.FUNC_CALL_INS{Name: fCallNode.Name.Name, Params: addrs, PutRet: *mem})
		*mem = *mem + 1
		return *mem - 1;
	}

	panic(fmt.Sprintf("[ERROR] Got unknown type of %v\n", node.NodeType()))
}

func (c *Compiler) compileStmt(node ast.Node, mem *int) {
	if node.NodeType() == ast.InfixExpr ||
		node.NodeType() == ast.IntLiteral ||
		node.NodeType() == ast.BoolLiteral ||
		node.NodeType() == ast.BoolInfix ||
		node.NodeType() == ast.EmptyExpr ||
		node.NodeType() == ast.ReturnExpr {
		c.compileExpr(node, mem)
		return
	}

	if node.NodeType() == ast.LetStmt {
		letStmt := node.(*ast.LetStmtNode)
		valAddr := c.compileExpr(letStmt.Value, mem)
		toEmit := bytecode.DECLARE_VAR_INS{
			Name: letStmt.Name,
			Addr: valAddr,
		}
		c.emit(&toEmit)
		return
	}

	if node.NodeType() == ast.VarReassign {
		letStmt := node.(*ast.VarReassignNode)
		valAddr := c.compileExpr(letStmt.NewVal, mem)
		toEmit := bytecode.DECLARE_VAR_INS{
			Name: letStmt.Var.Name,
			Addr: valAddr,
		}
		c.emit(&toEmit)
		return
	}

	if node.NodeType() == ast.IfStmt {
		ifNode := node.(*ast.IfStmtNode)

		condAddr := c.compileExpr(ifNode.Cond, mem)
		jmpFalse := &bytecode.JMP_IF_FALSE_INS{
			CondAddr:   condAddr,
			TargetAddr: -1, // placeholder
		}
		c.emit(jmpFalse)
		for _, stmt := range ifNode.Body {
			c.compileStmt(stmt, mem)
		}

		if len(ifNode.Alt) > 0 {
			jmp := &bytecode.JMP_INS{
				InstNum: -1, // placeholder
			}
			c.emit(jmp)

			jmpFalse.TargetAddr = len(c.ins)

			for _, stmt := range ifNode.Alt {
				c.compileStmt(stmt, mem)
			}

			jmp.InstNum = len(c.ins)
		} else {
			jmpFalse.TargetAddr = len(c.ins)
		}
		return
	}

	if node.NodeType() == ast.FuncDec {
		fDecNode := node.(*ast.FuncDecNode)

		c.emit(&bytecode.FUNC_DEC_START_INS{Name: fDecNode.Name, ParamCount: fDecNode.ParamCount})

		//Local frame memory counter

		localBestMem := c.currBestAdr;
		for _, val := range fDecNode.Body {
			c.compileStmt(val, &localBestMem)
		}
		// Emit function end (correct instruction)
		c.emit(&bytecode.FUNC_DEC_END_INS{})
		return
	}

	panic(fmt.Sprintf("[ERROR] Could not compile ast node of val %v, type %v\n", node, node.NodeType()))
}

func (c *Compiler) Compile(input ast.ProgramNode) []bytecode.Instruction {
	for _, val := range input.Statements {
		c.compileStmt(val, &c.currBestAdr)
	}
	for _, val := range c.ins{
		fmt.Printf("%v\n", val);
	}
	return c.ins
}
