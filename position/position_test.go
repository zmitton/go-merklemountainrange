package position

import (
	"reflect"
	"testing"
)

func TestGetNodePosition(t *testing.T) {
	t.Run("#GetNodePosition of a single very large leaf index", func(t *testing.T) {
		li := 1234567890
		pi := 2469135768
		if GetNodePosition(int64(li)).Index != int64(pi) {
			t.Errorf("leaf index %d should have position index %d", li, pi)
		}
	})
	t.Run("#GetNodePosition of many fundamental leaves", func(t *testing.T) {
		pis := []int64{0, 1, 3, 4, 7, 8, 10, 11, 15, 16, 18, 19, 22, 23, 25, 26,
			31, 32, 34, 35, 38, 39, 41, 42, 46, 47, 49, 50, 53, 54, 56, 57, 63, 64}
		for i := 0; i <= 33; i++ {
			if GetNodePosition(int64(i)).Index != pis[i] {
				t.Errorf("leaf index %d should have position index %d", int64(i), pis[i])
			}
		}
	})
}
func TestLeftChild(t *testing.T) {
	t.Run("#LeftChild", func(t *testing.T) {
		p := Position{2147483645, 29, true}
		lcp := Position{1610612733, 28, false}
		if LeftChild(p) != lcp {
			t.Errorf("left child should be %v, got %v", lcp, LeftChild(p))
		}
	})
	t.Run("#LeftChild", func(t *testing.T) {
		ps := []Position{{2, 1, false}, {5, 1, true}, {6, 2, false}, {9, 1, false},
			{12, 1, true}, {13, 2, true}, {14, 3, false}, {17, 1, false}, {20, 1, true},
			{21, 2, false}, {24, 1, false}, {27, 1, true}, {28, 2, true}, {29, 3, true},
			{30, 4, false}, {33, 1, false}, {36, 1, true}, {37, 2, false}}
		lcps := []Position{{0, 0, false}, {3, 0, false}, {2, 1, false}, {7, 0, false},
			{10, 0, false}, {9, 1, false}, {6, 2, false}, {15, 0, false}, {18, 0, false},
			{17, 1, false}, {22, 0, false}, {25, 0, false}, {24, 1, false}, {21, 2, false},
			{14, 3, false}, {31, 0, false}, {34, 0, false}, {33, 1, false}}
		for i := 0; i < len(lcps); i++ {
			lcp := LeftChild(ps[i])
			if lcp != lcps[i] {
				t.Errorf("left child of %v should be  %v, got %v", ps[i], lcps[i], lcp)
			}
		}
	})
}
func TestRightChild(t *testing.T) {
	t.Run("#RightChild", func(t *testing.T) {
		p := Position{2147483645, 29, true}
		lcp := Position{2147483644, 28, true}
		if RightChild(p) != lcp {
			t.Errorf("right child should be %v, got %v", lcp, RightChild(p))
		}
	})
	t.Run("#RightChild", func(t *testing.T) {
		ps := []Position{{2, 1, false}, {5, 1, true}, {6, 2, false}, {9, 1, false},
			{12, 1, true}, {13, 2, true}, {14, 3, false}, {17, 1, false}, {20, 1, true},
			{21, 2, false}, {24, 1, false}, {27, 1, true}, {28, 2, true}, {29, 3, true},
			{30, 4, false}, {33, 1, false}, {36, 1, true}, {37, 2, false}}
		rcps := []Position{{1, 0, true}, {4, 0, true}, {5, 1, true}, {8, 0, true},
			{11, 0, true}, {12, 1, true}, {13, 2, true}, {16, 0, true}, {19, 0, true},
			{20, 1, true}, {23, 0, true}, {26, 0, true}, {27, 1, true}, {28, 2, true},
			{29, 3, true}, {32, 0, true}, {35, 0, true}, {36, 1, true}}
		for i := 0; i < len(rcps); i++ {
			rcp := RightChild(ps[i])
			if rcp != rcps[i] {
				t.Errorf("right child of %v should be  %v, got %v", ps[i], rcps[i], rcp)
			}
		}
	})
}
func TestSibling(t *testing.T) {
	t.Run("#Sibling", func(t *testing.T) {
		ps := []Position{{62, 5, false}, {61, 4, true}, {5, 1, true},
			{65, 1, false}, {44, 2, true}}
		sps := []Position{{125, 5, true}, {30, 4, false}, {2, 1, false},
			{68, 1, true}, {37, 2, false}}
		for i := 0; i < len(sps); i++ {
			sp := Sibling(ps[i])
			if sp != sps[i] {
				t.Errorf("sibling of %v should be %v, got %v", ps[i], sps[i], sp)
			}
		}
	})
}
func TestParentIndex(t *testing.T) {
	t.Run("#ParentIndex of single large node index", func(t *testing.T) {
		p := Position{2147483645, 29, true}
		c := Position{2147483644, 28, true}
		if ParentIndex(c) != p.Index {
			t.Errorf("parent index of %v should be %d, got %d", c, p.Index, ParentIndex(c))
		}
	})
	t.Run("#ParentIndex of many fundamental positions", func(t *testing.T) {
		ps := []Position{{62, 5, false}, {61, 4, true}, {5, 1, true}, {65, 1, false}, {44, 2, true}}
		pis := []int64{126, 62, 6, 69, 45}
		for i := 0; i < len(pis); i++ {
			pi := ParentIndex(ps[i])
			if pi != pis[i] {
				t.Errorf("ParentIndex of %v should be %v, got %v", ps[i], pis[i], pi)
			}
		}
	})
}
func TestPeakPositions(t *testing.T) {
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(0)
		computedPeaks := PeakPositions(leafIndex)
		expectedPeaks := []Position{{0, 0, false}}
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(1)
		computedPeaks := PeakPositions(leafIndex)
		expectedPeaks := []Position{{2, 1, false}}
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(2)
		expectedPeaks := []Position{{2, 1, false}, {3, 0, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})

	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(9)
		expectedPeaks := []Position{{14, 3, false}, {17, 1, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(30)
		expectedPeaks := []Position{{30, 4, false}, {45, 3, false}, {52, 2, false}, {55, 1, false}, {56, 0, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(31)
		expectedPeaks := []Position{{62, 5, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(1073741823)
		expectedPeaks := []Position{{2147483646, 30, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(1073741822)
		expectedPeaks := []Position{
			{1073741822, 29, false}, {1610612733, 28, false}, {1879048188, 27, false},
			{2013265915, 26, false}, {2080374778, 25, false}, {2113929209, 24, false},
			{2130706424, 23, false}, {2139095031, 22, false}, {2143289334, 21, false},
			{2145386485, 20, false}, {2146435060, 19, false}, {2146959347, 18, false},
			{2147221490, 17, false}, {2147352561, 16, false}, {2147418096, 15, false},
			{2147450863, 14, false}, {2147467246, 13, false}, {2147475437, 12, false},
			{2147479532, 11, false}, {2147481579, 10, false}, {2147482602, 9, false},
			{2147483113, 8, false}, {2147483368, 7, false}, {2147483495, 6, false},
			{2147483558, 5, false}, {2147483589, 4, false}, {2147483604, 3, false},
			{2147483611, 2, false}, {2147483614, 1, false}, {2147483615, 0, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})
	t.Run("#PeakPosition of a single leaf index", func(t *testing.T) {
		leafIndex := int64(1077422)
		expectedPeaks := []Position{
			{2097150, 20, false}, {2129917, 14, false}, {2146300, 13, false},
			{2154491, 12, false}, {2154746, 7, false}, {2154809, 5, false},
			{2154824, 3, false}, {2154831, 2, false}, {2154834, 1, false},
			{2154835, 0, false}}
		computedPeaks := PeakPositions(leafIndex)
		if !reflect.DeepEqual(computedPeaks, expectedPeaks) {
			t.Errorf("peaks of %d should be %v, got %v", leafIndex, expectedPeaks, computedPeaks)
		}
	})

}

func TestLocalPeakPosition(t *testing.T) {
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{0, 0}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{0, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{1, 0}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{2, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{2, 2}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{3, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{9, 6}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{17, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{30, 14}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{56, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{30, 31}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{56, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{31, 30}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{62, 5, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{32, 30}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{63, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{31, 55}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{62, 5, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{32, 32}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{63, 0, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 33}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{65, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 34}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{65, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 35}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{65, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 36}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{69, 2, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 45}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{77, 3, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 47}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{77, 3, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 49}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{93, 4, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{33, 69}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{126, 6, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
	t.Run("#LocalPeakPosition of a single leaf index", func(t *testing.T) {
		leafIndexAndLength := [2]int64{1077422, 1090012}
		computedPeak := LocalPeakPosition(leafIndexAndLength[0], leafIndexAndLength[1])
		expectedPeak := Position{2162685, 15, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("peaks of leaf index %d length %d is wrong, got %v", leafIndexAndLength[0], leafIndexAndLength[1], computedPeak)
		}
	})
}

func TestGodPeakFromLeafIndex(t *testing.T) {
	t.Run("#GodPeakFromLeafIndex of a single leaf index", func(t *testing.T) {
		leafIndex := int64(0)
		computedPeak := GodPeakFromLeafIndex(leafIndex)
		expectedPeak := Position{2, 1, false}
		if !reflect.DeepEqual(computedPeak, expectedPeak) {
			t.Errorf("god peak of leaf index is %d length %v is wrong, got %v", leafIndex, expectedPeak, computedPeak)
		}
	})
	t.Run("#GodPeakFromLeafIndex of first 31 leaf indexes", func(t *testing.T) {
		expectedPeakArr := []Position{
			{2, 1, false}, {6, 2, false}, {6, 2, false}, {14, 3, false}, {14, 3, false},
			{14, 3, false}, {14, 3, false}, {30, 4, false}, {30, 4, false}, {30, 4, false},
			{30, 4, false}, {30, 4, false}, {30, 4, false}, {30, 4, false}, {30, 4, false},
			{62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false},
			{62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false},
			{62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false}, {62, 5, false},
			{62, 5, false}, {126, 6, false}}
		for i := 0; i < len(expectedPeakArr); i++ {
			if !reflect.DeepEqual(GodPeakFromLeafIndex(int64(i)), expectedPeakArr[i]) {
				t.Errorf("god peak of leaf index is %d length %v is wrong, got %v", i, expectedPeakArr[i], GodPeakFromLeafIndex(int64(i)))
			}
		}
	})
}

func TestMountainPositions(t *testing.T) {
	t.Run("#MountainPosiitons of a single leaf index", func(t *testing.T) {
		computed := MountainPositions(Position{0, 0, false}, 0)
		expected := [][]Position{}
		if !reflect.DeepEqual(computed, expected) {
			t.Errorf("mountain positions wrong, got %v", computed)
		}
	})
	t.Run("#MountainPosiitons of a single leaf index", func(t *testing.T) {
		computed := MountainPositions(Position{2, 1, false}, 0)
		expected := [][]Position{{{0, 0, false}, {1, 0, true}}}
		if !reflect.DeepEqual(computed, expected) {
			t.Errorf("mountain positions wrong, got %v", computed)
		}
	})

	t.Run("#MountainPosiitons of a few more positions", func(t *testing.T) {
		inputPeaks := []Position{
			{2, 1, false}, {2, 1, false}, {3, 0, false}, {6, 2, false},
			{6, 2, false}, {7, 0, false}, {9, 1, false}, {9, 1, false},
			{10, 0, false}, {14, 3, false}, {14, 3, false}, {14, 3, false},
			{14, 3, false}, {45, 3, false}, {45, 3, false}, {45, 3, false},
			{45, 3, false}, {45, 3, false}, {45, 3, false}, {45, 3, false},
			{45, 3, false}, {62, 5, false}, {62, 5, false},
		}
		inputIndexes := []int64{0, 1, 3, 0, 3, 7, 7, 8, 10, 3, 4, 7, 8, 31, 32, 34, 35, 38, 39, 41, 42, 22, 23}
		expectedPositionsArr := [][][]Position{
			{{{0, 0, false}, {1, 0, true}}},
			{{{0, 0, false}, {1, 0, true}}},
			{},
			{{{2, 1, false}, {5, 1, true}}, {{0, 0, false}, {1, 0, true}}},
			{{{2, 1, false}, {5, 1, true}}, {{3, 0, false}, {4, 0, true}}},
			{},
			{{{7, 0, false}, {8, 0, true}}},
			{{{7, 0, false}, {8, 0, true}}},
			{},
			{{{6, 2, false}, {13, 2, true}}, {{2, 1, false}, {5, 1, true}}, {{3, 0, false}, {4, 0, true}}},
			{{{6, 2, false}, {13, 2, true}}, {{2, 1, false}, {5, 1, true}}, {{3, 0, false}, {4, 0, true}}},
			{{{6, 2, false}, {13, 2, true}}, {{9, 1, false}, {12, 1, true}}, {{7, 0, false}, {8, 0, true}}},
			{{{6, 2, false}, {13, 2, true}}, {{9, 1, false}, {12, 1, true}}, {{7, 0, false}, {8, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{33, 1, false}, {36, 1, true}}, {{31, 0, false}, {32, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{33, 1, false}, {36, 1, true}}, {{31, 0, false}, {32, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{33, 1, false}, {36, 1, true}}, {{34, 0, false}, {35, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{33, 1, false}, {36, 1, true}}, {{34, 0, false}, {35, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{40, 1, false}, {43, 1, true}}, {{38, 0, false}, {39, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{40, 1, false}, {43, 1, true}}, {{38, 0, false}, {39, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{40, 1, false}, {43, 1, true}}, {{41, 0, false}, {42, 0, true}}},
			{{{37, 2, false}, {44, 2, true}}, {{40, 1, false}, {43, 1, true}}, {{41, 0, false}, {42, 0, true}}},
			{{{30, 4, false}, {61, 4, true}}, {{14, 3, false}, {29, 3, true}}, {{21, 2, false}, {28, 2, true}}, {{24, 1, false}, {27, 1, true}}, {{22, 0, false}, {23, 0, true}}},
			{{{30, 4, false}, {61, 4, true}}, {{14, 3, false}, {29, 3, true}}, {{21, 2, false}, {28, 2, true}}, {{24, 1, false}, {27, 1, true}}, {{22, 0, false}, {23, 0, true}}},
		}
		for i := 0; i < len(inputPeaks); i++ {
			computed := MountainPositions(inputPeaks[i], inputIndexes[i])
			if !reflect.DeepEqual(computed, expectedPositionsArr[i]) {
				t.Errorf("mountain posisitons of peak %v, leaf %d is wrong, got %v", inputPeaks[i], inputIndexes[i], computed)
			}
		}
	})
}
