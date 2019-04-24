package gotupolis

import (
	"sync"
)

// TupleSpaceTree is a binary-tree-based implementation of a tuple space.
// This implementation uses a binary tree as underlying data structure, which helps
// with search performance, but is still quite naïve because if a lack of automated balancing.
// This implementation is not meant for use in production but merely as a proof of concept.
// Shoutout to
//			https://flaviocopes.com/golang-data-structure-binary-search-tree/
// for some inspiration.
type TupleSpaceTree struct {
	root *tstNode
	lock sync.RWMutex
}

type tstNode struct {
	leftChild  *tstNode
	rightChild *tstNode
	value      Tuple
}

// In reads a tuple that matches the argument and remove it from the space.
func (TupleSpaceTree) In(tuple Tuple) Tuple {

}

// Read a tuple that matches the argument.
func (TupleSpaceTree) Read(tuple Tuple) Tuple {

}

// Out write a tuple into the tuple space.
func (TupleSpaceTree) Out(tuple Tuple) {

}
