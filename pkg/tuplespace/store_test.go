package gotupolis

import (
	"testing"

	option "github.com/micutio/goptional"
)

func TestStore(t *testing.T) {
	store := getStoreImpl()
	tup := MakeTuple(I(1), I(2), F(3.14), S("Foo!"), T(MakeTuple(S("hurz"))))

	// test insertion
	store.Out(tup)

	// test read
	var tupleOpt1 option.Maybe[Tuple] = store.Read(tup)
	if !tupleOpt1.IsPresent() {
		t.Errorf("Error: cannot find inserted tuple %v", tup)
	}

	// test out
	// since we've only read the tuple, it should still be in the store
	var tupleOpt2 option.Maybe[Tuple] = store.In(tup)
	if !tupleOpt2.IsPresent() {
		t.Errorf("Error: cannot find inserted tuple %v", tup)
	}

	// now the tuple should be gone
	var tupleOpt3 option.Maybe[Tuple] = store.In(tup)
	if tupleOpt3.IsPresent() {
		t.Errorf("Error: store should not contain tuple %v anymore", tup)
	}
}

func getStoreImpl() Store {
	return NewSimpleStore()
}
