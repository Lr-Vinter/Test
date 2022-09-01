package service

import (
	"testproject/internal/dbcontroller"

	//"testproject/internal/clicontroller"
	"testproject/internal/cli"
	"testproject/internal/commands"
)

type Service struct {
	db            *dbcontroller.DBController
	interpretator *cli.CLIcontroller
}

func NewService(db *dbcontroller.DBController, interpretator *cli.CLIcontroller) *Service {
	return &Service{
		db:            db,
		interpretator: interpretator,
	}
}

func (s *Service) setCMD(args []string) error {
	err := s.db.AddData(args)
	if err != nil {
		return err
	}

	//s.clicontroller.WriteResponse("")
	return nil
}

func (s *Service) getCMD(args []string) error {
	answer, err := s.db.RetrieveData(args)
	if err != nil {
		return err
	}

	//fmt.Println("value", answer)
	s.interpretator.WriteResponse(answer)
	return nil
}

func (s *Service) exitCMD(args []string) error {
	err := s.db.Close(args)
	if err != nil {
		return err
	}

	//s.clicontroller.WriteResponse("exit")
	return nil
}

func (s *Service) RegisterCommands() { //
	s.interpretator.SendCommandDescription(commands.Set, s.setCMD)
	s.interpretator.SendCommandDescription(commands.Get, s.getCMD)
	s.interpretator.SendCommandDescription(commands.Exit, s.exitCMD)
}
