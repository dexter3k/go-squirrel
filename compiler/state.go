package compiler

import (
	"github.com/dexter3k/go-squirrel/sqvm"
)

type state struct {
	strings  map[string]int
	strTable []string
}

func newState() *state {
	return &state{
		strings:  map[string]int{},
		strTable: []string{},
	}
}

func (s *state) makeString(str string) int {
	if i, present := s.strings[str]; present {
		return i
	}
	s.strings[str] = len(s.strings)
	strTable = append(strTable, str)
	return len(s.strings) - 1
}

func (s state) makeFuncProto() (*sqvm.FuncProto, error) {
	return nil, nil
}

func (s *state) popTarget() {
}
