package main

import (
	"testproject/internal/cmdinterpretator"
	"testproject/internal/dbcontroller"
	"testproject/internal/clicontroller"
	"testproject/internal/service"
)

func main() {

	db := dbcontroller.NewDBController() // 
	i:= cmdinterpretator.NewCMDInterpetator() //
	c := clicontroller.NewController()

	s := service.NewService(db, c, i)
	s.RegisterCommands()

	s.Execute()
	//for {
	//	c.ReceiveAndRunCMD()
	//	cmd, args := c.ReturnCmdAndArgs()
	//	i.Run(cmd, args)
	//}

}
