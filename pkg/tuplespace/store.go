package gotupolis

import (
	"fmt"

	"github.com/tidwall/btree"

	option "github.com/micutio/goptional"
)

// The Store defines an interface that any concrete implementation of a tuple space has to
// follow.
type Store interface {

	// Read a tuple that matches the argument and remove it from the space.
	In(query Tuple) option.Maybe[Tuple]

	// Read a tuple that matches the argument.
	Read(query Tuple) option.Maybe[Tuple]

	// Write a tuple into the tuple space.
	Out(tuple Tuple)
}

type SimpleStore struct {
	tree *btree.BTreeG[Tuple]
}

func NewSimpleStore() *SimpleStore {
	return &SimpleStore{tree: btree.NewBTreeG(TupleOrder)}
}

func (store *SimpleStore) In(query Tuple) option.Maybe[Tuple] {
	tuple, found := store.tree.Get(query)
	if found {
		if tuple.IsMatching(query) {
			store.tree.Delete(tuple)
			return option.NewJust(tuple)
		} else {
			fmt.Printf("[In] tuple %v does not match query %v\n", tuple, query)
		}
	}
	return option.NewNothing[Tuple]()
}

func (store *SimpleStore) Read(query Tuple) option.Maybe[Tuple] {
	tuple, found := store.tree.Get(query)
	if found {
		if tuple.IsMatching(query) {
			return option.NewJust(tuple)
		} else {
			fmt.Printf("[Read] tuple %v does not match query %v\n", tuple, query)
		}
	}
	return option.NewNothing[Tuple]()
}

func (store *SimpleStore) Out(tuple Tuple) {
	if !tuple.IsDefined() {
		fmt.Printf("[Out] Warning: attempt to store undefined tuple %v \n", tuple)
	} else {
		store.tree.Set(tuple)
	}
}
