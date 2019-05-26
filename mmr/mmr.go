package mmr

import (
	"bytes"
	"errors"
	"math"
	"merklemountainrange/db"
	"merklemountainrange/position"
)

// import "math"

type digest func(args ...[]byte) []byte

type Mmr struct {
	digest      digest // when to use pointers???
	db          db.Db
	_leafLength int64
	// _nodeLength int64
	// todo: add lock / semephore shit
}

func H(args ...[]byte) []byte {
	return []byte{212}
}

func NewMmr(_digest digest, _db db.Db) Mmr {
	this := Mmr{digest: _digest, db: _db, _leafLength: -1}
	return this
}

func (mmr Mmr) GetNodeLength() int64 {
	return position.GetNodePosition(mmr.GetLeafLength()).Index
}
func (mmr Mmr) GetLeafLength() int64 {
	if mmr._leafLength == -1 {
		mmr._leafLength = mmr.db.GetLeafLength()
	}
	return mmr._leafLength
}

func (mmr Mmr) getNodeValue(p position.Position) []byte {
	// caller's responsibility to request a position within leafLength
	nodeValue, ok := mmr.db.Get(p.Index)
	if !ok {
		if p.Height > 0 { // implied node
			leftChildValue := mmr.getNodeValue(position.LeftChild(p))
			rightChildValue := mmr.getNodeValue(position.RightChild(p))
			nodeValue = mmr.digest(leftChildValue, rightChildValue)
		} else {
			panic(errors.New("Missing node in db"))
		}
	}
	return nodeValue
}
func (mmr Mmr) hashUp(positionPairs [][]position.Position) {
	for i := len(positionPairs) - 1; i >= 0; i-- {
		leftValue := mmr.getNodeValue(positionPairs[i][0])
		rightValue := mmr.getNodeValue(positionPairs[i][1])
		writeIndex := position.ParentIndex(positionPairs[i][0])
		mmr.db.Set(writeIndex, mmr.digest(leftValue, rightValue))
	}
}
func (mmr Mmr) setLeafLength(leafLength int64) {
	mmr.db.SetLeafLength(leafLength)
	mmr._leafLength = leafLength
}
func (mmr Mmr) verifyPath(currentPosition position.Position, currentValue []byte, leafPosition position.Position) []byte { // verifies as it walks
	if currentPosition.Index == leafPosition.Index { // base case
		return currentValue
	} else {
		leftChildPosition := position.LeftChild(currentPosition)
		rightChildPosition := position.RightChild(currentPosition)
		leftValue := mmr.getNodeValue(leftChildPosition)
		rightValue := mmr.getNodeValue(rightChildPosition)
		if !bytes.Equal(currentValue, mmr.digest(leftValue, rightValue)) {
			panic(errors.New("Hash mismatch of node # and its children"))
		}
		if leafPosition.Index > currentPosition.Index-int64(math.Pow(2, currentPosition.Height))-int64(currentPosition.Height+1) { //umm yeah, check this line
			return mmr.verifyPath(rightChildPosition, rightValue, leafPosition)
		} else {
			return mmr.verifyPath(leftChildPosition, leftValue, leafPosition)
		}
	}
}

// func (mmr Mmr) GetProof(leafIndexes []int64, referenceTreeLength int64) Mmr{ // returns a sparse MMR containing the leaves specified
// 	let proofMmr
// 	await this.lock.acquire()
// 	try{
// 		referenceTreeLength = referenceTreeLength || await this.getLeafLength()

// 		let positions = MMR.proofPositions(leafIndexes, referenceTreeLength)
// 		let nodes = {}

// 		let nodeIndexes = Object.keys(positions)
// 		await Promise.all(nodeIndexes.map( async (i) => {
// 			let nodeValue = await this._getNodeValue(positions[i])
// 			nodes[i] = nodeValue
// 		}))
// 		proofMmr = new MMR(this.digest, new MemoryBasedDb(referenceTreeLength, nodes))

// 	}finally{
// 		this.lock.release()
// 		return proofMmr
// 	}
// }
