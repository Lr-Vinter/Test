package clicontroller

import (
	"testproject/internal/commands"
	"fmt"
	"testproject/internal/inputanalyzer"
	"os"
)

type CLIController struct {
	Inputanalyzer *inputanalyzer.InputAnalyzer

	cmd    commands.Command
	args   []string
}

func NewController() *CLIController {
	c := new(CLIController)
	c.Inputanalyzer = inputanalyzer.NewInputAnalyzer(os.Stdin, 100)
	return c
}

func (c *CMDInterpretator) WriteResponse(answer string) { //
	fmt.Println("get result ", answer)
}

func (c *CMDInterpretator) ReceiveAndRunCMD() error {
	err := c.Inputanalyzer.GetCmdWithArgs(&c.cmd, &c.args)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	return nil
}

func (c *CLIController) ReturnCmdAndArgs() (cmd commands.Command, args []string) {
	return c.cmd, c.args
}
