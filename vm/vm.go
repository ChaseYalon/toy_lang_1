package vm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	s := NewMainScope()
	return &Vm{
		Ins:       []bytecode.Instruction{},
		Ram:       [1 << 11]any{},
		insPtr:    0,
		MainScope: s,
		currScope: s,
		CallStack: []*CallFrame{
			{
				Local_scope: s, //Main frame
				ArgAddrs:    []int{},
				ResumeAddr:  -1,
				PutVal:      -1,
			},
		},
	}
}

func (v *Vm) convertToInt(input any) (int, error) {
	switch t := input.(type) {
	case int:
		return t, nil
	case bool:
		bIn := input.(bool)
		if bIn {
			return 1, nil
		} else {
			return 0, nil
		}
	case string:
		s := input.(string)
		i, err := strconv.Atoi(s)
		if err != nil {
			return -1, err
		}
		return i, nil
	default:
		return -1, fmt.Errorf("[ERROR] Could not convert %v of type %v to an int", input, t)
	}

}
func (v *Vm) convertToBool(input any) (bool, error) {
	switch t := input.(type) {
	case int:
		i := input.(int)
		if i > 0 {
			return true, nil
		}
		return false, nil
	case bool:
		return input.(bool), nil
	case string:
		s := input.(string)
		if s == "true" {
			return true, nil
		}
		if s == "false" {
			return false, nil
		}
		if s == "" {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("[ERROR] Could not convert type %v into a bool", t)
	}
}
func (v *Vm) convertToString(input any) (string, error) {
	switch t := input.(type) {
	case int:
		return strconv.Itoa(t), nil
	case bool:
		b := input.(bool)
		res := strconv.FormatBool(b)
		return res, nil
	case string:
		s := input.(string)
		return s, nil
	default:
		return "", fmt.Errorf("[ERROR] Could not convert type of %v into a string", t)
	}
}

func (v *Vm) executeIns(ins bytecode.Instruction, local_scope *Scope) {
	switch ins.OpType() {
	case bytecode.LOAD_INT:
		op := ins.(*bytecode.LOAD_INT_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.LOAD_STRING:
		op := ins.(*bytecode.LOAD_STRING_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.INFIX_INT:
		op := ins.(*bytecode.INFIX_INS)

		leftVal, err1 := v.convertToInt(v.Ram[op.Left_addr])
		rightVal, err2 := v.convertToInt(v.Ram[op.Right_addr])

		if err1 != nil || err2 != nil {
			if op.Operation == 1 || op.Operation == 9 || op.Operation == 10 {
				leftStr, serr1 := v.convertToString(v.Ram[op.Left_addr])
				rightStr, serr2 := v.convertToString(v.Ram[op.Right_addr])
				if serr1 != nil || serr2 != nil {
					panic(fmt.Sprintf("[ERROR] Op: %v, Err 1: %v, Err2: %v", op, err1, err2))
				}
				if op.Operation == 1 {
					v.Ram[op.Save_to_addr] = leftStr + rightStr
				}
				if op.Operation == 9 {
					v.Ram[op.Save_to_addr] = leftStr == rightStr
				}
				if op.Operation == 10 {
					v.Ram[op.Save_to_addr] = leftStr != rightStr
				}
				return

			}
		}
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
		case 13:
			v.Ram[op.Save_to_addr] = leftVal % rightVal
		case 11:
			bLeft, _ := v.convertToBool(v.Ram[op.Left_addr])
			bRight, _ := v.convertToBool(v.Ram[op.Right_addr])
			v.Ram[op.Save_to_addr] = bLeft && bRight
		case 12:
			bLeft, _ := v.convertToBool(v.Ram[op.Left_addr])
			bRight, _ := v.convertToBool(v.Ram[op.Right_addr])
			v.Ram[op.Save_to_addr] = bLeft || bRight
		default:
			panic(fmt.Sprintf("[ERROR] Unknown operator, got %v\n", op.Operation))
		}
	case bytecode.DECLARE_VAR:
		op := ins.(*bytecode.DECLARE_VAR_INS)
		local_scope.setVar(op.Name, op.Addr)
	case bytecode.REF_VAR:
		op := ins.(*bytecode.REF_VAR_INS)
		varAddr := local_scope.getVar(op.Name)
		if v.Ram[varAddr] == nil {
			panic(fmt.Sprintf("[ERROR] Variable %s at address %d contains nil value", op.Name, varAddr))
		}
		v.Ram[op.SaveTo] = v.Ram[varAddr]
	case bytecode.LOAD_BOOL:
		op := ins.(*bytecode.LOAD_BOOL_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.JMP:
		op := ins.(*bytecode.JMP_INS)
		v.insPtr = op.InstNum
	case bytecode.JMP_IF_FALSE:
		op := ins.(*bytecode.JMP_IF_FALSE_INS)
		cond, _ := v.convertToBool(v.Ram[op.CondAddr])
		if !cond {
			v.insPtr = op.TargetAddr
		}
	case bytecode.FUNC_DEC_START:
		op := ins.(*bytecode.FUNC_DEC_START_INS)
		name := op.Name
		startPtr := v.insPtr
		local_scope.setFunc(name, startPtr)
		delta := 1
		for {
			if v.insPtr+delta >= len(v.Ins) {
				panic("[ERROR] Reached end of instructions while scanning for FUNC_DEC_END")
			}

			if v.Ins[v.insPtr+delta].OpType() == bytecode.FUNC_DEC_END {
				v.insPtr = v.insPtr + delta
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
			fScope.setVar(fDec.ParamNames[i], valAddr)
		}
		v.currScope = fScope
		v.insPtr = startAddr
	case bytecode.CALL_BUILTIN:
		op := ins.(*bytecode.CALL_BUILTIN_INS)

		switch op.Name {
		case "println":
			input := v.Ram[op.Params[0]]
			sIn, ok := input.(string)
			if !ok {
				byteIn := input.(bytecode.Instruction)
				panic(fmt.Sprintf("[ERROR] println must be passed a string, got %v\n", byteIn.OpType()))
			}
			fmt.Println(sIn)
		case "print":
			input := v.Ram[op.Params[0]]
			sIn, ok := input.(string)
			if !ok {
				byteIn := input.(bytecode.Instruction)
				panic(fmt.Sprintf("[ERROR] print must be passed a string, got %v\n", byteIn.OpType()))
			}
			fmt.Print(sIn)
		case "input":
			prompt := v.Ram[op.Params[0]]
			sIn, ok := prompt.(string)
			if !ok {
				byteIn := prompt.(bytecode.Instruction)
				panic(fmt.Sprintf("[ERROR] print must be passed a string, got %v\n", byteIn.OpType()))
			}
			fmt.Print(sIn)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimRight(input, "\r\n")
			v.Ram[op.PutRet] = input
		case "str":
			in := v.Ram[op.Params[0]]
			s, err := v.convertToString(in)
			if err != nil {
				panic(fmt.Sprintf("%v\n", err))
			}
			v.Ram[op.PutRet] = s
		case "bool":
			in := v.Ram[op.Params[0]]
			b, err := v.convertToBool(in)
			if err != nil {
				panic(fmt.Sprintf("%v\n", err))
			}
			v.Ram[op.PutRet] = b
		case "int":
			in := v.Ram[op.Params[0]]
			i, err := v.convertToInt(in)
			if err != nil {
				panic(fmt.Sprintf("%v\n", err))
			}
			v.Ram[op.PutRet] = i
		}

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
			continue
		case bytecode.FUNC_DEC_END:
			if len(v.CallStack) > 0 {
				callItem := v.CallStack[len(v.CallStack)-1]
				v.insPtr = callItem.ResumeAddr
				v.CallStack = v.CallStack[:len(v.CallStack)-1]
				v.currScope = v.CallStack[len(v.CallStack)-1].Local_scope

				continue
			}

		}
		v.executeIns(v.Ins[v.insPtr], v.currScope)
		v.insPtr++
	}

	if shouldPrint {
		var vals map[string]any = make(map[string]any)
		for key, val := range v.MainScope.Vars {
			vals[key] = v.Ram[val]
		}
		fmt.Printf("%+v\n", vals)
	}
	toReturnVars := make(map[string]any)
	for key, val := range v.MainScope.Vars {
		toReturnVars[key] = v.Ram[val]
	}
	return &v.Ram, toReturnVars
}
