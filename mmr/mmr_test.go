package mmr

import (
	"fmt"
	"time"

	"github.com/zmitton/go-merklemountainrange/db"
	"github.com/zmitton/go-merklemountainrange/digest"

	"testing"
	// "crypto/sha256"
	// "github.com/zmitton/go-flyclient"
)

func TestInstance(t *testing.T) {
	fileBasedDb1 := db.OpenFilebaseddb("../db/testdata/etcleafdata.mmr")
	fileBasedMmr1 := NewMmr(digest.FlyHash, fileBasedDb1)

	memoryBasedDb1 := db.NewMemorybaseddb(map[int64][]byte{}, 0)
	memoryBasedMmr1 := NewMmr(digest.FlyHash, memoryBasedDb1)

	leafLength := fileBasedDb1.GetLeafLength()
	etcLeafData := make([][]byte, 0)
	for i := int64(0); i < leafLength; i++ {
		etcLeafData = append(etcLeafData, fileBasedMmr1.Get(i))
	}

	t.Run("#Append to mmr", func(t *testing.T) {
		leafLength := fileBasedDb1.GetLeafLength()
		for i := int64(0); i < leafLength; i++ {
			memoryBasedMmr1.Append(fileBasedMmr1.Get(i), i)
		}
	})

	t.Run("#GetLeafLength", func(t *testing.T) {
		leafLength1 := fileBasedMmr1.GetLeafLength()
		leafLength2 := memoryBasedMmr1.GetLeafLength()
		if leafLength1 != 1000 {
			t.Errorf("leafLength1 should be equal to 1000, not %d", leafLength1)
		}
		if leafLength2 != 1000 {
			t.Errorf("leafLength2 should be equal to 1000, not %d", leafLength2)
		}
	})

	t.Run("#GetNodeLength", func(t *testing.T) {
		nodeLength1 := fileBasedMmr1.GetNodeLength()
		nodeLength2 := memoryBasedMmr1.GetNodeLength()
		if nodeLength1 != 1994 {
			t.Errorf("nodeLength1 should be equal to 1994, not %d", nodeLength1)
		}
		if nodeLength2 != 1994 {
			t.Errorf("nodeLength2 should be equal to 1994, not %d", nodeLength2)
		}
	})

	t.Run("Performance/Benchmarks", func(t *testing.T) {
		const NUM_LOOPS = 1000

		// in-memory based
		preTime := time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			memoryBasedMmr1.Get(i)

		}
		fmt.Print("\nin-memory #Get (w ~1000 leaves):        \t", time.Since(preTime)/NUM_LOOPS)

		preTime = time.Now()
		// memoryBasedMmr1.GetVerified(0)
		for i := int64(0); i < NUM_LOOPS; i++ {
			memoryBasedMmr1.GetVerified(i)
		}
		fmt.Print("\nin-memory #GetVerified (w ~1000 leaves):\t", time.Since(preTime)/NUM_LOOPS)

		// sampleLeaf := digest.FlyHash([]byte{'z'})
		// sampleLeaf := make([]byte, 64)
		sampleLeaf := digest.FlyHash(make([]byte, 128))
		// fmt.Print("CCCC")
		// sampleLeaf = digest.FlyHash([]byte{})[0:]
		// sampleLeaf = append(sampleLeaf, sampleLeaf...)
		// memoryBasedMmr1.Append(sampleLeaf, 1000)
		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			memoryBasedMmr1.Append(sampleLeaf, i+1000)
		}
		fmt.Print("\nin-memory #Append (w ~1000 leaves):      \t", time.Since(preTime)/NUM_LOOPS)

		//file based
		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			fileBasedMmr1.Get(i)
		}
		fmt.Print("\n\nfile-based #Get (w ~1000 leaves):       \t", time.Since(preTime)/NUM_LOOPS)

		tempFileBasedDb := db.CreateFilebaseddb("../db/testdata/temp.mmr", 64)
		tempFileBasedMmr := NewMmr(digest.FlyHash, tempFileBasedDb)

		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			fileBasedMmr1.GetVerified(i)
		}
		fmt.Print("\nfile based #GetVerified (w ~1000 leaves):\t", time.Since(preTime)/NUM_LOOPS)

		// tempFileBasedMmr.Append(sampleLeaf, 1)
		// fmt.Print(tempFileBasedMmr.Get(0))
		preTime = time.Now()
		// tempFileBasedDb.Set(sampleLeaf, 0)
		for i := int64(0); i < NUM_LOOPS; i++ {
			// fmt.Print("HERE", i)
			tempFileBasedMmr.Append(sampleLeaf, i)
			// tempFileBasedMmr.setLeafLength(1000)
			// tempFileBasedDb.SetLeafLength(1000)

			// tempFileBasedDb.Set(sampleLeaf, i)
		}

		fmt.Print("\nfile based #Append (w ~1000 leaves):     \t", time.Since(preTime)/NUM_LOOPS, "\n\n")

		// tempFileBasedMmr.Append(sampleLeaf, -1)
		// tempFileBasedMmr.Append(sampleLeaf, 1000)

		// os.Remove("../db/testdata/temp.mmr")

	})

}
