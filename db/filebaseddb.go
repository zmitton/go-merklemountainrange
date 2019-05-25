package db

import (
	"encoding/binary"
	"errors"
	"os"
	// "time"
)

const wordsize = 64

type Filebaseddb struct {
	Fd *os.File
}

func NewFilebaseddb(path string) *Filebaseddb {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(errors.New("Error opening or creating filedb"))
	}
	db := Filebaseddb{Fd: file}

	readBuffer := make([]byte, wordsize)
	n, err := db.Fd.ReadAt(readBuffer, 0)
	if n == 0 {
		db.SetLeafLength(0)
		db.Fd.Sync()
	}

	return &db
}
func (db Filebaseddb) Get(index int64) []byte {
	return db.get(wordsize * (index + 1))
}
func (db Filebaseddb) Set(index int64, value []byte) {
	db.set(wordsize*(index+1), value)
}
func (db Filebaseddb) GetLeafLength() int64 {
	leafLength := db.get(0)
	return int64(binary.BigEndian.Uint64(leafLength[wordsize-8 : wordsize]))
}
func (db Filebaseddb) SetLeafLength(_value int64) {
	index := int64(wordsize - 8)
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(_value))
	db.set(index, value)
	db.Fd.Sync()
}

func (db Filebaseddb) get(index int64) []byte {
	readBuffer := make([]byte, wordsize)
	n, err := db.Fd.ReadAt(readBuffer, int64(index))

	if err != nil || n == 0 {
		panic(errors.New("Error reading from Filebaseddb"))
	}

	return readBuffer
}
func (db Filebaseddb) set(index int64, value []byte) {
	n, err := db.Fd.WriteAt(value, int64(index))

	if err != nil || n == 0 {
		panic(errors.New("Error writing to Filebaseddb"))
	}
}
