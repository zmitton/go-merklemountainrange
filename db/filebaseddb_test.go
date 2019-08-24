package db

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	db1 := OpenFilebaseddb("../db/testdata/etcleafdata.mmr")
	db2 := CreateFilebaseddb("../db/testdata/temp.mmr", 64)

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

	db2.Set(db1node0, 0)
	db2node0, _ := db2.Get(0)
	if db2node0[0:1][0] != []byte{212}[0] {
		t.Errorf("index 0 should start with 'd4'")
	}

	db2.Set(db1node1, 1)
	db2node1, _ := db2.Get(1)
	if db2node1[0:1][0] != []byte{136}[0] {
		t.Errorf("index 1 should start with '88'")
	}

	db2.Set(db1node2, 2)
	db2node2, _ := db2.Get(2)
	if db2node2[0:1][0] != []byte{135}[0] {
		t.Errorf("index 2 should start with '87'")
	}
}
