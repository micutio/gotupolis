package gotupolis

// The TupleSpace defines an interface that any concrete implementation has to follow.
type TupleSpace interface {

	// Read a tuple that matches the argument and remove it from the space.
	In(tuple *Tuple) *Tuple

	// Read a tuple that matches the argument.
	Read(tuple *Tuple) *Tuple

	// Write a tuple into the tuple space.
	Out(tuple *Tuple)
}
