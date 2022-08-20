package gotupolis

import (
	"github.com/tidwall/btree"

	option "github.com/micutio/goptional"
)

// The Store defines an interface that any concrete implementation of a tuple space has to
// follow.
type Store interface {

	// Read a tuple that matches the argument and remove it from the space.
	In(tuple Tuple) option.Maybe[Tuple]

	// Read a tuple that matches the argument.
	Read(tuple Tuple) option.Maybe[Tuple]

	// Write a tuple into the tuple space.
	Out(tuple Tuple)
}

type SimpleStore struct {
	tree *btree.BTreeG[Tuple]
}

func New() *SimpleStore {
	return &SimpleStore{tree: btree.NewBTreeG(TupleOrder)}
}
