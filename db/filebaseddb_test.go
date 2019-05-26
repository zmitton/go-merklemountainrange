package db

import (
	"fmt"
	"testing"
)

// import "github.com/ethereum/go-ethereum/common/math"

func Test(t *testing.T) {
	db1 := NewFilebaseddb("etcLeafData.mmr")
	db2 := NewFilebaseddb("temp.mmr")

	leafLength := db1.GetLeafLength()
	if leafLength != 1000 {
		fmt.Printf("db1node0 %d \n", leafLength)
		t.Errorf("leafLength should be equal to 1000, not %d", leafLength)
	}

	db1node0, _ := db1.Get(0)
	// fmt.Printf("db1node0 %08b \n", db1node0[0:1][0])
	// fmt.Printf("db1node0 %08b \n", []byte{212}[0])
	if db1node0[0:1][0] != []byte{212}[0] {
		t.Errorf("index 0 should start with 'd4'")
	}

	db1node1, _ := db1.Get(1)
	if db1node1[0:1][0] != []byte{136}[0] {
		t.Errorf("index 1 should start with '88'")
	}

	db1node2, _ := db1.Get(2)
	if db1node2[0:1][0] != []byte{135}[0] {
		t.Errorf("index 2 should start with '87'")
	}

	db2.Set(0, db1node0)
	db2node0, _ := db2.Get(0)
	if db2node0[0:1][0] != []byte{212}[0] {
		t.Errorf("index 0 should start with 'd4'")
	}

	db2.Set(1, db1node1)
	db2node1, _ := db2.Get(1)
	if db2node1[0:1][0] != []byte{136}[0] {
		t.Errorf("index 1 should start with '88'")
	}

	db2.Set(2, db1node2)
	db2node2, _ := db2.Get(2)
	if db2node2[0:1][0] != []byte{135}[0] {
		t.Errorf("index 2 should start with '87'")
	}

}
