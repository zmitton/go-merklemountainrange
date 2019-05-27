package position

import (
	"errors"
	"math"
)

type Position struct {
	Index     int64   // index
	Height    float64 // height
	Rightness bool    // whether it is a right (or left) child
}

func LeftChild(p Position) Position {
	if p.Height <= 0 {
		panic(errors.New("Height 0 does not have child"))
	}
	return Position{
		Index:     p.Index - int64(math.Pow(2, p.Height)),
		Height:    p.Height - 1,
		Rightness: false,
	}
}
func RightChild(p Position) Position {
	if p.Height <= 0 {
		panic(errors.New("Height 0 does not have child"))
	}
	return Position{
		Index:     p.Index - 1,
		Height:    p.Height - 1,
		Rightness: true,
	}
}
func Sibling(p Position) Position {
	multiplier := int64(1)
	if p.Rightness {
		multiplier = -1
	}
	return Position{
		Index:     p.Index + multiplier*(int64(math.Pow(2, p.Height+1))-1),
		Height:    p.Height,
		Rightness: !p.Rightness,
	}
}
func ParentIndex(p Position) int64 {
	if p.Rightness {
		return p.Index + 1
	} else {
		return p.Index + int64(math.Pow(2, p.Height+1))
	}
}
func GodPeakFromLeafIndex(_leafIndex int64) Position {
	peakHeight := float64(0)
	leafIndex := float64(_leafIndex)
	for math.Pow(2, peakHeight) <= leafIndex+1 {
		peakHeight++
	}
	return Position{
		Index:     int64(math.Pow(2, (peakHeight+1)) - 2),
		Height:    peakHeight,
		Rightness: false,
	}
}
func GetNodePosition(leafIndex int64) Position {
	currentPosition := GodPeakFromLeafIndex(leafIndex)
	accumulator := int64(0)
	serviceRange := int64(0)
	for currentPosition.Height > 0 {
		serviceRange = int64(math.Pow(2, currentPosition.Height-1))
		if leafIndex >= accumulator+serviceRange {
			currentPosition = RightChild(currentPosition)
			accumulator += serviceRange
		} else {
			currentPosition = LeftChild(currentPosition)
		}
	}
	return currentPosition
}

func hasPosition(nodes map[int64]Position, p Position) bool {
	_, has := nodes[p.Index]
	if !has && p.Height > 0 {
		if hasPosition(nodes, LeftChild(p)) && hasPosition(nodes, RightChild(p)) {
			has = true
		}
	}
	return has
}

func PeakPositions(leafIndex int64) []Position {
	currentPosition := GodPeakFromLeafIndex(leafIndex)
	peakPositions := make([]Position, 0)
	for leafIndex >= 0 {
		currentPosition = LeftChild(currentPosition)
		if leafIndex >= int64(math.Pow(2, currentPosition.Height-1)) {
			peakPositions = append(peakPositions, currentPosition)
			currentPosition = Sibling(currentPosition)
			leafIndex -= int64(math.Pow(2, currentPosition.Height)) // leafIndex becomes a kindof accumulator
		}
	}
	return peakPositions
}

func LocalPeakPosition(leafIndex int64, leafLength int64) Position {
	if leafLength > leafIndex {
		leafIndex = leafLength - 1
	}
	return localPeakPosition(leafIndex, PeakPositions(leafIndex))
}
func localPeakPosition(leafIndex int64, peakPositions []Position) Position {
	currentRange := int64(0)
	localPeak := Position{Index: -1, Height: -1, Rightness: false}
	for i, peakPosition := range peakPositions {
		currentRange = int64(math.Pow(2, peakPosition.Height))
		if leafIndex < currentRange {
			localPeak = peakPositions[i]
			break
		} else {
			leafIndex -= currentRange
		}
	}
	return localPeak
}

func MountainPositions(currentPosition Position, targetIndex int64) [][]Position { // positions to hash after appending
	mountainPositions := make([][]Position, 0)
	for currentPosition.Height > 0 {
		children := []Position{LeftChild(currentPosition), RightChild(currentPosition)}
		mountainPositions = append(mountainPositions, children)
		if targetIndex > currentPosition.Index-int64(math.Pow(2, currentPosition.Height)-currentPosition.Height+1) {
			currentPosition = children[1]
		} else {
			currentPosition = children[0]
		}
	}
	return mountainPositions
}

func ProofPositions(leafIndexes []int64, referenceTreeLength int64) map[int64]Position {
	positions := make(map[int64]Position)
	finalPeakPositions := PeakPositions(referenceTreeLength - 1)
	// add peak positions
	for _, finalPeakPosition := range finalPeakPositions { // log(n)/2
		positions[finalPeakPosition.Index] = finalPeakPosition
	}
	//add local mountain proof positions for each leaf
	for _, leafIndex := range leafIndexes { // k*2log(n) // k is num leaves to prove
		if leafIndex >= referenceTreeLength {
			panic(errors.New("leafIndex must be less than leaf length"))
		}
		nodePosition := GetNodePosition(leafIndex)
		finalLocalPeak := localPeakPosition(leafIndex, finalPeakPositions)
		mountainPositions := MountainPositions(finalLocalPeak, nodePosition.Index)
		for _, mountainPosition := range mountainPositions {
			positions[mountainPosition[0].Index] = mountainPosition[0]
			positions[mountainPosition[1].Index] = mountainPosition[1]
		}
	}
	// find implied positions (ones which can be calculated based on child positions that are present)
	impliedIndexes := make([]int64, 0)
	for k, v := range positions { // k*log(n)
		// for (let j = 0; j < positionIndexes.length; j++) { // k*log(n)
		if v.Height > 0 {
			hasLeftChild := hasPosition(positions, LeftChild(v))
			hasRightChild := hasPosition(positions, RightChild(v))
			if hasLeftChild && hasRightChild {
				// don't remove them yet because recursion will be slower
				impliedIndexes = append(impliedIndexes, k)
			}
		}
	}
	// finally remove implied nodes
	for _, impliedIndex := range impliedIndexes { // k*log(n)
		delete(positions, impliedIndex)
	}
	return positions
}
