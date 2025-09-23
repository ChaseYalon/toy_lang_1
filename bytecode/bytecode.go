package bytecode

import (
	"fmt"
)

type OpLabel int

const (
	LOAD_INT OpLabel = iota
	INFIX_INT
	DECLARE_VAR
	REF_VAR
)

func (l OpLabel) String() string {
	switch l {
	case LOAD_INT:
		return "LOAD_INT"
	case INFIX_INT:
		return "INFIX_INT"
	case DECLARE_VAR:
		return "DECLARE_VAR"
	case REF_VAR:
		return "REF_VAR"
	default:
		return "UNDEFINED"
	}
}

type Instruction interface {
	OpType() OpLabel
	String() string
}

type LOAD_INT_INS struct {
	Address int
	Value   int
}

func (l *LOAD_INT_INS) OpType() OpLabel {
	return LOAD_INT
}
func (l *LOAD_INT_INS) String() string {
	return fmt.Sprintf("LOAD_INT %d   %d", l.Address, l.Value)
}

type INFIX_INT_INS struct {
	Left_addr    int
	Right_addr   int
	Save_to_addr int
	Operation    int // 1 for add 2 for subtract 3 for multiply 4 for divide
}

func (a *INFIX_INT_INS) OpType() OpLabel {
	return INFIX_INT
}
func (a *INFIX_INT_INS) String() string {
	return fmt.Sprintf("INFIX_INT %d   %d   %d   %d", a.Left_addr, a.Right_addr, a.Save_to_addr, a.Operation)
}

type DECLARE_VAR_INS struct {
	Name string
	Addr int //Points to the address with the variable
}

func (d *DECLARE_VAR_INS) OpType() OpLabel {
	return DECLARE_VAR
}
func (d *DECLARE_VAR_INS) String() string {
	return fmt.Sprintf("DECLARE_VAR   %v   %v   ", d.Name, d.Addr)
}

type REF_VAR_INS struct {
	Name   string
	SaveTo int
}

func (d *REF_VAR_INS) OpType() OpLabel {
	return REF_VAR
}
func (d *REF_VAR_INS) String() string {
	return fmt.Sprintf("REF_VAR   %v   %d   ", d.Name, d.SaveTo)
}
