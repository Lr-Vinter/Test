package cli

import (
	"errors"
	"fmt"
	"testproject/internal/commands"
)

type CommandCallback func(args []string) error

type parser interface {
	GetCmdWithArgs() (commands.Command, []string, error)
}

type CLIcontroller struct {
	Inputanalyzer parser
	callbacks     map[commands.Command]CommandCallback

	cmd  commands.Command
	args []string
}

func NewCLIController(InputAnalyzer parser) *CLIcontroller {
	return &CLIcontroller{
		Inputanalyzer: InputAnalyzer,
		callbacks:     make(map[commands.Command]CommandCallback),
	}
}

func (ci *CLIcontroller) Execute() error {
	cmd, args, err := ci.Inputanalyzer.GetCmdWithArgs()
	if err != nil {
		return err
	}

	err = ci.run(cmd, args)
	if err != nil {
		return err
	}
	return nil
}

//***

func (ci *CLIcontroller) SendCommandDescription(cmd commands.Command, f CommandCallback) {
	ci.callbacks[cmd] = f
}

func (ci *CLIcontroller) run(cmd commands.Command, args []string) error {
	if f, ok := ci.callbacks[cmd]; ok {
		return f(args)
	}
	return errors.New("CI : no such command registered")
}

func (ci *CLIcontroller) WriteResponse(answer string) { //
	fmt.Println("get result ", answer)
}

func (ci *CLIcontroller) Dump() {
	for k, v := range ci.callbacks {
		fmt.Println("cmd", k, "value", v)
	}
}
