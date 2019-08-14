package digest

import (
	"encoding/binary"
	"errors"
)

type Digest func(input []byte) []byte

// type Db interface {
// 	Get(int64) ([]byte, bool)
// 	Set(int64, []byte)
// 	GetLeafLength() int64
// 	SetLeafLength(int64)
// }

// converts output of hash func from [32]byte (fixed array) to []byte (slice)
func Wrapper32(hashFunc func([]byte) [32]byte) func([]byte) []byte {
	return func(input []byte) []byte {
		fixedOutput := hashFunc(input)
		return fixedOutput[:]
	}
}
func Wrapper64(hashFunc func([]byte) [64]byte) func([]byte) []byte {
	return func(input []byte) []byte {
		fixedOutput := hashFunc(input)
		return fixedOutput[:]
	}
}

func FlyHash(input []byte) []byte {
	chunks := [][]byte{}
	for i := 63; i < len(input); i += 64 {
		chunks = append(chunks, input[0:i+1])
	}
	difficultySum := uint64(0)
	for i := 0; i < len(chunks); i++ {
		difficultySum += uint64(chunks[i][63])
	}
	difficultyBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(difficultyBytes, difficultySum)

	output := make([]byte, 56)
	output = append(output, difficultyBytes...)
	if len(output) != 64 {
		panic(errors.New("ALERT"))
	}
	return output
}

func FlyHash1(input []byte) []byte {

	a := uint64(input[63])
	b := uint64(input[127])

	c := make([]byte, 8)
	binary.LittleEndian.PutUint64(c, a+b)
	output := make([]byte, 56)
	output = append(output, c...)
	if len(output) != 64 {
		panic(errors.New("ALERT"))
	}
	return output
}
