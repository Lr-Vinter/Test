package datastorage

import (
	"bufio"
	"os"
	"errors"
	"log"
)

type DataStorage struct {
	//filelist   *list.List
	fi         *os.File
	filelength int64

	reader *bufio.Reader
	writer *bufio.Writer
}

func NewDataStorage() *DataStorage {
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
		log.Fatal(err)
	}

	d.reader = bufio.NewReader(d.fi)
	_, err = d.reader.Discard(position)
	if err != nil {
		err = errors.New("DS: discard method failed")
		return "", err
	}
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
		return err
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