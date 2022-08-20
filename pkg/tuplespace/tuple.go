package gotupolis

import (
	"fmt"
	"math"
	"strings"
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
	// ANY indicates any possible type of the above, functioning as a wildcard.
	ANY uint = 5
	// NONE indicates an invalid type
	NONE uint = 0

	// FLOATPRECISION sets the error for floating point comparison
	FLOATPRECISION float64 = 0.0000001

	// LT is the `less than` return value for order comparisons
	LT int = -1
	// EQ is the `equals` return value for order comparisons
	EQ int = 0
	// GT is the `greater than` return value for order comparisons
	GT int = 1
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Elem acts as an element container, holding a generic element and its type indication.
type Elem struct {
	elemType  uint
	elemValue interface{}
}

func (e Elem) String() string {
	switch e.elemType {
	case INT:
		return fmt.Sprintf("%v", e.elemValue.(int32))
	case FLOAT:
		return fmt.Sprintf("%v", e.elemValue.(float64))
	case STRING:
		return fmt.Sprintf("\"%v\"", e.elemValue.(string))
	case TUPLE:
		return e.elemValue.(Tuple).String()
	case ANY:
		return "_"
	case NONE:
		return "nil"
	default:
		return fmt.Sprintf("Error: invalid elem type %T", e.elemValue)
	}
}

// Tuple element constructors /////////////////////////////////////////////////

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
func Any() Elem {
	return Elem{ANY, nil}
}

// Match two elements for equality, which is true either if they are of the same type and value
// or one or both are wildcards.
func (e Elem) isMatching(other Elem) bool {
	if e.elemType == INT && other.elemType == INT {
		return e.elemValue.(int32) == other.elemValue.(int32)
	}

	if e.elemType == FLOAT && other.elemType == FLOAT {
		return (math.Abs(e.elemValue.(float64)-other.elemValue.(float64)) < FLOATPRECISION)
	}

	if e.elemType == STRING && other.elemType == STRING {
		return e.elemValue.(string) == other.elemValue.(string)
	}
	if e.elemType == TUPLE && other.elemType == TUPLE {
		return e.elemValue.(Tuple).IsMatching(other.elemValue.(Tuple))
	}

	if e.elemType == NONE || other.elemType == NONE {
		return false
	}

	if e.elemType == ANY || other.elemType == ANY {
		return true
	}
	return false
}

// Comparator function, used for determining ordering of two elements.
// The order between elements of different type is arbitrary, but consistent.
// ANY < tuple < string < double < int < nil
// The order between elements of the same type is the builtin in golang
// Note to self: discussion about value receiver vs pointer receiver:
//
//	https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver-in-golang
//
// Returns 1 if this e > other, -1 if e < other, 0 if both are equal
func (e Elem) order(other Elem) int {
	switch e.elemType {

	case ANY:
		if other.elemType == ANY {
			return EQ
		}
		return LT

	case TUPLE:
		switch other.elemType {
		case ANY:
			return GT
		case TUPLE:
			return e.elemValue.(Tuple).order(other.elemValue.(Tuple))
		default:
			return LT
		}

	case STRING:
		switch other.elemType {
		case ANY:
		case TUPLE:
			return GT
		case STRING:
			if e.elemValue.(string) < other.elemValue.(string) {
				return LT
			}
			if e.elemValue.(string) == other.elemValue.(string) {
				return EQ
			}
			return GT
		case FLOAT:
		case INT:
		case NONE:
			return LT
		}

	case FLOAT:
		switch other.elemType {
		case ANY:
		case TUPLE:
		case STRING:
			return GT
		case FLOAT:
			if e.elemValue.(float64) < other.elemValue.(float64) {
				return LT
			}
			if e.elemValue.(float64) == other.elemValue.(float64) {
				return EQ
			}
			return GT
		case INT:
		case NONE:
			return LT
		}

	case INT:
		if other.elemType == INT {
			if e.elemValue.(int32) < other.elemValue.(int32) {
				return LT
			}
			if e.elemValue.(int32) == other.elemValue.(int32) {
				return EQ
			}
			return GT
		}
		if other.elemType == NONE {
			return GT
		}
		return LT
	default:
		return LT
	}
	return LT
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

func (t Tuple) String() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("(")

	size := len(t.elements)
	for i, e := range t.elements {
		strBuilder.WriteString(e.String())
		if i < size-1 {
			strBuilder.WriteString("|")
		}
	}
	strBuilder.WriteString(")")
	return strBuilder.String()
}

// MakeTuple creates a new Tuple instance from the given parameters.
func MakeTuple(element ...Elem) Tuple {
	var resultTuple Tuple
	// For debugging only.
	// for _, e := range element {
	// 	fmt.Printf("element %v of type %T (reflect %v)\n", e, e, reflect.TypeOf(e))
	// }
	resultTuple.elements = element
	// fmt.Printf("resulting tuple: %v", resultTuple)
	return resultTuple
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

// Comparator function, used for determining ordering of two tuples.
func (t Tuple) order(other Tuple) int {
	tSize := len(t.elements)
	otherSize := len(other.elements)

	shorterSize := min(tSize, otherSize)

	// Check each element for equality.
	for i := 0; i < shorterSize; i++ {
		if ord := t.elements[i].order(other.elements[i]); ord != EQ {
			return ord
		}
	}

	if tSize < otherSize {
		return LT
	}

	if tSize == otherSize {
		return EQ
	}

	return GT
}

// Comparator function for ordering tuples, used for tidwall/btree.
// Returns `true` if t1 is considered _less than_ t2, `false` otherwise.
func TupleOrder(t1, t2 Tuple) bool {
	return t1.order(t2) == LT
}
