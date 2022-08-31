package main

import (
	"testproject/internal/cmdinterpretator"
	"testproject/internal/dbcontroller"
	//"testproject/internal/clicontroller"
	"testproject/internal/service"
)

func main() {

	db := dbcontroller.NewDBController() // 
	i:= cmdinterpretator.NewCMDInterpetator() //

	s := service.NewService(db, i)
	s.RegisterCommands()

	i.Execute()

}
