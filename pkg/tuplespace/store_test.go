package gotupolis

import (
	"testing"
)

func TestStore(t *testing.T) {
	store := getStoreImpl()
	tup := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))
	query := MakeTuple(I(1), Any(), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))

	// insert for testing defined queries
	isTupInserted := store.Out(tup)
	if !isTupInserted {
		t.Errorf("Error: tuple %v was not inserted, but it should have been", tup)
	}

	isQueryInserted := store.Out(query)
	if isQueryInserted {
		t.Errorf("Error: query %v was inserted, but should not have been, since it contains "+
			"wildcards", query)
	}

	tupleOpt1a := store.Read(tup)
	if !tupleOpt1a.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", tup)
	}

	// since we've only read the tuple, it should still be in the store
	tupleOpt1b := store.In(tup)
	if !tupleOpt1b.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", tup)
	}

	// now the tuple should be gone
	tupleOpt1c := store.In(tup)
	if tupleOpt1c.IsPresent() {
		t.Errorf("Error: store should not contain tuple %v anymore", tup)
	}

	// insert for testing with not defined queries
	store.Out(tup)

	tupleOpt2a := store.Read(query)
	if !tupleOpt2a.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", query)
	}

	// since we've only read the tuple, it should still be in the store
	tupleOpt2b := store.In(query)
	if !tupleOpt2b.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", query)
	}

	// now the tuple should be gone
	tupleOpt2c := store.In(query)
	if tupleOpt2c.IsPresent() {
		t.Errorf("Error: store should not contain tuple %v anymore", tup)
	}

}

func getStoreImpl() Store {
	return NewSimpleStore()
}
