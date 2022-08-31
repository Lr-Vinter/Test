package service

import (
	"testproject/internal/dbcontroller"
	//"testproject/internal/clicontroller"
	"testproject/internal/cmdinterpretator"
	"testproject/internal/commands"
)

type Service struct {
	db            *dbcontroller.DBController
	interpretator *cmdinterpretator.CMDInterpretator
}

func NewService(db *dbcontroller.DBController, interpretator *cmdinterpretator.CMDInterpretator) *Service {
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

func (s *Service) GetCMD(args []string) error {
	answer, err := s.db.RetrieveData(args)
	if err != nil {
		return err
	}

	s.interpretator.WriteResponse(answer)
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
	s.interpretator.SendCommandDescription(commands.Set, s.setCMD)
	s.interpretator.SendCommandDescription(commands.Get, s.GetCMD)
	s.interpretator.SendCommandDescription(commands.Exit, s.ExitCMD)
}
