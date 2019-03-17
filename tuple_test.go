package gotupolis

import "testing"

func TestMakeTuple(t *testing.T) {
	tup1 := MakeTuple(INTe(1), INTe(2), FLOATe(3.14), STRINGe("Furz!"), TUPLEe(MakeTuple(STRINGe("hurz"))))
	tup2 := MakeTuple(INTe(1), ANYe(), FLOATe(3.14), STRINGe("Furz!"), ANYe())

	if !tup1.IsMatching(tup2) {
		t.Errorf("Error: tuples %v and %v not matching", tup1, tup2)
	}
}
