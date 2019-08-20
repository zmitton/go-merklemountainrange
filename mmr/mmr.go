package mmr

import (
	"errors"

	// "fmt"
	"math"

	"github.com/zmitton/go-merklemountainrange/db"
	"github.com/zmitton/go-merklemountainrange/digest"
	"github.com/zmitton/go-merklemountainrange/position"
)

type Mmr struct {
	digest digest.Digest
	db     db.Db
	// consider: do i even need semephore stuff?
}

func New(_digest digest.Digest, _db db.Db) *Mmr {
	return &Mmr{digest: _digest, db: _db}
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
	// fmt.Print("NODEVALUE_ ", len(nodeValue))

	if !ok {
		if p.Height > 0 { // node can be implied by its children
			leftChildValue := mmr.getNodeValue(position.LeftChild(p))
			rightChildValue := mmr.getNodeValue(position.RightChild(p))
			nodeValue = mmr.digest(append(leftChildValue[:], rightChildValue[:]...))
		} else {
			// fmt.Print("GGG", p)
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
		// thing := append(leftValue[:], rightValue[:]...)
		// thing := make([]byte, 64)
		// if
		// sampleLeaf := make([]byte, 64)
		// sampleLeaf = sha3.NewLegacyKeccak256().Sum([]byte{})[0:]
		// sampleLeaf = append(sampleLeaf, sampleLeaf...)

		value := mmr.digest(append(leftValue[:], rightValue[:]...))
		// fmt.Print("LLL ", len(value))
		mmr.db.Set(value, writeIndex)
	}
}
func (mmr *Mmr) setLeafLength(leafLength int64) {
	mmr.db.SetLeafLength(leafLength)
}
func (mmr *Mmr) verifyPath(currentPosition position.Position, currentValue []byte, leafPosition position.Position) []byte { // verifies as it walks
	// fmt.Print("\n\ncurrentpo: ", currentValue)
	// fmt.Print("\n\nlength: ", len(currentValue))

	if currentPosition.Index == leafPosition.Index { // base case
		return currentValue
	} else {
		// fmt.Print("\n\n curpos ", currentPosition, " leafPos ", leafPosition)

		leftChildPosition := position.LeftChild(currentPosition)
		rightChildPosition := position.RightChild(currentPosition)
		leftValue := mmr.getNodeValue(leftChildPosition)
		rightValue := mmr.getNodeValue(rightChildPosition)
		// fmt.Print("\n\n leftChildPos ", leftChildPosition)
		// if !bytes.Equal(currentValue, mmr.digest(leftValue, rightValue)) {
		// fmt.Print("\n\nleafPos: ", leafPosition, "\ncurrentVal ", currentValue[0], "\ndigVal ", mmr.digest(append(leftValue[:], rightValue[:]...)))
		// if bytes.Equal(currentValue, mmr.digest(append(leftValue[:], rightValue[:]...))) {
		// 	panic(errors.New("Hash mismatch of node # and its children"))
		// }
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
	db := db.NewMemorybaseddb(make(map[int64][]byte), referenceTreeLength[0])
	for _, position := range positions {
		db.Set(mmr.getNodeValue(position), position.Index)
	}
	// fmt.Print(db) //for testing only
	return New(mmr.digest, db)
}

func (mmr *Mmr) Get(leafIndex int64) ([]byte, bool) {
	leafLength := mmr.GetLeafLength()
	if leafIndex >= leafLength {
		return []byte{}, false
		// panic(errors.New("Leaf not in tree"))
	}
	// leaf, _ := mmr.db.Get(position.GetNodePosition(leafIndex).Index)
	// return leaf
	return mmr.db.Get(position.GetNodePosition(leafIndex).Index)
}

func (mmr *Mmr) GetVerified(leafIndex int64) []byte {
	leafLength := mmr.GetLeafLength()
	if leafIndex >= leafLength {
		panic(errors.New("Leaf not in tree"))
	}
	leafPosition := position.GetNodePosition(leafIndex)
	localPeakPosition := position.LocalPeakPosition(leafIndex, leafLength)
	// fmt.Print("here")
	// fmt.Print(": ", leafIndex, leafLength, " :")
	localPeakValue := mmr.getNodeValue(localPeakPosition)
	// fmt.Print("there")
	// fmt.Print("\n\nlocalPeakPos", localPeakPosition)
	// fmt.Print("leafLength", mmr.GetLeafLength())
	// fmt.Print("nodeLength", mmr.GetNodeLength())
	// fmt.Print("\n\nlocalpeakVALUE", localPeakValue)
	// fmt.Print("\nleafPOSITION", leafPosition)
	return mmr.verifyPath(localPeakPosition, localPeakValue, leafPosition)
}
func (mmr *Mmr) Append(value []byte, leafIndex ...int64) {
	leafLength := mmr.GetLeafLength()
	if len(leafIndex) == 0 || leafIndex[0] == leafLength {
		nodePosition := position.GetNodePosition(leafLength)
		mountainPositions := position.MountainPositions(position.LocalPeakPosition(leafLength, leafLength), nodePosition.Index)
		if len(value) != 64 {
			// fmt.Print("SHIT", len(value))
		}
		mmr.db.Set(value, nodePosition.Index)
		// fmt.Print("LEAFINDEX", leafIndex, "NODEPOSITION", nodePosition, "MP", mountainPositions)
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

// func (mmr *Mmr) GetDb() db.Db {
// 	return mmr.db
// }
