package mmr

import (
	"errors"
	"math"
	"merklemountainrange/db"
	"merklemountainrange/position"
)

// import "math"

type digest func(p [][]byte) []byte

type Mmr struct {
	digest      digest // when to use pointers???
	db          db.Db
	_leafLength int64
	_nodeLength int64
	// todo: add lock / semephore shit
}

func H(p [][]byte) []byte {
	return []byte{212}
}

func NewMmr(_digest digest, _db db.Db) *Mmr {
	this := Mmr{digest: _digest, db: _db, _leafLength: -1, _nodeLength: -1}
	return &this
}

func (mmr Mmr) GetLeafLength() int64 {
	if mmr._leafLength == -1 {
		mmr._leafLength = mmr.db.GetLeafLength()
	}
	return mmr._leafLength
}

func LeftChild(p *position.Position) *position.Position {
	if p.Height <= 0 {
		panic(errors.New("Height 0 does not have child"))
	}
	return &position.Position{
		Index:     p.Index - int64(math.Pow(2, float64(p.Height))),
		Height:    p.Height - 1,
		Rightness: false,
	}
}
func RightChild(p *position.Position) *position.Position {
	if p.Height <= 0 {
		panic(errors.New("Height 0 does not have child"))
	}
	return &position.Position{
		Index:     p.Index - 1,
		Height:    p.Height - 1,
		Rightness: true,
	}
}
func Sibling(p *position.Position) *position.Position {
	multiplier := int64(1)
	if p.Rightness {
		multiplier = -1
	}
	return &position.Position{
		Index:     p.Index + multiplier*(int64(math.Pow(2, float64(p.Height+1)))-1),
		Height:    p.Height,
		Rightness: !p.Rightness,
	}
}
func ParentIndex(p *position.Position) int64 {
	if p.Rightness {
		return p.Index + 1
	} else {
		return p.Index + int64(math.Pow(2, float64(p.Height+1)))
	}
}

