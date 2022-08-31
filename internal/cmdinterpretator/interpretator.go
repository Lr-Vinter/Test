package cmdinterpretator

import (
	"errors"
	"testproject/internal/commands"
	"fmt"
)

type CommandCallback func(args []string) error

type CMDInterpretator struct {
	callbacks map[commands.Command]CommandCallback
}

func NewCMDInterpetator() *CMDInterpretator {
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

func (ci *CMDInterpretator) Dump() {
	for k, v := range ci.callbacks{
		fmt.Println("cmd", k, "value", v)
	}
}
