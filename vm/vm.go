package vm

import (
	"fmt"
	"toy_lang/bytecode"
)

type Vm struct {
	Ins  []bytecode.Instruction
	insPtr int
	Ram  [1 << 11]any
	Vars map[string]int
}

func NewVm() *Vm {
	return &Vm{
		Ins:  []bytecode.Instruction{},
		Ram:  [1 << 11]any{},
		Vars: make(map[string]int),
		insPtr: 0,
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
		panic(fmt.Sprintf("[ERROR] Could not convert type %v into a bool\n", t));
	}
}

func (v *Vm) executeIns(ins bytecode.Instruction) {
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
			bLeft := v.convertToBool(v.Ram[op.Left_addr]);
			bRight := v.convertToBool(v.Ram[op.Right_addr]);
			v.Ram[op.Save_to_addr] = bLeft && bRight;
		case 12:
			bLeft := v.convertToBool(v.Ram[op.Left_addr]);
			bRight := v.convertToBool(v.Ram[op.Right_addr]);
			v.Ram[op.Save_to_addr] = bLeft || bRight;
		default:
			panic(fmt.Sprintf("[ERROR] Unknown operator, got %v\n", op.Operation))
		}
	case bytecode.DECLARE_VAR:
		op := ins.(*bytecode.DECLARE_VAR_INS)
		v.Vars[op.Name] = op.Addr
	case bytecode.REF_VAR:
		op := ins.(*bytecode.REF_VAR_INS)
		v.Ram[op.SaveTo] = v.Ram[v.Vars[op.Name]]
	case bytecode.LOAD_BOOL:
		op := ins.(*bytecode.LOAD_BOOL_INS)
		v.Ram[op.Address] = op.Value
	case bytecode.JMP:
		op := ins.(*bytecode.JMP_INS);
		v.insPtr = op.InstNum;
	case bytecode.JMP_IF_FALSE:
		op := ins.(*bytecode.JMP_IF_FALSE_INS);
		cond := v.convertToBool(v.Ram[op.CondAddr]);
		if !cond{
			v.insPtr = op.TargetAddr;
		}
	}
}

func (v *Vm) Execute(instructions []bytecode.Instruction, shouldPrint bool) (*[1 << 11]any, map[string]any) {
	v.Ins = instructions;

	for v.insPtr < len(v.Ins){
		v.executeIns(v.Ins[v.insPtr]);
		v.insPtr++;
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
	for key, val := range v.Vars {
		toReturnVars[key] = v.Ram[val]
	}
	return &v.Ram, toReturnVars
}
