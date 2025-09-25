package vm

import (
	"fmt"
)

type Scope struct {
	Vars   map[string]int
	Funcs  map[string]int //Ptr to start of function
	Parent *Scope
	isMain bool
}

func (s *Scope) newChild() *Scope {
	return &Scope{
		Vars:   make(map[string]int),
		Parent: s,
		isMain: false,
		Funcs:  make(map[string]int),
	}
}

func (s *Scope) getVar(input string) int {
	lVar, found := s.Vars[input]
	if !found && s.isMain {
		panic(fmt.Sprintf("[ERROR] Variable \"%v\" is undefined, vars are %v\n", input, s.Vars))
	}
	if !found {
		return s.Parent.getVar(input)
	}

	return lVar
}

func (s *Scope) setVar(name string, address int) {
	s.Vars[name] = address
}
func (s *Scope) getFunc(name string) int {
	lF, found := s.Funcs[name]
	if !found && s.isMain {
		panic(fmt.Sprintf("[ERROR] Function \"%v\" is undefined\n", name))
	}
	if !found {
		return s.Parent.getFunc(name)
	}
	return lF
}

func (s *Scope) setFunc(name string, startAddr int) {
	s.Funcs[name] = startAddr
}

func NewMainScope() *Scope {
	return &Scope{
		Vars:   make(map[string]int),
		Parent: nil,
		isMain: true,
		Funcs:  make(map[string]int),
	}
}
