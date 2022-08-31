package cmdinterpretator

import (
	"errors"
	"testproject/internal/commands"
	"testproject/internal/inputanalyzer"
	"fmt"
	"os"
)

type CommandCallback func(args []string) error

type CMDInterpretator struct {
	Inputanalyzer *inputanalyzer.InputAnalyzer
	callbacks map[commands.Command]CommandCallback

	cmd    commands.Command
	args   []string
}

func NewCMDInterpetator() *CMDInterpretator {
	return &CMDInterpretator{Inputanalyzer : inputanalyzer.NewInputAnalyzer(os.Stdin, 100),
			callbacks: make(map[commands.Command]CommandCallback),
		}
}

func (ci *CMDInterpretator) Execute() {
	for {
		cmd, args, err := ci.Inputanalyzer.GetCmdWithArgs()
		if err != nil {
			fmt.Println(err)
		}

		err = ci.Run(cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//***

func (ci *CMDInterpretator) SendCommandDescription(cmd commands.Command, f CommandCallback) {
	ci.callbacks[cmd] = f
}

func (ci *CMDInterpretator) Run(cmd commands.Command, args []string) error {
	if f, ok := ci.callbacks[cmd]; ok {
		return f(args)
	}
	return errors.New("CI : no such command registered")
}

func (ci *CMDInterpretator) WriteResponse(answer string) { //
	fmt.Println("get result ", answer)
}

func (ci *CMDInterpretator) Dump() {
	for k, v := range ci.callbacks{
		fmt.Println("cmd", k, "value", v)
	}
}
