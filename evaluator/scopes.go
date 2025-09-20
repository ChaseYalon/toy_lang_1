package evaluator;
import(
	"toy_lang/ast"
	"fmt"
)

type Scope struct {
	Vars   v_map
	Funcs  f_map
	Parent *Scope
}

func (s *Scope) getVar(name string) (ast.Node, bool) {
	if val, ok := s.Vars[name]; ok {
		return val, true
	}
	if s.Parent != nil {
		return s.Parent.getVar(name)
	}
	return nil, false
}

func (s *Scope) getFunc(name string) (ast.FuncDecNode, bool) {
	if val, ok := s.Funcs[name]; ok {
		return val, true
	}
	if s.Parent != nil {
		return s.Parent.getFunc(name)
	}
	return ast.FuncDecNode{}, false
}

func (s *Scope) declareFunc(f ast.FuncDecNode) {
	s.Funcs[f.Name] = f
}

func (s *Scope) declareVar(name string, val ast.Node) {
	s.Vars[name] = val
}

func (s *Scope) assignVar(name string, val ast.Node) bool {
	// Fixed: Added StringLiteral to the condition
	if val.NodeType() == ast.IntLiteral || val.NodeType() == ast.BoolLiteral || val.NodeType() == ast.StringLiteral {
		if _, ok := s.Vars[name]; ok {
			s.Vars[name] = val
			return true
		}
		if s.Parent != nil {
			return s.Parent.assignVar(name, val)
		}
		return false
	}
	panic(fmt.Sprintf("[ERROR] Tried to assign non primitive value to variable, got %v\n", val))
}

func (s *Scope) newChild() *Scope {
	return &Scope{
		Vars:   make(v_map),
		Funcs:  make(f_map),
		Parent: s,
	}
}

func (s *Scope) String() string {
	return fmt.Sprintf("Vars: %+v, Parent: %v\n", s.Vars, s.Parent)
}