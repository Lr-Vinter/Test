package service

import (
	"testproject/internal/dbcontroller"
	"testproject/internal/clicontroller"
	"testproject/internal/cmdinterpretator"
	"testproject/internal/commands"
	"fmt"
)

type Service struct {
	db            *dbcontroller.DBController
	clicontroller *clicontroller.CLIController
	interpretator *cmdinterpretator.CMDInterpretator
}

func NewService(db *dbcontroller.DBController, clicontroller *clicontroller.CLIController, interpretator *cmdinterpretator.CMDInterpretator) *Service {
	return &Service{
		db:            db,
		clicontroller: clicontroller,
		interpretator: interpretator,
	}
}

func (s *Service) Execute() {
	for {
		err := s.clicontroller.ReceiveAndRunCMD()
		if(err != nil) {
			fmt.Println(err)
		}

		cmd, args := s.clicontroller.ReturnCmdAndArgs()
		err = s.interpretator.Run(cmd, args)
		if(err != nil) {
			fmt.Println(err)
		}
	}
}

func (s *Service) SetCMD(args []string) error {
	err := s.db.AddData(args)
	if err != nil {
		return err
	}

	//s.clicontroller.WriteResponse("")
	return nil
}

func (s *Service) GetCMD(args []string) error {
	answer, err := s.db.RetrieveData(args)
	if err != nil {
		return err
	}

	s.clicontroller.WriteResponse(answer)
	return nil
}

func (s *Service) ExitCMD(args []string) error {
	err := s.db.Close(args)
	if err != nil {
		return err
	}

	//s.clicontroller.WriteResponse("exit")
	return nil
}

func (s *Service) RegisterCommands() { //
	s.interpretator.SendCommandDescription(commands.Set, s.SetCMD)
	s.interpretator.SendCommandDescription(commands.Get, s.GetCMD)
	s.interpretator.SendCommandDescription(commands.Exit, s.ExitCMD)
}
