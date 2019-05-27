package mmr

import (
	"bytes"
	"errors"
	"math"
	"merklemountainrange/db"
	"merklemountainrange/digest"
	"merklemountainrange/position"
	"sync"
)

type Mmr struct {
	digest     digest.Digest
	db         db.Db
	leafLength int64
	// consider: do i even need semephore stuff
}

func H(args ...[]byte) []byte {
	return []byte{212}
}

// leafLength of -1 for new
func NewMmr(_digest digest.Digest, _db db.Db) Mmr {
	this := Mmr{digest: _digest, db: _db}
	return this
}

func (mmr Mmr) GetNodeLength() int64 {
	return position.GetNodePosition(mmr.GetLeafLength()).Index
}
func (mmr Mmr) GetLeafLength() int64 {
	// if mmr.leafLength == 0 {
	// 	mmr.leafLength = mmr.db.GetLeafLength()
	// }
	// return mmr.leafLength
	return mmr.db.GetLeafLength()
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
	// mmr.leafLength = leafLength
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

func (mmr Mmr) GetProof(leafIndexes []int64, referenceTreeLength int64) Mmr { // returns a sparse MMR containing the leaves specified
	if referenceTreeLength == -1 { // variatic hack
		referenceTreeLength = mmr.GetLeafLength()
	}

	positions := position.ProofPositions(leafIndexes, referenceTreeLength)
	db := db.NewMemorybaseddb(make(map[int64][]byte), referenceTreeLength)
	var wg sync.WaitGroup
	wg.Add(len(positions))
	for _, position := range positions {
		go func() {
			db.Set(position.Index, mmr.getNodeValue(position))
			wg.Done()
		}()
	}
	wg.Wait()

	return NewMmr(mmr.digest, db)
}

func (mmr Mmr) Get(leafIndex int64) []byte {
	leafLength := mmr.GetLeafLength()
	if leafIndex >= leafLength {
		panic(errors.New("Leaf not in tree"))
	}
	leafPosition := position.GetNodePosition(leafIndex)
	localPeakPosition := position.LocalPeakPosition(leafIndex, leafLength)
	localPeakValue := mmr.getNodeValue(localPeakPosition)
	return mmr.verifyPath(localPeakPosition, localPeakValue, leafPosition)
}
func (mmr Mmr) Append(value []byte, leafIndex int64) {
	leafLength := mmr.GetLeafLength()
	if leafIndex == -1 || leafIndex == leafLength {
		nodePosition := position.GetNodePosition(leafLength)
		mountainPositions := position.MountainPositions(position.LocalPeakPosition(leafLength, leafLength), nodePosition.Index)
		mmr.db.Set(nodePosition.Index, value)
		mmr.hashUp(mountainPositions)
		mmr.setLeafLength(leafLength + 1)
	} else {
		panic(errors.New("Can only append to end of MMR"))
	}
}
func (mmr Mmr) AppendMany(values [][]byte, startLeafIndex int64) {
	if startLeafIndex == -1 {
		startLeafIndex = mmr.GetLeafLength()
	}
	for i, value := range values {
		mmr.Append(value, startLeafIndex+int64(i))
	}
}
func (mmr Mmr) GetRoot(leafIndex int64) []byte {
	var peakValues [][]byte
	if leafIndex == -1 {
		leafIndex = mmr.GetLeafLength() - 1
	}
	peakPositions := position.PeakPositions(leafIndex)
	for _, peakPosition := range peakPositions {
		peakValues = append(peakValues, mmr.getNodeValue(peakPosition))
	}
	// note: a single peak differs from its MMR root in that it gets hashed a second time
	return mmr.digest(peakValues...)
}

// logically deletes everything after (and including) leafIndex.
// todo: consider side affects. test more
func (mmr Mmr) Delete(leafIndex int64) {
	leafLength := mmr.GetLeafLength()
	if leafIndex < leafLength {
		mmr.setLeafLength(leafIndex)
	}
}
