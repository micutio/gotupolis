package gotupolis

import "testing"

func TestMakeTuple(t *testing.T) {
	tup1 := MakeTuple(I(1), I(2), F(3.14), S("Furz!"), T(*MakeTuple(S("hurz"))))
	tup2 := MakeTuple(I(1), A(), F(3.14), S("Furz!"), A())
	tup3 := MakeTuple(I(1), A(), F(3.140001), S("Furz!"), A())
	tup4 := MakeTuple(I(1), I(2), F(3.14), S("Furz!"), T(*MakeTuple(A())))

	if !tup1.IsMatching(tup2) {
		t.Errorf("Error: tuples %v and %v not matching", tup1, tup2)
	}

	if tup1.IsMatching(tup3) {
		t.Errorf("Error: tuples %v and %v should not match", tup1, tup3)
	}

	if !tup1.IsMatching(tup4) {
		t.Errorf("Error: tuples %v and %v not matching", tup1, tup4)
	}
}

func TestElemOrdering(t *testing.T) {
	e1 := I(-1)
	e2 := I(-1)
	e3 := I(2)
	e4 := F(3.14)

	if ord := e1.order(&e2); ord != EQ {
		t.Errorf("Error: order of tuples %v and %v should be EQ, not %v", e1, e2, ord)
	}

	if ord := e1.order(&e3); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", e1, e3, ord)
	}

	if ord := e3.order(&e2); ord != GT {
		t.Errorf("Error: order of tuples %v and %v should be GT, not %v", e3, e2, ord)
	}

	if ord := e4.order(&e1); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", e4, e1, ord)
	}

	tup1 := MakeTuple(I(1), A(), F(3.14), S("Furz!"), A())
	tup2 := MakeTuple(I(1), A(), F(3.1400), S("Furz!"), A())
	tup4 := MakeTuple(I(1), I(2), F(3.14), S("Furz!"), T(*MakeTuple(A())))

	if ord := tup1.order(tup2); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", tup1, tup2, ord)
	}
}
