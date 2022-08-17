package gotupolis

// The Store defines an interface that any concrete implementation of a tuple space has to
// follow.
type Store interface {

	// Read a tuple that matches the argument and remove it from the space.
	In(tuple Tuple) Tuple

	// Read a tuple that matches the argument.
	Read(tuple Tuple) Tuple

	// Write a tuple into the tuple space.
	Out(tuple Tuple)
}
