package mmr

import (
	"errors"
	"math"

	"github.com/zmitton/go-merklemountainrange/db"
	"github.com/zmitton/go-merklemountainrange/digest"
	"github.com/zmitton/go-merklemountainrange/position"
)

type Mmr struct {
	digest digest.Digest
	db     db.Db
	// consider: do i even need a semephore?
}

func New(_digest digest.Digest, _db db.Db) *Mmr {
	return &Mmr{digest: _digest, db: _db}
}

func FromSerialized(_digest digest.Digest, serializedDb []byte) *Mmr {
	return New(_digest, db.FromSerialized(serializedDb)) /*uses a memoryBasedDb*/
}

func (mmr *Mmr) GetNodeLength() int64 {
	return position.GetNodePosition(mmr.GetLeafLength()).Index
}
func (mmr *Mmr) GetLeafLength() int64 {
	return mmr.db.GetLeafLength() // caching is handled by db
}

func (mmr *Mmr) getNodeValue(p position.Position) []byte {
	// caller's responsibility to request a position within leafLength
	nodeValue, ok := mmr.db.Get(p.Index)

	if !ok {
		if p.Height > 0 { // nodeValue may be "implied" by its children (exists functionally)
			leftChildValue := mmr.getNodeValue(position.LeftChild(p))
			rightChildValue := mmr.getNodeValue(position.RightChild(p))
			nodeValue = mmr.digest(append(leftChildValue[:], rightChildValue[:]...))
		} else {
			panic(errors.New("Missing node in db"))
		}
	}
	return nodeValue
}
func (mmr *Mmr) hashUp(positionPairs [][]position.Position) {
	for i := len(positionPairs) - 1; i >= 0; i-- {
		leftValue := mmr.getNodeValue(positionPairs[i][0])
		rightValue := mmr.getNodeValue(positionPairs[i][1])
		writeIndex := position.ParentIndex(positionPairs[i][0])
		value := mmr.digest(append(leftValue[:], rightValue[:]...))
		mmr.db.Set(value, writeIndex)
	}
}
func (mmr *Mmr) setLeafLength(leafLength int64) {
	mmr.db.SetLeafLength(leafLength)
}
func (mmr *Mmr) verifyPath(currentPosition position.Position, currentValue []byte, leafPosition position.Position) []byte { // verifies as it walks
	if currentPosition.Index == leafPosition.Index { // base case
		return currentValue
	} else {

		leftChildPosition := position.LeftChild(currentPosition)
		rightChildPosition := position.RightChild(currentPosition)
		leftValue := mmr.getNodeValue(leftChildPosition)
		rightValue := mmr.getNodeValue(rightChildPosition)

		if leafPosition.Index > currentPosition.Index-int64(math.Pow(2, currentPosition.Height))-int64(currentPosition.Height)+1 {
			return mmr.verifyPath(rightChildPosition, rightValue, leafPosition)
		} else {
			return mmr.verifyPath(leftChildPosition, leftValue, leafPosition)
		}
	}
}

// returns a sparse MMR containing the leaves specified
func (mmr *Mmr) GetProof(leafIndexes []int64, referenceTreeLength ...int64) *Mmr {
	if len(referenceTreeLength) == 0 { // variatic hack
		referenceTreeLength = append(referenceTreeLength, mmr.GetLeafLength())
	}
	positions := position.ProofPositions(leafIndexes, referenceTreeLength[0])
	nodes := make(map[int64][]byte)
	db := db.NewMemorybaseddb(referenceTreeLength[0], nodes)
	for _, position := range positions {
		db.Set(mmr.getNodeValue(position), position.Index)
	}
	return New(mmr.digest, db)
}
func (mmr *Mmr) Serialize() []byte {
	return mmr.db.Serialize()
}

// Use `Get()` in most cases (to be safe). Only use this for extra speed when
// interacting with a full mmr and not dealing with proofs.
func (mmr *Mmr) GetUnverified(leafIndex int64) ([]byte, bool) {
	leafLength := mmr.GetLeafLength()
	if leafIndex >= leafLength {
		return []byte{}, false
	}
	return mmr.db.Get(position.GetNodePosition(leafIndex).Index)
}

func (mmr *Mmr) Get(leafIndex int64) []byte {
	leafLength := mmr.GetLeafLength()
	if leafIndex >= leafLength {
		panic(errors.New("Leaf not in tree"))
	}
	leafPosition := position.GetNodePosition(leafIndex)
	localPeakPosition := position.LocalPeakPosition(leafIndex, leafLength)
	localPeakValue := mmr.getNodeValue(localPeakPosition)

	return mmr.verifyPath(localPeakPosition, localPeakValue, leafPosition)
}
func (mmr *Mmr) Append(value []byte, leafIndex ...int64) {
	leafLength := mmr.GetLeafLength()
	if len(leafIndex) == 0 || leafIndex[0] == leafLength {
		nodePosition := position.GetNodePosition(leafLength)
		localPeak := position.LocalPeakPosition(leafLength, leafLength)
		mountainPositions := position.MountainPositions(localPeak, nodePosition.Index)
		mmr.db.Set(value, nodePosition.Index)
		mmr.hashUp(mountainPositions)
		mmr.setLeafLength(leafLength + 1)
	} else {
		panic(errors.New("Can only append to end of MMR"))
	}
}
func (mmr *Mmr) AppendMany(values [][]byte, startLeafIndex ...int64) {
	if len(startLeafIndex) == 0 {
		startLeafIndex = append(startLeafIndex, mmr.GetLeafLength())
	}
	for i, value := range values {
		mmr.Append(value, startLeafIndex[0]+int64(i))
	}
}
func (mmr *Mmr) GetRoot(leafIndex ...int64) []byte {
	var peakValues []byte
	if len(leafIndex) == 0 {
		leafIndex = append(leafIndex, mmr.GetLeafLength()-1)
	}
	peakPositions := position.PeakPositions(leafIndex[0])
	for _, peakPosition := range peakPositions {
		peakValues = append(peakValues, mmr.getNodeValue(peakPosition)...)
	}
	// note: a single peak differs from its MMR root in that it gets hashed a second time
	return mmr.digest(peakValues)
}

// logically deletes everything after (and including) leafIndex.
// todo: consider side affects. test more
func (mmr *Mmr) Delete(leafIndex int64) {
	leafLength := mmr.GetLeafLength()
	if leafIndex < leafLength {
		mmr.setLeafLength(leafIndex)
	}
}
func (mmr *Mmr) Db() db.Db {
	return mmr.db
}
