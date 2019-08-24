package digest

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"
)

type Digest func(input []byte) []byte

// // converts output of hash func from [32]byte (fixed array) to []byte (slice)
// func Wrapper32(hashFunc func([]byte) [32]byte) func([]byte) []byte {
// 	return func(input []byte) []byte {
// 		fixedOutput := hashFunc(input)
// 		return fixedOutput[:]
// 	}
// }
// func Wrapper64(hashFunc func([]byte) [64]byte) func([]byte) []byte {
// 	return func(input []byte) []byte {
// 		fixedOutput := hashFunc(input)
// 		return fixedOutput[:]
// 	}
// }

func keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func Keccak256FlyHash(input []byte) []byte {
	if len(input)%64 != 0 {
		panic(errors.New("ALERT, input must be 64 byte chucks for this hash func"))
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

	finalHash := keccak256(input)

	output := finalHash[0:32]

	output = append(output, make([]byte, 24)...)
	output = append(output, difficultyBytes...)

	return output
}
