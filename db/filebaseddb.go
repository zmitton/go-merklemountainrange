/*
This uses a standard file. The entire file is treated as an array where
nodes are stored as 64byte elements and `get/set` operations are done
by reading random accessed data at their `index` multiplied by the word
size of `64` (plus one word because the first 64 bytes holds the
`leafLength` data)
*/

package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	// "time"
)

type Filebaseddb struct {
	fd               *os.File
	cachedLeafLength int64
}

func NewFilebaseddb(path string) *Filebaseddb {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(errors.New("Error opening or creating filedb"))
	}
	db := Filebaseddb{fd: file}

	readBuffer := make([]byte, wordsize)
	n, err := db.fd.ReadAt(readBuffer, 0)
	var leafLength int64
	if n == 0 {
		leafLength = 0
	} else {
		leafLength = db.GetLeafLength()
	}
	db.SetLeafLength(leafLength) // sets the cache

	return &db
}
func (db Filebaseddb) Get(index int64) (value []byte, ok bool) {
	return db.get(wordsize * (index + 1))
}
func (db Filebaseddb) Set(index int64, value []byte) {
	db.set(wordsize*(index+1), value)
}
func (db Filebaseddb) GetLeafLength() int64 {
	if db.cachedLeafLength == 0 {
		leafLengthBuffer, _ := db.get(0)
		db.cachedLeafLength = int64(binary.BigEndian.Uint64(leafLengthBuffer[wordsize-8 : wordsize]))
	}
	return db.cachedLeafLength
}
func (db Filebaseddb) SetLeafLength(leafLength int64) {
	index := int64(wordsize - 8)
	leafLengthBuffer := make([]byte, 8)
	binary.BigEndian.PutUint64(leafLengthBuffer, uint64(leafLength))
	db.cachedLeafLength = leafLength
	db.set(index, leafLengthBuffer)
	db.fd.Sync() // save/flush only after length data is updated
}

func (db Filebaseddb) get(index int64) ([]byte, bool) {
	value := make([]byte, wordsize)
	n, err := db.fd.ReadAt(value, int64(index))

	if err != nil || n == 0 {
		panic(errors.New("Error reading from Filebaseddb"))
	}
	ok := bytes.Equal(make([]byte, wordsize), value)

	return value, ok
}
func (db Filebaseddb) set(index int64, value []byte) {
	n, err := db.fd.WriteAt(value, int64(index))

	if err != nil || n == 0 {
		panic(errors.New("Error writing to Filebaseddb"))
	}
}
