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
	LOAD_BOOL
	JMP
	JMP_IF_FALSE
	FUNC_DEC_START
	FUNC_DEC_END
	FUNC_CALL
	RETURN
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
	case LOAD_BOOL:
		return "LOAD_BOOL"
	case JMP:
		return "JMP"
	case JMP_IF_FALSE:
		return "JMP_IF_FALSE"
	case FUNC_DEC_START:
		return "FUNC_DEC_START"
	case FUNC_DEC_END:
		return "FUNC_DEC_END"
	case FUNC_CALL:
		return "FUNC_CALL"
	case RETURN:
		return "RETURN"
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
	return fmt.Sprintf("LOAD_INT ADDR(%d), VAL(%d)", l.Address, l.Value)
}

type INFIX_INS struct {
	Left_addr    int
	Right_addr   int
	Save_to_addr int
	Operation    int // 1 for add 2 for subtract 3 for multiply 4 for divide 5 for less then 6 for less then eqt 7 for greater then 8 for greater then eqt 9 for equals, 10 for not equals, 11 for and, 12 for or
}

func (a *INFIX_INS) OpType() OpLabel {
	return INFIX_INT
}
func (a *INFIX_INS) String() string {
	return fmt.Sprintf("INFIX_INT LEFT_ADDR(%d), RIGHT_ADDR(%d), SAVE_TO_ADDR(%d), OPERATOR(%d)", a.Left_addr, a.Right_addr, a.Save_to_addr, a.Operation)
}

type DECLARE_VAR_INS struct {
	Name string
	Addr int //Points to the address with the variable
}

func (d *DECLARE_VAR_INS) OpType() OpLabel {
	return DECLARE_VAR
}
func (d *DECLARE_VAR_INS) String() string {
	return fmt.Sprintf("DECLARE_VAR NAME(%v), VAL_ADDR(%v)", d.Name, d.Addr)
}

type REF_VAR_INS struct {
	Name   string
	SaveTo int
}

func (d *REF_VAR_INS) OpType() OpLabel {
	return REF_VAR
}
func (d *REF_VAR_INS) String() string {
	return fmt.Sprintf("REF_VAR  NAME(%v) SAVE_TO(%d)   ", d.Name, d.SaveTo)
}

type LOAD_BOOL_INS struct {
	Address int
	Value   bool
}

func (l *LOAD_BOOL_INS) OpType() OpLabel {
	return LOAD_BOOL
}
func (l *LOAD_BOOL_INS) String() string {
	return fmt.Sprintf("LOAD_BOOL ADDR(%v) VALUE(%v)", l.Address, l.Value)
}



type JMP_INS struct{
	InstNum int // Address, not change 
}
func (j *JMP_INS)OpType() OpLabel{
	return JMP
}
func (j *JMP_INS)String() string{
	return fmt.Sprintf("JUMP_TO ADDR(%d)", j.InstNum);
}



type JMP_IF_FALSE_INS struct{
	CondAddr int //Address with the bool
	TargetAddr int
}
func (j *JMP_IF_FALSE_INS)OpType() OpLabel{
	return JMP_IF_FALSE
}
func (j *JMP_IF_FALSE_INS)String() string{
	return fmt.Sprintf("JUMP_TO_IF CondADDR(%d) JmpToAddr(%d)", j.CondAddr, j.TargetAddr)
}



type FUNC_DEC_START_INS struct{
	Name string
	ParamCount int
}
func (f *FUNC_DEC_START_INS)OpType() OpLabel{
	return FUNC_DEC_START
}
func (f *FUNC_DEC_START_INS)String() string{
	return fmt.Sprintf("FUNC_DEC_START NAME(%v) PARAMS(%d)", f.Name, f.ParamCount);
}



type FUNC_DEC_END_INS struct{} //For all intents and purposes return, if the user does not explicitly
func (f *FUNC_DEC_END_INS) OpType()OpLabel{
	return FUNC_DEC_END
}
func (f *FUNC_DEC_END_INS)String()string{
	return "FUNC_DEC_END"
}



type FUNC_CALL_INS struct{
	Params []int //Pointers to where the parameters are held, in order so for fn add(a, b) it would need to point to a first and b second
	Name string
	PutRet int
}
func (f *FUNC_CALL_INS) OpType() OpLabel{
	return FUNC_CALL
}
func (f *FUNC_CALL_INS)String() string{
	return fmt.Sprintf("FUNC_CALL PARAMS(%+d) NAME(%v) SAVE_VAL(%d)", f.Params, f.Name, f.PutRet);
}


type RETURN_INS struct{
	Ptr int
}
func (r *RETURN_INS) OpType() OpLabel{
	return RETURN;
}
func (r *RETURN_INS) String() string{
	return fmt.Sprintf("RETURN %d", r.Ptr);
}