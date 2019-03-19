package gotupolis

import "testing"

func TestMakeTuple(t *testing.T) {
	tup1 := MakeTuple(I(1), I(2), F(3.14), S("Furz!"), T(MakeTuple(S("hurz"))))
	tup2 := MakeTuple(I(1), A(), F(3.14), S("Furz!"), A())
	tup3 := MakeTuple(I(1), A(), F(3.140001), S("Furz!"), A())
	tup4 := MakeTuple(I(1), I(2), F(3.14), S("Furz!"), T(MakeTuple(A())))

	if !tup1.IsMatching(&tup2) {
		t.Errorf("Error: tuples %v and %v not matching", tup1, tup2)
	}

	if tup1.IsMatching(&tup3) {
		t.Errorf("Error: tuples %v and %v should not match", tup1, tup3)
	}

	if !tup1.IsMatching(&tup4) {
		t.Errorf("Error: tuples %v and %v not matching", tup1, tup4)
	}
}
