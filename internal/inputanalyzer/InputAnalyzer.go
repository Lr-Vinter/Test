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

	exprtype    map[string]ExprStruct
	parsedinput []string
}

func NewInputAnalyzer(input *os.File, size int) *InputAnalyzer {
	ia := new(InputAnalyzer)
	ia.reader = bufio.NewReaderSize(input, 100)
	ia.exprtype = make(map[string]ExprStruct)
	ia.specifyExprStruct()

	fmt.Println("New analyzer")
	return ia
}

func (ia *InputAnalyzer) specifyExprStruct() {
	ia.exprtype["set"] = ExprStruct{commands.Set, []string{"a", "100"}}
	ia.exprtype["get"] = ExprStruct{commands.Get, []string{"a"}}
	ia.exprtype["exit"] = ExprStruct{commands.Exit, []string{}}
}

func (ia *InputAnalyzer) parseInput() error {
	input, err := ia.reader.ReadString(10)
	if err != nil {
		return errors.New("IA: Failed to read input string")
	}

	buff := strings.TrimRight(input, "\r\n")
	ia.parsedinput = strings.Split(buff, " ")

	for i := 0; i < len(ia.parsedinput); i++ {
		if len(ia.parsedinput[i]) < 1 {
			ia.parsedinput = append(ia.parsedinput[:i],
				ia.parsedinput[i+1:]...)
			i--
		}
	}

	if len(ia.parsedinput) == 0 {
		return errors.New("IA: Zero line")
	}
	return nil
}

func (ia *InputAnalyzer) GetCmdWithArgs() (commands.Command, []string, error) { ///////
	err := ia.parseInput()
	if err != nil {
		return commands.Exit, []string{}, err
	}

	if value, ok := ia.exprtype[ia.parsedinput[0]]; ok {
		if len(ia.parsedinput) == len(value.args)+1 {
			cmd, args := value.cmd, ia.parsedinput[1:]
			return cmd, args, nil
		}
		return commands.Exit, []string{}, errors.New("Wrong arg number...")
	}

	return commands.Exit, []string{}, errors.New("Wrong input, command not found...")
}

//func (inputanalyzer *InputAnalyzer) dump() {
//	fmt.Println("input", inputanalyzer.input)
//	fmt.Printf("parsed input %q\n", inputanalyzer.parsedinput)
//}
