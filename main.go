package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"testproject/internal/cmdinterpretator"
	"testproject/internal/commands"
	"testproject/internal/index"
	"testproject/internal/inputanalyzer"
)

//*****

type CLIController struct {
	inputanalyzer *inputanalyzer.InputAnalyzer
	exec          *cmdinterpretator.CMDInterpretator

	answer string // answer from DBController
	cmd    commands.Command
	args   []string
}

func NewController(db *DBController) *CLIController {
	c := new(CLIController)
	c.inputanalyzer = new(inputanalyzer.InputAnalyzer) //  ?
	c.inputanalyzer = inputanalyzer.NewInputAnalyzer(os.Stdin, 100)
	c.exec = cmdinterpretator.NewCMDInterpetator()
	return c
}

func (c *CLIController) WriteResponse(answer string) { //
	fmt.Println(answer)
}

func (c *CLIController) ReceiveAndRunCMD() error {
	var err error
	err = c.inputanalyzer.ParseInput()
	err = c.inputanalyzer.GetCmdWithArgs(&c.cmd, &c.args)
	if err != nil {
		fmt.Println(err)
	}

	err = c.exec.Run(c.cmd, c.args)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.answer)

	return err
}

//*****
//*****

type Service struct {
	db            *DBController
	clicontroller *CLIController
	interpretator *cmdinterpretator.CMDInterpretator
}

func NewService(db *DBController, clicontroller *CLIController, interpretator *cmdinterpretator.CMDInterpretator) *Service {
	return &Service{
		db:            db,
		clicontroller: clicontroller,
		interpretator: interpretator,
	}
}

func (s *Service) setCMD(args []string) error {
	err := s.db.AddData(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//s.clicontroller.WriteResponse("")
	return nil
}

func (s *Service) getCMD(args []string) error {
	answer, err := s.db.RetrieveData(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.clicontroller.WriteResponse(answer)
	return nil
}

func (s *Service) exitCMD(args []string) error {
	err := s.db.Close(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//s.clicontroller.WriteResponse("exit")
	return nil
}

func (s *Service) registerCommands() error { //
	s.interpretator.SendCommandDescription(commands.Set, s.setCMD)
	s.interpretator.SendCommandDescription(commands.Set, s.getCMD)
	s.interpretator.SendCommandDescription(commands.Set, s.exitCMD)
	return nil
}

type DBController struct {
	clicontroller *CLIController
	datastorage   *DataStorage
	index         *index.Index
}

func NewDBController() *DBController {
	db := new(DBController)

	db.datastorage = new(DataStorage)
	db.datastorage = SetDataStorage() // DB
	db.index = index.NewIndexMap()

	return db
}

func (db *DBController) AddData(args []string) error {
	startposition := db.datastorage.GetFileLength("data.txt")
	db.index.Indexmap[args[0]] = int(startposition) + 1 + len(args[0])

	writedata := strings.Join(args[:], " ")
	err := db.datastorage.WriteData("data.txt", writedata+string("\n"))
	if err != nil {
		err = errors.New("DB : Failed to Execute Set CMD")
	}

	db.clicontroller.answer = "DB : Set successfully completed"
	return err
}

func (db *DBController) Close(args []string) error {
	return nil
}

func (db *DBController) RetrieveData(args []string) (string, error) {
	if len(args) != 1 {
		errstr := "DB : wrong arguments number for retrieve data func"
		err := errors.New(errstr)
		return errstr, err
	}

	key := args[0]
	position := db.index.Indexmap[key]

	answer, err := db.datastorage.GetDataByPos("data.txt", position)
	if err != nil {
		err = errors.New("DB : failed to load data")
	}

	return answer, err
}

//***

type DataStorage struct {
	//filelist   *list.List
	fi         *os.File
	filelength int64

	reader *bufio.Reader
	writer *bufio.Writer
}

func SetDataStorage() *DataStorage {
	d := new(DataStorage)
	//first, _ := os.Create("data.txt")

	//d.filelist = list.New()
	//d.filelist.PushBack(first)

	return d
}

func (ds *DataStorage) AddFile(filename string) error {
	//fi, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//fi, err := os.Create(filename)
	//if err != nil {
	//	err = errors.New("failed to create file")
	//	return err
	//}
	//ds.filelist.PushBack(fi)
	//fi.Close()

	return nil
}

func (d *DataStorage) GetDataByPos(filename string, position int) (string, error) {
	var err error // а нельзя инициализировать в некст строчке (вроде)
	d.fi, err = os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		err = errors.New("DS: failed to open file")
	}

	d.reader = bufio.NewReader(d.fi)
	d.reader.Discard(position)
	data, err := d.reader.ReadBytes(10)
	if err != nil {
		err = errors.New("DS: failed to read data from file")
	}

	return string(data), err
}

func (d *DataStorage) WriteData(filename string, data string) error {
	var err error
	d.fi, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = errors.New("DS: failed to open file")
	}

	d.writer = bufio.NewWriter(d.fi)
	_, err = d.writer.WriteString(data)
	if err != nil {
		err = errors.New("DS: failed to write data to file")
	}

	d.writer.Flush()
	d.filelength, _ = d.fi.Seek(0, 1)
	d.fi.Close()
	return err
}

func (d *DataStorage) GetFileLength(filename string) int64 {
	return d.filelength
}

//***
// index
//***

func main() {

	db := NewDBController()
	c := NewController(db)
	i := cmdinterpretator.NewCMDInterpetator()
	s := NewService(db, c, i)

	s.registerCommands()

	for {
		s.db.clicontroller.ReceiveAndRunCMD()
	}

}
