package gotupolis

import "testing"

// More tests to do:
// - (Elem) String()
// - (Tuple) String()

func TestMakeTuple(t *testing.T) {
	tup1 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))
	tup2 := MakeTuple(I(1), A(), F(3.14), S("Foo!"), A())
	tup3 := MakeTuple(I(1), A(), F(3.140001), S("Foo!"), A())
	tup4 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(A())))

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

	tup1 := MakeTuple(I(1), A(), F(3.14), S("Foo!"), A())
	tup2 := MakeTuple(I(1), A(), F(3.1400), S("Foo!"), A())
	tup4 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(A())))
	tup5 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(I(5), S("Bar!"))))
	tup6 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(S("Bar!"))))

	if ord := tup1.order(tup2); ord != EQ {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", tup1, tup2, ord)
	}

	if ord := tup1.order(tup4); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be EQ, not %v", tup1, tup4, ord)
	}

	if ord := tup4.order(tup5); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be GT, not %v", tup4, tup5, ord)
	}

	if ord := tup5.order(tup6); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", tup5, tup6, ord)
	}

	if ord := tup4.order(tup6); ord != LT {
		t.Errorf("Error: order of tuples %v and %v should be LT, not %v", tup4, tup6, ord)
	}
}

func TestToString(t *testing.T) {
	// Test element to string conversions
	// For numbers: run positive, negative and zero
	e1 := I(-3)
	e2 := I(+45)
	e3 := I(0)
	e4 := F(-.14)
	e5 := F(3.14)
	e6 := F(0.0)
	e7 := S("真棒！")

	if e1.String() != "-3" {
		t.Errorf("Error: element %v should equate to string `-3`", e1.String())
	}

	if e2.String() != "45" {
		t.Errorf("Error: element %v should equate to string `45`", e2.String())
	}

	if e3.String() != "0" {
		t.Errorf("Error: element %v should equate to string `0`", e3.String())
	}

	if e4.String() != "-0.14" {
		t.Errorf("Error: element %v should equate to string `-0.14`", e4.String())
	}

	if e5.String() != "3.14" {
		t.Errorf("Error: element %v should equate to string `3.14`", e5.String())
	}

	if e6.String() != "0" {
		t.Errorf("Error: element %v should equate to string `0`", e6.String())
	}

	if e7.String() != "\"真棒！\"" {
		t.Errorf("Error: element %v should equate to string \"真棒！\"", e7.String())
	}

	// Test tuple to string conversions
	tup1 := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(I(5), S("Bar!"))))

	if tup1.String() != "(1|2|3.14|\"Foo!\"|(5|\"Bar!\"))" {
		t.Errorf("Error: element %v should equate to string \"(1|2|3.14|\"Foo!\"|(5|\"Bar!\"))\"", e7.String())
	}
}
