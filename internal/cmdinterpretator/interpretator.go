package cmdinterpretator

import (
	"errors"
	"testproject/internal/commands"
)

type CommandCallback func(args []string) error

type CMDInterpretator struct {
	callbacks map[commands.Command]CommandCallback
}

func NewCMDInterpetator() *CMDInterpretator {
	//ci = new(CMDInterpretator{callbacks: make(map[command]CommandCallback), db: db}) ? поч так нельзя
	ci := &CMDInterpretator{callbacks: make(map[commands.Command]CommandCallback)}
	return ci
}

func (ci *CMDInterpretator) SendCommandDescription(cmd commands.Command, f CommandCallback) {
	ci.callbacks[cmd] = f
}

func (ci *CMDInterpretator) Run(cmd commands.Command, args []string) error {
	if f, ok := ci.callbacks[cmd]; ok {
		return f(args)
	}
	return errors.New("CI : no such command registered")
}
