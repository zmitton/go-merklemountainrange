package digest

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"
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

// func FlyHash(input []byte) []byte {
// 	chunks := [][]byte{}
// 	for i := 63; i < len(input); i += 64 {
// 		chunks = append(chunks, input[0:i+1])
// 	}
// 	difficultySum := uint64(0)
// 	for i := 0; i < len(chunks); i++ {
// 		difficultySum += uint64(chunks[i][63])
// 	}
// 	difficultyBytes := make([]byte, 8)
// 	binary.LittleEndian.PutUint64(difficultyBytes, difficultySum)

// 	output := make([]byte, 56)
// 	output = append(output, difficultyBytes...)
// 	if len(output) != 64 || len(input)%64 != 0 {
// 		panic(errors.New("ALERT, input must be 64 byte chucks"))
// 	}
// 	return output
// }

func keccak256(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		hash.Write(d)
	}
	return hash.Sum(nil)
}

func Keccak256FlyHash(input []byte) []byte {
	hasher := keccak256
	// hasher := func(data ...[]byte) []byte {
	// 	hash := sha3.NewLegacyKeccak256()
	// 	for _, d := range data {
	// 		hash.Write(d)
	// 	}
	// 	return hash.Sum(nil)
	// }
	if len(input)%64 != 0 {
		panic(errors.New("ALERT, input must be 64 byte chucks"))
	}
	chunks := [][]byte{}
	for i := 63; i < len(input); i += 64 {
		chunks = append(chunks, input[i-63:i+1])
	}
	difficultySum := uint64(0)
	for i := 0; i < len(chunks); i++ {
		difficultySum += binary.BigEndian.Uint64(chunks[i][56:64])
	}
	difficultyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(difficultyBytes, difficultySum)

	finalHash := hasher(input)

	output := finalHash[0:32]

	output = append(output, make([]byte, 24)...)
	output = append(output, difficultyBytes...)

	if len(output) != 64 {
		panic(errors.New("remove after debugging"))
	}

	return output
}
