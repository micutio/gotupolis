package gotupolis

import (
	"testing"

	option "github.com/micutio/goptional"
)

func TestStore(t *testing.T) {
	store := getStoreImpl()
	tup := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))
	query := MakeTuple(I(1), Any(), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))

	// insert for testing defined queries
	store.Out(tup)

	var tupleOpt1a option.Maybe[Tuple] = store.Read(tup)
	if !tupleOpt1a.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", tup)
	}

	// since we've only read the tuple, it should still be in the store
	var tupleOpt1b option.Maybe[Tuple] = store.In(tup)
	if !tupleOpt1b.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", tup)
	}

	// now the tuple should be gone
	var tupleOpt1c option.Maybe[Tuple] = store.In(tup)
	if tupleOpt1c.IsPresent() {
		t.Errorf("Error: store should not contain tuple %v anymore", tup)
	}

	// insert for testing with not defined queries
	store.Out(tup)

	var tupleOpt2a option.Maybe[Tuple] = store.Read(query)
	if !tupleOpt2a.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", query)
	}

	// since we've only read the tuple, it should still be in the store
	var tupleOpt2b option.Maybe[Tuple] = store.In(query)
	if !tupleOpt2b.IsPresent() {
		t.Errorf("Error: cannot find tuple %v", query)
	}

	// now the tuple should be gone
	var tupleOpt2c option.Maybe[Tuple] = store.In(query)
	if tupleOpt2c.IsPresent() {
		t.Errorf("Error: store should not contain tuple %v anymore", tup)
	}

}

func getStoreImpl() Store {
	return NewSimpleStore()
}
