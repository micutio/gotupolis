package tuplespace

import (
	"fmt"

	"github.com/tidwall/btree"
)

// The Store defines an interface that any concrete implementation of a tuple space has to follow.
// The tuple space assumes the store implementation to be thread-safe in order to allow concurrent
// access.
type Store interface {
	// In returns a tuple that matches the argument and remove it from the space.
	In(query Tuple) Result

	// Read returns a tuple that matches the argument.
	Read(query Tuple) Result

	// Out writes a tuple into the tuple space.
	Out(tuple Tuple) bool
}

// The BTreeStore is a simple in-memory implementation of a store.
type BTreeStore struct {
	tree *btree.BTreeG[Tuple]
}

// NewSimpleStore creates an empty store instance which is ready for use.
func NewSimpleStore() *BTreeStore {
	return &BTreeStore{tree: btree.NewBTreeG(TupleOrder)}
}

// In implements the `In` function of the `Store` interface.
func (store *BTreeStore) In(query Tuple) Result {
	tuple, found := store.tree.Get(query)
	if found {
		if tuple.IsMatching(query) {
			store.tree.Delete(tuple)
			return Result{tuple, true}
		}

		fmt.Printf("[In] tuple %v does not match query %v\n", tuple, query)
	}
	return Result{Tuple{}, false}
}

// Read implements the `Read` function of the `Store` interface.
func (store *BTreeStore) Read(query Tuple) Result {
	tuple, found := store.tree.Get(query)
	if found {
		if tuple.IsMatching(query) {
			return Result{tuple, true}
		}

		fmt.Printf("[Read] tuple %v does not match query %v\n", tuple, query)
	}
	return Result{Tuple{}, false}
}

// Out implements the `Out` function of the `Store` interface
// Returns `true` if the tuple was inserted, false otherwise.
func (store *BTreeStore) Out(tuple Tuple) bool {
	if !tuple.IsDefined() {
		fmt.Printf("[Out] Warning: attempt to store undefined tuple %v \n", tuple)
		return false
	}

	store.tree.Set(tuple)
	return true
}
