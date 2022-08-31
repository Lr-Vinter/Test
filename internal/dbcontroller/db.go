package dbcontroller

import (
	"testproject/internal/datastorage"
	"testproject/internal/index"
	"errors"
	"strings"
)

type DBController struct {
	datastorage   *datastorage.DataStorage
	index         *index.Index
}

func NewDBController() *DBController {
	db := new(DBController)

	db.datastorage = datastorage.NewDataStorage() // DB
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

	return err
}

func (db *DBController) Close(args []string) error {
	return nil
}

func (db *DBController) RetrieveData(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("DB : wrong arguments number for retrieve data func")
	}

	key := args[0]
	position := db.index.Indexmap[key]

	answer, err := db.datastorage.GetDataByPos("data.txt", position)
	if err != nil {
		err = errors.New("DB : failed to load data")
	}
	
	return answer, err
}
