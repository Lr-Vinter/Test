package main

import (
	"fmt"
	"os"
	"testproject/internal/cli"
	"testproject/internal/dbcontroller"
	"testproject/internal/inputanalyzer"
	"testproject/internal/service"
)

func main() {

	ia := inputanalyzer.NewInputAnalyzer(os.Stdin, 100)
	db := dbcontroller.NewDBController() //
	i := cli.NewCLIController(ia)        //

	service := service.NewService(db, i)
	service.RegisterCommands()
	for {
		err := i.Execute()
		if err != nil {
			fmt.Println(err)
		}
	}
}
