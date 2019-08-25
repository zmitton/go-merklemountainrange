package mmr

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/zmitton/go-merklemountainrange/db"
	"github.com/zmitton/go-merklemountainrange/digest"

	"testing"
)

func TestInstance(t *testing.T) {
	os.Remove("../db/testdata/temp.mmr")

	fileBasedDb1 := db.OpenFilebaseddb("../db/testdata/etcleafdata.mmr")
	fileBasedMmr1 := New(digest.Keccak256FlyHash, fileBasedDb1)

	memoryBasedDb1 := db.NewMemorybaseddb(map[int64][]byte{}, 0)
	memoryBasedMmr1 := New(digest.Keccak256FlyHash, memoryBasedDb1)

	tempFileBasedDb := db.CreateFilebaseddb("../db/testdata/temp.mmr", 64)
	tempFileBasedMmr := New(digest.Keccak256FlyHash, tempFileBasedDb)

	leafLength := fileBasedDb1.GetLeafLength()
	etcLeafData := make([][]byte, 0)

	sampleLeaf0, _ := fileBasedMmr1.GetUnverified(0)
	sampleLeaf1, _ := fileBasedMmr1.GetUnverified(1)

	for i := int64(0); i < leafLength; i++ {
		leaf, _ := fileBasedMmr1.GetUnverified(i)
		etcLeafData = append(etcLeafData, leaf)
	}

	t.Run("#Append to in-memory mmr", func(t *testing.T) {
		leafLength := fileBasedDb1.GetLeafLength()
		for i := int64(0); i < leafLength; i++ {
			leaf, _ := fileBasedMmr1.GetUnverified(i)
			memoryBasedMmr1.Append(leaf, i)
		}
	})

	t.Run("#GetLeafLengths of both ", func(t *testing.T) {
		leafLength1 := fileBasedMmr1.GetLeafLength()
		leafLength2 := memoryBasedMmr1.GetLeafLength()
		if leafLength1 != 1000 {
			t.Errorf("leafLength1 should be equal to 1000, not %d", leafLength1)
		}
		if leafLength2 != 1000 {
			t.Errorf("leafLength2 should be equal to 1000, not %d", leafLength2)
		}
	})

	t.Run("#GetNodeLengths of both", func(t *testing.T) {
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
			memoryBasedMmr1.GetUnverified(i)
		}
		fmt.Print("\nin-memory #GetUnverified (w ~1000 leaves):\t", time.Since(preTime)/NUM_LOOPS)

		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			memoryBasedMmr1.Get(i)
		}
		fmt.Print("\nin-memory #Get (w ~1000 leaves):\t\t", time.Since(preTime)/NUM_LOOPS)

		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			// memoryBasedMmr1.Append(sampleLeaf0, i+1000)
		}
		fmt.Print("\nin-memory #Append (w ~1000 leaves):      \t", time.Since(preTime)/NUM_LOOPS)

		//file based
		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			fileBasedMmr1.GetUnverified(i)
		}
		fmt.Print("\n\nfile-based #GetUnverified (w ~1000 leaves):\t", time.Since(preTime)/NUM_LOOPS)

		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			fileBasedMmr1.Get(i)
		}
		fmt.Print("\nfile based #Get (w ~1000 leaves):\t\t", time.Since(preTime)/NUM_LOOPS)
		preTime = time.Now()
		for i := int64(0); i < NUM_LOOPS; i++ {
			tempFileBasedMmr.Append(etcLeafData[i], i)
		}
		fmt.Print("\nfile based #Append (w ~1000 leaves):     \t", time.Since(preTime)/NUM_LOOPS, "\n\n")
	})

	t.Run("#a single parent is correct digest of its chilren", func(t *testing.T) {
		x := []byte{}
		copy(x, sampleLeaf0)
		x = digest.Keccak256FlyHash(append(sampleLeaf0, sampleLeaf1...))
		y, _ := fileBasedDb1.Get(2)
		if !bytes.Equal(x, y) {
			t.Errorf("got incorrect hash parent for H(leaf0,leaf1)")
		}
	})
	t.Run("#full replication of etcTestData", func(t *testing.T) {
		nodeLength := fileBasedMmr1.GetNodeLength()

		for i := int64(0); i < nodeLength; i++ {
			fixtureNode, _ := fileBasedDb1.Get(i)
			tempNode, _ := tempFileBasedDb.Get(i)
			if !bytes.Equal(fixtureNode, tempNode) {
				fmt.Print("Expected\n", fmt.Sprintf("%x", fixtureNode))
				fmt.Print("\nGot\n", fmt.Sprintf("%x", tempNode))
				t.Errorf("got incorrect db node")
			}
		}
	})

	t.Run("#GetRoot of leaf 999", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(999)
		if fmt.Sprintf("%x", root) != "1d6e5c69d70d3ac8847ccf63f61303f607382bd988d0d8b559ce53e3305e7b6700000000000000000000000000000000000000000000000000001400691fd2d6" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
	})
	t.Run("#GetRoot of leaf 0", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(0)
		node0, _ := memoryBasedDb1.Get(0)
		if fmt.Sprintf("%x", root) != "4ef0d4100c84abf7f877cde7ae268676b3bab9341cdac33ae7c5de5ca8d865660000000000000000000000000000000000000000000000000000000400000000" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
		if !bytes.Equal(root, digest.Keccak256FlyHash(node0)) {
			t.Errorf(fmt.Sprintf("%x", digest.Keccak256FlyHash(node0)))
		}
	})
	t.Run("#GetRoot of leaf 1", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(1)
		node2, _ := memoryBasedDb1.Get(2)
		if fmt.Sprintf("%x", root) != "b31315e249f2d6814113099367c3ab9092eb94a9e6ea2f3539b78a2da8589ee400000000000000000000000000000000000000000000000000000007ff800000" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
		if !bytes.Equal(root, digest.Keccak256FlyHash(node2)) {
			t.Errorf(fmt.Sprintf("%x", digest.Keccak256FlyHash(node2)))
		}
	})
	t.Run("#GetRoot of leaf 2", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(2)
		node2, _ := memoryBasedDb1.Get(2)
		node3, _ := memoryBasedDb1.Get(3)
		if fmt.Sprintf("%x", root) != "76bf715a5208daa07aabcd414d6759ad6e55254617909235730c17089153e2bc0000000000000000000000000000000000000000000000000000000bfe801000" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
		expected := []byte{}
		expected = append(expected, node2...)
		expected = digest.Keccak256FlyHash(append(expected, node3...))
		if !bytes.Equal(root, expected) {
			t.Errorf(fmt.Sprintf("%x", expected))
		}
	})
	t.Run("#GetRoot of leaf 3", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(3)
		node6, _ := memoryBasedDb1.Get(6)
		if fmt.Sprintf("%x", root) != "c1371662e5123efcdf1d1fa786d040b8434df32e9af9e1e25c34dcde5332f14b0000000000000000000000000000000000000000000000000000000ffd003ffe" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
		if !bytes.Equal(root, digest.Keccak256FlyHash(node6)) {
			t.Errorf(fmt.Sprintf("%x", digest.Keccak256FlyHash(node6)))
		}
	})

	t.Run("#GetRoot of leaf 59", func(t *testing.T) {
		root := memoryBasedMmr1.GetRoot(59)
		node62, _ := memoryBasedDb1.Get(62)
		node93, _ := memoryBasedDb1.Get(93)
		node108, _ := memoryBasedDb1.Get(108)
		node115, _ := memoryBasedDb1.Get(115)
		if fmt.Sprintf("%x", root) != "5847b83a713a3da7cadf901093768652f0eb9b2fb058e67961434ec5b7bf34ae000000000000000000000000000000000000000000000000000000f1fba4e525" {
			t.Errorf("root of leaf index 999 was incorrect, got:\n %s", fmt.Sprintf("%x", root))
		}
		expected := []byte{}
		expected = append(expected, node62...)
		expected = append(expected, node93...)
		expected = append(expected, node108...)
		expected = append(expected, node115...)
		expected = digest.Keccak256FlyHash(expected)
		if !bytes.Equal(root, expected) {
			t.Errorf(fmt.Sprintf("%x", expected))
		}
	})

	t.Run("#Delete from a fileBasedMmr", func(t *testing.T) {
		leafLength := tempFileBasedDb.GetLeafLength()
		_, has34 := tempFileBasedMmr.GetUnverified(34)
		if leafLength != 1000 || !has34 {
			t.Errorf("ERRRRR")
		}
		tempFileBasedMmr.Delete(34)
		leafLength = tempFileBasedDb.GetLeafLength()

		_, has32 := tempFileBasedMmr.GetUnverified(32)
		_, has33 := tempFileBasedMmr.GetUnverified(33)

		_, has34 = tempFileBasedMmr.GetUnverified(34)
		_, has35 := tempFileBasedMmr.GetUnverified(35)
		if !has32 || !has33 || leafLength != 34 || has34 || has35 {
			t.Errorf("deleting '34' should remove everything after leaf 33")
		}
		tempFileBasedMmr.Append(sampleLeaf1, 34) // shouldnt throw an error
		leaf34, has34 := tempFileBasedMmr.GetUnverified(34)

		if !has34 || !bytes.Equal(leaf34, sampleLeaf1) {
			t.Errorf("leaf34 should have gotten rewritten as sampleLeaf1")
		}
	})

	t.Run("#Delete everything after leaf index 33", func(t *testing.T) {
		leafLength := memoryBasedDb1.GetLeafLength()
		_, has34 := memoryBasedMmr1.GetUnverified(34)
		if leafLength != 1000 || !has34 {
			t.Errorf("ERRRRR")
		}
		memoryBasedMmr1.Delete(34)
		leafLength = memoryBasedDb1.GetLeafLength()

		_, has32 := memoryBasedMmr1.GetUnverified(32)
		_, has33 := memoryBasedMmr1.GetUnverified(33)
		// memoryBasedMmr1.Get(34) // should throw
		// memoryBasedMmr1.Get(35) // should throw

		_, has34 = memoryBasedMmr1.GetUnverified(34)
		_, has35 := memoryBasedMmr1.GetUnverified(35)
		if !has32 || !has33 || leafLength != 34 || has34 || has35 {
			t.Errorf("deleting '34' should remove everything after leaf 33")
		}
	})

	t.Run("#GetProof", func(t *testing.T) {
		proofMmr := memoryBasedMmr1.GetProof([]int64{18})
		serialized := proofMmr.Serialize()
		fmt.Print(fmt.Sprintf("\n\nProofhex:\n0x%x\n\n", serialized))
		computedPrfMmr := FromSerialized(digest.Keccak256FlyHash, serialized) // testing it doesnt throw
		_, has18 := computedPrfMmr.GetUnverified(18)
		_, has19 := computedPrfMmr.GetUnverified(19)
		_, has17 := computedPrfMmr.GetUnverified(17)
		_, has20 := computedPrfMmr.GetUnverified(20)
		computedPrfMmr.Get(18) // expect not to throw
		computedPrfMmr.Get(19) // expect not to throw
		// computedPrfMmr.Get(17) // expect to throw
		// computedPrfMmr.Get(20) // expect to throw
		if !has18 || !has19 || has17 || has20 {
			t.Errorf("Proof contained wrong subset of valued")
		}
	})

	os.Remove("../db/testdata/temp.mmr")
}
