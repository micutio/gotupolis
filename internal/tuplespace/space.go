package tuplespace

// TODO: Refer to https://github.com/sgjp/go-tuplespace for inspiration

// The Space contains the actual store and handles concurrent read and write access to it.
type Space struct {
	store Store
}

type Result struct {
	tuple Tuple
	OK bool
}

// NewSpace creates a new space instance that uses the default store implementation `SimpleStore`.
func NewSpace() *Space {
	return &Space{store: NewSimpleStore()}
}

// MakeSpace creates a new space that uses the given store implementation.
func MakeSpace(store Store) *Space {
	return &Space{store}
}

// In retrieves a tuple that matches the query from the space and removes it.
// The tuple may contain wildcards. If it does and matches multiple tuples in the space, then an
// arbitrary match will be returned as a result.
func (s *Space) In(query Tuple) <-chan Result {
	c := make(chan Result)
	go func() {
		c <- s.store.In(query)
	}()
	return c
}

// Read retrieves a tuple that matches the query from the space but do not remove it.
// The tuple may contain wildcards. If it does and matches multiple tuples in the space, then an
// arbitrary match will be returned as a result.
func (s *Space) Read(query Tuple) <-chan Result {
	c := make(chan Result)
	go func() {
		c <- s.store.Read(query)
	}()
	return c
}

// Out inserts a tuple into the tuple space.
// The tuple must be defined, i.e.: NOT contain any wildcards or `None`, otherwise it will not be
// inserted.
func (s *Space) Out(query Tuple) <-chan bool {
	c := make(chan bool)
	go func() {
		c <- s.store.Out(query)
	}()
	return c
}
