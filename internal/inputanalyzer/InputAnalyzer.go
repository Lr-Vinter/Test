package inputanalyzer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"testproject/internal/commands"
)

type ExprStruct struct {
	cmd  commands.Command
	args []string
}

type InputAnalyzer struct {
	reader *bufio.Reader
	input  string

	exprtype    map[string]ExprStruct
	parsedinput []string
}

func NewInputAnalyzer(input *os.File, size int) *InputAnalyzer {
	ia := new(InputAnalyzer)
	ia.reader = bufio.NewReaderSize(input, 100)
	ia.exprtype = make(map[string]ExprStruct)
	ia.SpecifyExprStruct()

	fmt.Println("New analyzer")
	return ia
}

func (ia *InputAnalyzer) SpecifyExprStruct() {
	ia.exprtype["set"] = ExprStruct{commands.Set, []string{"a", "100"}}
	ia.exprtype["get"] = ExprStruct{commands.Get, []string{"a"}}
	ia.exprtype["exit"] = ExprStruct{commands.Exit, []string{}}
}

func (ia *InputAnalyzer) GetInput() error {
	err := errors.New("Failed to read string")
	ia.input, err = ia.reader.ReadString(10)
	return err
}

func (ia *InputAnalyzer) ParseInput() error {
	err := ia.GetInput()
	buff := strings.TrimRight(ia.input, "\r\n")
	ia.parsedinput = strings.Split(buff, " ")

	for i := 0; i < len(ia.parsedinput); i++ {
		if len(ia.parsedinput[i]) < 1 {
			ia.parsedinput = append(ia.parsedinput[:i],
				ia.parsedinput[i+1:]...)
			i--
		}
	}

	return err
}

func (ia *InputAnalyzer) GetCmdWithArgs(cmd *commands.Command, args *[]string) error {
	if value, ok := ia.exprtype[ia.parsedinput[0]]; ok {
		if len(ia.parsedinput) == len(value.args)+1 { // cmd length = 1
			*cmd = value.cmd
			*args = ia.parsedinput[1:]
			return nil
		}
		err := errors.New("Wrong arg number...")
		return err
	}
	err := errors.New("Wrong input, command not found...")
	return err
}

func (inputanalyzer *InputAnalyzer) Dump() {
	fmt.Println("input", inputanalyzer.input)
	fmt.Printf("parsed input %q\n", inputanalyzer.parsedinput)
}
