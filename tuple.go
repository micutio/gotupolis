package gotupolis

import (
	"fmt"
	"math"
	"reflect"
)

const (
	// INT indicates 32bit-integers.
	INT uint = 1
	// FLOAT indicates double precision (64bit) floating point numbers.
	FLOAT uint = 2
	// STRING indicates... well... strings.
	STRING uint = 3
	// TUPLE indicates a nested tuple.
	TUPLE uint = 4
	// WILDCARD indicates any possible type of the above.
	WILDCARD uint = 5
	// NONE indicates an invalid type
	NONE uint = 0

	// FLOATPRECISION sets the error for floating point comparison
	FLOATPRECISION float64 = 0.0000001
)

// Elem acts as an element container, holding a generic element and its type indication.
type Elem struct {
	elemType  uint
	elemValue interface{}
}

// I instantiates an int-type tuple element.
func I(intVal int32) Elem {
	return Elem{INT, intVal}
}

// F instantiates a double precision (64bit) float64-type tuple element.
func F(floatVal float64) Elem {
	return Elem{FLOAT, floatVal}
}

// S instantiates a string-type tuple element.
func S(stringVal string) Elem {
	return Elem{STRING, stringVal}
}

// T instantiates a Tuple-type tuple element.
func T(tupleVal Tuple) Elem {
	return Elem{TUPLE, tupleVal}
}

// A instantiates a Wildcard tuple element.
func A() Elem {
	return Elem{WILDCARD, nil}
}

// Match two elements for equality, which is true either if they are of the same type and value
// or one or both are wildcards.
func (e Elem) isMatching(other Elem) bool {
	if e.elemType == 1 && other.elemType == 1 {
		return e.elemValue.(int32) == other.elemValue.(int32)
	}

	if e.elemType == 2 && other.elemType == 2 {
		return (math.Abs(e.elemValue.(float64)-other.elemValue.(float64)) < FLOATPRECISION)
	}

	if e.elemType == 3 && other.elemType == 3 {
		return e.elemValue.(string) == other.elemValue.(string)
	}
	if e.elemType == 4 && other.elemType == 4 {
		return e.elemValue.(Tuple).IsMatching(other.elemValue.(Tuple))
	}

	if e.elemType == 0 || other.elemType == 0 {
		return false
	}

	if e.elemType == 5 || other.elemType == 5 {
		return true
	}
	return false
}

// Tuple can contain elements of five different data types:
// - integers
// - floating point numbers
// - strings
// - tuples themselves
// - wildcards
type Tuple struct {
	elements []Elem
}

// IsMatching checks two tuples for equality, which is true if
// - they are of the same lenght AND
// - each element of one matches the others
func (t Tuple) IsMatching(other Tuple) bool {
	tSize := len(t.elements)
	otherSize := len(other.elements)

	// Check length of tuples first
	if tSize != otherSize {
		return false
	}

	// Check each element for equality.
	for i := 0; i < tSize; i++ {
		if !t.elements[i].isMatching(other.elements[i]) {
			return false
		}
	}

	return true
}

// MakeTuple creates a new Tuple instance from the given parameters.
func MakeTuple(element ...Elem) Tuple {
	var resultTuple Tuple
	for _, e := range element {
		fmt.Printf("element %v of type %T (reflect %v)\n", e, e, reflect.TypeOf(e))
	}
	resultTuple.elements = element
	fmt.Printf("resulting tuple: %v", resultTuple)
	return resultTuple
}
