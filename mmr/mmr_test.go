package mmr

import "testing"
import "fmt"
import "merklemountainrange/db"

func Test(t *testing.T) {
	db1 := db.NewFilebaseddb("../db/testdata/etcleafdata.mmr")
	mmr := NewMmr(H, db1)

	leafLength := mmr.GetLeafLength()
	if leafLength != 1000 {
		fmt.Printf("db1node0 %d \n", leafLength)
		t.Errorf("leafLength should be equal to 1000, not %d", leafLength)
	}

	if false {
		t.Errorf("ZZZZZZZ")
	}

}
