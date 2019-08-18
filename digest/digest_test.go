package digest

import (
	"fmt"

	// "github.com/zmitton/go-merklemountainrange/db"
	// "github.com/zmitton/go-merklemountainrange/digest"

	"testing"

	"golang.org/x/crypto/sha3"
	// "crypto/sha256"
	// "github.com/zmitton/go-flyclient"
	// "golang.org/x/crypto/sha3"
)

func TestInstance(t *testing.T) {

	// NewLegacyKeccak256
	t.Run("#Append to in-memory mmr", func(t *testing.T) {
		thing := FlyHash
		thing2 := sha3.NewLegacyKeccak256().Sum
		x := thing([]byte{'z'})
		y := thing2([]byte{})
		fmt.Print("\nJJJJJJ \n", x)
		fmt.Print("\nJJJJJJ \n", fmt.Sprintf("%x", y))
	})

}
