package cli

import (
	"testing"
	"testproject/internal/commands"
)

type MockParser struct {
	cmd  commands.Command
	args []string
	err  error
}

func (m *MockParser) GetCmdWithArgs() (commands.Command, []string, error) {
	return m.cmd, m.args, m.err
}

var T_conn commands.Command = "conn"
var T_comp commands.Command = "comp"

//***

type FuncAndCounter struct {
	counter int
}

func (m *FuncAndCounter) compare(args []string) error {
	m.counter++
	return nil
}

func TestExecute(t *testing.T) {

	client := NewCLIController(&MockParser{T_conn, []string{"df", "ddf"}, nil})
	fac := &FuncAndCounter{0}

	client.SendCommandDescription(T_conn, fac.compare)

	output := client.Execute()
	if fac.counter != 1 || output != nil {
		t.Fatal(fac.counter, "-", output)
	}
}
