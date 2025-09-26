package vm

import (
	"fmt"
	"toy_lang/bytecode"
)

type CallFrame struct {
	Local_scope *Scope
	ArgAddrs    []int
	ResumeAddr  int //Where to put the stack pointer after execution
	PutVal      int //Where the return value (if any) should be placed in memory
}

type Vm struct {
	Ins       []bytecode.Instruction
	insPtr    int
	Ram       [1 << 11]any
	MainScope *Scope
	CallStack []*CallFrame
	currScope *Scope
}

func NewVm() *Vm {
	s :=NewMainScope();
	return &Vm{
		Ins:       []bytecode.Instruction{},
		Ram:       [1 << 11]any{},
		insPtr:    0,
		MainScope: s,
		currScope: s,
		CallStack: []*CallFrame{
			&CallFrame{
				Local_scope: s, //Main frame
				ArgAddrs: []int{},
				ResumeAddr: -1,
				PutVal: -1,
			},
		},
	}
}

func (v *Vm) convertToInt(input any) int {
	switch t := input.(type) {
	case int:
		return t
	case bool:
		bIn := input.(bool)
		if bIn {
			return 1
		} else {
			return 0
		}
	default:
		panic(fmt.Sprintf("[ERROR] Could not convert %v of type %v to an int", input, t))
	}

}
func (v *Vm) convertToBool(input any) bool {
	switch t := input.(type) {
	case int:
		i := input.(int)
		if i > 0 {
			return true
		}
		return false
	case bool:
		return input.(bool)
	default:
		panic(fmt.Sprintf("[ERROR] Could not convert type %v into a bool\n", t))
	}
}

func (v *Vm) executeIns(ins bytecode.Instruction, local_scope *Scope) {
	switch ins.OpType() {
	case bytecode.LOAD_INT:
		op := ins.(*bytecode.LOAD_INT_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.INFIX_INT:
		op := ins.(*bytecode.INFIX_INS)

		leftVal := v.convertToInt(v.Ram[op.Left_addr])
		rightVal := v.convertToInt(v.Ram[op.Right_addr])
		switch op.Operation {
		case 1:
			v.Ram[op.Save_to_addr] = leftVal + rightVal
			return
		case 2:
			v.Ram[op.Save_to_addr] = leftVal - rightVal
			return
		case 3:
			v.Ram[op.Save_to_addr] = leftVal * rightVal
			return
		case 4:
			v.Ram[op.Save_to_addr] = leftVal / rightVal
			return
		case 5:
			v.Ram[op.Save_to_addr] = leftVal < rightVal
		case 6:
			v.Ram[op.Save_to_addr] = leftVal <= rightVal
		case 7:
			v.Ram[op.Save_to_addr] = leftVal > rightVal
		case 8:
			v.Ram[op.Save_to_addr] = leftVal >= rightVal
		case 9:
			v.Ram[op.Save_to_addr] = leftVal == rightVal
		case 10:
			v.Ram[op.Save_to_addr] = leftVal != rightVal
		case 11:
			bLeft := v.convertToBool(v.Ram[op.Left_addr])
			bRight := v.convertToBool(v.Ram[op.Right_addr])
			v.Ram[op.Save_to_addr] = bLeft && bRight
		case 12:
			bLeft := v.convertToBool(v.Ram[op.Left_addr])
			bRight := v.convertToBool(v.Ram[op.Right_addr])
			v.Ram[op.Save_to_addr] = bLeft || bRight
		default:
			panic(fmt.Sprintf("[ERROR] Unknown operator, got %v\n", op.Operation))
		}
	case bytecode.DECLARE_VAR:
		op := ins.(*bytecode.DECLARE_VAR_INS)
		local_scope.setVar(op.Name, op.Addr)
	case bytecode.REF_VAR:
		op := ins.(*bytecode.REF_VAR_INS)
		v.Ram[op.SaveTo] = v.Ram[local_scope.getVar(op.Name)]
	case bytecode.LOAD_BOOL:
		op := ins.(*bytecode.LOAD_BOOL_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.JMP:
		op := ins.(*bytecode.JMP_INS)
		v.insPtr = op.InstNum
	case bytecode.JMP_IF_FALSE:
		op := ins.(*bytecode.JMP_IF_FALSE_INS)
		cond := v.convertToBool(v.Ram[op.CondAddr])
		if !cond {
			v.insPtr = op.TargetAddr
		}
	case bytecode.FUNC_DEC_START:
		op := ins.(*bytecode.FUNC_DEC_START_INS)
		name := op.Name
		startPtr := v.insPtr
		local_scope.setFunc(name, startPtr)
		fmt.Println("in func_dec_start")
		delta := 1
		for {
			if v.insPtr+delta >= len(v.Ins) {
				panic("[ERROR] Reached end of instructions while scanning for FUNC_DEC_END")
			}

			if v.Ins[v.insPtr+delta].OpType() == bytecode.FUNC_DEC_END {
				v.insPtr = v.insPtr + delta + 1
				fmt.Printf("next instr: %v\n", v.Ins[v.insPtr]);
				return
			}

			delta++
			if delta > len(v.Ins) {
				panic("[ERROR] Infinite loop while scanning for FUNC_DEC_END")
			}
		}

	case bytecode.FUNC_CALL:
		op := ins.(*bytecode.FUNC_CALL_INS)
		startAddr := v.MainScope.getFunc(op.Name)

		// Use the current scope as parent, so parameter addresses exist
		fScope := local_scope.newChild()
		callFrame := &CallFrame{
			Local_scope: fScope,
			ArgAddrs:    op.Params,
			ResumeAddr:  v.insPtr + 1,
			PutVal:      op.PutRet,
		}
		v.CallStack = append(v.CallStack, callFrame)

		fDec := v.Ins[startAddr].(*bytecode.FUNC_DEC_START_INS)

		// Map parameters to RAM addresses provided by the compiler
		for i, valAddr := range op.Params {
			fmt.Printf("appending variable of name %v to value %v\n", fDec.ParamNames[i], valAddr);
			fScope.setVar(fDec.ParamNames[i], valAddr)
		}
		v.currScope = fScope;
		v.insPtr = startAddr + 1

	}
}

func (v *Vm) Execute(instructions []bytecode.Instruction, shouldPrint bool) (*[1 << 11]any, map[string]any) {
	v.Ins = instructions

	for v.insPtr < len(v.Ins) {
		ins := v.Ins[v.insPtr]

		switch ins.OpType() {
		case bytecode.RETURN:
			retIns := ins.(*bytecode.RETURN_INS)
			callItem := v.CallStack[len(v.CallStack)-1]
			v.Ram[callItem.PutVal] = v.Ram[retIns.Ptr]
			v.insPtr = callItem.ResumeAddr
			v.CallStack = v.CallStack[:len(v.CallStack)-1]
			v.currScope = v.CallStack[len(v.CallStack)-1].Local_scope
			continue // skip v.insPtr++ because we already updated it
		case bytecode.FUNC_DEC_END:
			if len(v.CallStack) > 0 {
				callItem := v.CallStack[len(v.CallStack)-1]
				v.insPtr = callItem.ResumeAddr
				v.CallStack = v.CallStack[:len(v.CallStack)-1]
				v.currScope = v.CallStack[len(v.CallStack)-1].Local_scope

				continue
			}

		}
		fmt.Printf("Executing instruciton: %v\n", v.Ins[v.insPtr])
		v.executeIns(v.Ins[v.insPtr], v.currScope)
		v.insPtr++
	}

	if shouldPrint {
		fmt.Printf("RAM: ")
		for i := range 2048 {
			if v.Ram[i] != nil {
				fmt.Printf("{%v}, ", v.Ram[i])
			}
		}
		fmt.Print("\n")
	}
	toReturnVars := make(map[string]any)
	for key, val := range v.MainScope.Vars {
		toReturnVars[key] = v.Ram[val]
	}
	return &v.Ram, toReturnVars
}
