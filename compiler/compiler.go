package compiler

import (
	"fmt"
	"toy_lang/ast"
	"toy_lang/bytecode"
	"toy_lang/token"
)

type WhileStmtRecord struct {
	CondAddr int
	brks     []*bytecode.JMP_INS
}

type Compiler struct {
	ins         []bytecode.Instruction
	currBestAdr int
	//Hash set to check if a function name is predefined
	builtinFuncs map[string]string
	whileStack   []*WhileStmtRecord
}

func NewCompiler() *Compiler {
	return &Compiler{
		ins:         []bytecode.Instruction{},
		currBestAdr: 0,
		builtinFuncs: map[string]string{
			"println": "println",
			"print":   "print",
			"input":   "input",
			"str":     "str",
			"int":     "int",
			"bool":    "bool",
		},
		whileStack: []*WhileStmtRecord{},
	}
}

func (c *Compiler) pushWhile(stmt *WhileStmtRecord) {
	c.whileStack = append(c.whileStack, stmt)
}

func (c *Compiler) popWhile() *WhileStmtRecord {
	if len(c.whileStack) == 0 {
		panic("popWhile: stack is empty")
	}
	last := c.whileStack[len(c.whileStack)-1]
	c.whileStack = c.whileStack[:len(c.whileStack)-1]
	return last
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
		*mem = *mem + 1
		c.emit(&toRet)
		return *mem - 1
	}
	if node.NodeType() == ast.FloatLiteral {
		fNode := node.(*ast.FloatLiteralNode)
		toRet := bytecode.LOAD_FLOAT_INS{
			Address: *mem,
			Value:   fNode.Value,
		}
		*mem = *mem + 1
		c.emit(&toRet)
		return *mem - 1
	}
	if node.NodeType() == ast.StringLiteral {
		op := node.(*ast.StringLiteralNode)
		toRet := bytecode.LOAD_STRING_INS{
			Address: *mem,
			Value:   op.Value,
		}
		*mem = *mem + 1
		c.emit(&toRet)
		return *mem - 1
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
		case token.MODULO:
			opInstr = 13
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

	if node.NodeType() == ast.ReferenceExpr {
		refExpr := node.(*ast.ReferenceExprNode)
		toRet := bytecode.REF_VAR_INS{
			Name:   refExpr.Name,
			SaveTo: *mem,
		}
		*mem = *mem + 1
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
		return *mem - 1
	}

	if node.NodeType() == ast.FuncCall {
		fCallNode := node.(*ast.FuncCallNode)

		name := fCallNode.Name.Name
		if _, exists := c.builtinFuncs[name]; exists {
			var addrs []int = []int{}
			for _, val := range fCallNode.Params {
				addrs = append(addrs, c.compileExpr(val, mem))
			}
			c.emit(&bytecode.CALL_BUILTIN_INS{Name: fCallNode.Name.Name, Params: addrs, PutRet: *mem})
			*mem = *mem + 1
			return *mem - 1
		}

		var addrs []int = []int{}
		for _, val := range fCallNode.Params {
			addrs = append(addrs, c.compileExpr(val, mem))
		}
		c.emit(&bytecode.FUNC_CALL_INS{Name: fCallNode.Name.Name, Params: addrs, PutRet: *mem})
		*mem = *mem + 1
		return *mem - 1
	}

	panic(fmt.Sprintf("[ERROR] Got unknown type of %v\n", node.NodeType()))
}

func (c *Compiler) compileStmt(node ast.Node, mem *int) {
	if node.NodeType() == ast.InfixExpr ||
		node.NodeType() == ast.IntLiteral ||
		node.NodeType() == ast.BoolLiteral ||
		node.NodeType() == ast.BoolInfix ||
		node.NodeType() == ast.EmptyExpr ||
		node.NodeType() == ast.ReturnExpr ||
		node.NodeType() == ast.FloatLiteral {
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
			jmpFalse.TargetAddr = len(c.ins) - 1
		}
		return
	}
	if node.NodeType() == ast.WhileStmt {
		whileNode := node.(*ast.WhileStmtNode)
		condInsNum := len(c.ins)
		condAddr := c.compileExpr(whileNode.Cond, mem)

		jmpFalse := &bytecode.JMP_IF_FALSE_INS{
			CondAddr:   condAddr,
			TargetAddr: -1, // placeholder
		}
		c.emit(jmpFalse)
		var breaks []*bytecode.JMP_INS
		c.pushWhile(&WhileStmtRecord{
			CondAddr: condAddr,
			brks:     breaks,
		})
		for _, stmt := range whileNode.Body {
			c.compileStmt(stmt, mem)
		}
		//This section here is awful and will cause bugs
		if condInsNum == 0 {
			condInsNum = 1
		}
		c.emit(&bytecode.JMP_INS{InstNum: condInsNum - 1})
		lastWhile := c.popWhile()
		jmpFalse.TargetAddr = len(c.ins)
		for _, val := range lastWhile.brks {
			val.InstNum = len(c.ins)
		}
		return
	}
	if node.NodeType() == ast.ContinueStmt {
		lastWhile := c.popWhile()
		c.emit(&bytecode.JMP_INS{InstNum: lastWhile.CondAddr})
		c.pushWhile(lastWhile)
		return
	}
	if node.NodeType() == ast.BreakSmt {
		lastWhile := c.popWhile()
		ins := &bytecode.JMP_INS{InstNum: -1} //Placeholder
		lastWhile.brks = append(lastWhile.brks, ins)
		c.emit(ins)

		c.pushWhile(lastWhile)
		return
	}

	if node.NodeType() == ast.FuncDec {
		fDecNode := node.(*ast.FuncDecNode)
		var paramNames []string
		for _, val := range fDecNode.Params {
			paramNames = append(paramNames, val.Name)
		}
		c.emit(&bytecode.FUNC_DEC_START_INS{Name: fDecNode.Name, ParamCount: fDecNode.ParamCount, ParamNames: paramNames})

		//Local frame memory counter

		//localBestMem := c.currBestAdr
		for _, val := range fDecNode.Body {
			c.compileStmt(val, &c.currBestAdr)
		}
		c.emit(&bytecode.FUNC_DEC_END_INS{})
		return
	}

	panic(fmt.Sprintf("[ERROR] Could not compile ast node of val %v, type %v\n", node, node.NodeType()))
}

func (c *Compiler) Compile(input ast.ProgramNode) []bytecode.Instruction {
	for _, val := range input.Statements {
		c.compileStmt(val, &c.currBestAdr)
	}
	return c.ins
}
