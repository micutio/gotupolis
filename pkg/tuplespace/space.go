package gotupolis

import (
	option "github.com/micutio/goptional"
)

// TODO: Refer to https://github.com/sgjp/go-tuplespace for inspiration

// The Space contains the actual store and handles concurrent read and write access to it.
type Space struct {
	store Store
}

// Create a new space instance that uses the default store implementation `SimpleStore`
func NewSpace() *Space {
	return &Space{store: NewSimpleStore()}
}

// Create a new space that uses the given store implementation
func MakeSpace(store Store) *Space {
	return &Space{store}
}

// Retrieve a tuple that matches the query from the space and remove it.
// The tuple may contain wildcards. If it does and matches multiple tuples in the space, then an
// arbitrary match will be returned as a result.
func (s *Space) In(query Tuple) <-chan option.Maybe[Tuple] {
	c := make(chan option.Maybe[Tuple])
	go func() {
		c <- s.store.In(query)
	}()
	return c
}

// Retrieve a tuple that matches the query from the space but do not remove it.
// The tuple may contain wildcards. If it does and matches multiple tuples in the space, then an
// arbitrary match will be returned as a result.
func (s *Space) Read(query Tuple) <-chan option.Maybe[Tuple] {
	c := make(chan option.Maybe[Tuple])
	go func() {
		c <- s.store.Read(query)
	}()
	return c
}

// Insert a tuple into the tuple space.
// The tuple must be defined, i.e.: NOT contain any wildcards or `None`, otherwise it will not be
// inserted.
func (s *Space) Out(query Tuple) <-chan bool {
	c := make(chan bool)
	go func() {
		c <- s.store.Out(query)
	}()
	return c
}
