package tuplespace

import (
	"fmt"
	"math"
	"strings"
)

const (
	// FloatPrecision sets the error for floating point comparison.
	FloatPrecision float64 = 0.0000001
)

type TupleElement uint

const (
	// ElemInt indicates 32bit-integers.
	ElemInt TupleElement = 1
	// ElemFloat indicates double precision (64bit) floating point numbers.
	ElemFloat = 2
	// ElemString indicates... well... strings.
	ElemString = 3
	// ElemTuple indicates a nested tuple.
	ElemTuple = 4
	// ElemAny indicates any possible type of the above, functioning as a wildcard.
	ElemAny = 5
	// ElemNone indicates an invalid type.
	ElemNone = 0
)

const (
	// LT is the `less than` return value for order comparisons.
	LT int = -1
	// EQ is the `equals` return value for order comparisons.
	EQ int = 0
	// GT is the `greater than` return value for order comparisons.
	GT int = 1
)

// Elem acts as an element container, holding a generic element and its type indication.
type Elem struct {
	elemType  TupleElement
	elemValue any
}

func (e Elem) GetType() TupleElement {
	return e.elemType
}

func (e Elem) GetValue() any {
	return e.elemValue
}

func (e Elem) String() string {
	switch e.elemType {
	case ElemInt:
		return fmt.Sprintf("%v", e.elemValue.(int))
	case ElemFloat:
		return fmt.Sprintf("%v", e.elemValue.(float64))
	case ElemString:
		return fmt.Sprintf("\"%v\"", e.elemValue.(string))
	case ElemTuple:
		return e.elemValue.(Tuple).String()
	case ElemAny:
		return "_"
	case ElemNone:
		return "nil"
	default:
		panic(fmt.Sprintf("Error: invalid elem type %T", e.elemValue))
	}
}

// Tuple element constructors /////////////////////////////////////////////////

// I instantiates an int-type tuple element.
func I(intVal int) Elem {
	return Elem{ElemInt, intVal}
}

// F instantiates a double precision (64bit) float64-type tuple element.
func F(floatVal float64) Elem {
	return Elem{ElemFloat, floatVal}
}

// S instantiates a string-type tuple element.
func S(stringVal string) Elem {
	return Elem{ElemString, stringVal}
}

// T instantiates a Tuple-type tuple element.
func T(tupleVal Tuple) Elem {
	return Elem{ElemTuple, tupleVal}
}

// Any instantiates a Wildcard tuple element.
func Any() Elem {
	return Elem{ElemAny, nil}
}

func None() Elem {
	return Elem{ElemNone, nil}
}

// IsDefined returns true if the element is defined, false if it is a wildcard or none.
func (e Elem) IsDefined() bool {
	switch e.elemType {
	case ElemInt:
		return true
	case ElemFloat:
		return true
	case ElemString:
		return true
	case ElemTuple:
		return e.elemValue.(Tuple).IsDefined()
	case ElemAny:
		return false
	case ElemNone:
		return false
	default:
		panic(fmt.Sprintf("Error: invalid elem type %T", e.elemValue))
	}
}

// Match two elements for equality, which is true either if they are of the same type and value
// or one or both are wildcards.
func (e Elem) isMatching(other Elem) bool {
	if e.elemType == ElemInt && other.elemType == ElemInt {
		return e.elemValue.(int) == other.elemValue.(int)
	}

	if e.elemType == ElemFloat && other.elemType == ElemFloat {
		return math.Abs(e.elemValue.(float64)-other.elemValue.(float64)) < FloatPrecision
	}

	if e.elemType == ElemString && other.elemType == ElemString {
		return e.elemValue.(string) == other.elemValue.(string)
	}
	if e.elemType == ElemTuple && other.elemType == ElemTuple {
		return e.elemValue.(Tuple).IsMatching(other.elemValue.(Tuple))
	}

	if e.elemType == ElemNone || other.elemType == ElemNone {
		return false
	}

	if e.elemType == ElemAny || other.elemType == ElemAny {
		return true
	}
	return false
}

// Comparator function, used for determining ordering of two elements.
// The order between elements of different type is arbitrary, but consistent.
// ElemAny < tuple < string < double < int < nil
// The order between elements of the same type is the builtin in golang
// Note to self: discussion about value receiver vs pointer receiver:
//
//	https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver-in-golang
//
// Returns 1 if this e > other, -1 if e < other, 0 if both are equal.
func (e Elem) order(other Elem) int {
	switch e.elemType {
	case ElemAny:
		return EQ

	case ElemTuple:
		return e.orderTuples(other)

	case ElemString:
		return e.orderStrings(other)

	case ElemFloat:
		return e.orderFloats(other)

	case ElemInt:
		return e.orderInts(other)

	default:
		return LT
	}
}

func (e Elem) orderTuples(other Elem) int {
	switch other.elemType {
	case ElemAny:
		return EQ
	case ElemTuple:
		return e.elemValue.(Tuple).order(other.elemValue.(Tuple))
	default:
		return LT
	}
}

func (e Elem) orderStrings(other Elem) int {
	switch other.elemType {
	case ElemAny:
		return EQ
	case ElemTuple:
		return GT
	case ElemString:
		if e.elemValue.(string) < other.elemValue.(string) {
			return LT
		}
		if e.elemValue.(string) == other.elemValue.(string) {
			return EQ
		}
		return GT
	case ElemFloat:
	case ElemInt:
	case ElemNone:
		return LT
	}
	return LT
}

func (e Elem) orderFloats(other Elem) int {
	switch other.elemType {
	case ElemAny:
		return EQ
	case ElemTuple:
	case ElemString:
		return GT
	case ElemFloat:
		if e.elemValue.(float64) < other.elemValue.(float64) {
			return LT
		}
		if e.elemValue.(float64) == other.elemValue.(float64) {
			return EQ
		}
		return GT
	case ElemInt:
	case ElemNone:
		return LT
	}
	return LT
}

func (e Elem) orderInts(other Elem) int {
	if other.elemType == ElemAny {
		return EQ
	}
	if other.elemType == ElemInt {
		if e.elemValue.(int) < other.elemValue.(int) {
			return LT
		}
		if e.elemValue.(int) == other.elemValue.(int) {
			return EQ
		}
		return GT
	}
	if other.elemType == ElemNone {
		return GT
	}
	return LT
}

// Tuple can contain elements of five different data types:
// - integers
// - floating point numbers
// - strings
// - tuples themselves
// - wildcards
// .
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

// IsDefined returns true if the tuple does not contain any wildcards or none fields.
func (t Tuple) IsDefined() bool {
	for _, v := range t.elements {
		if !v.IsDefined() {
			// break out of the loop if any element is undefined
			return false
		}
	}
	return true
}

// IsMatching checks two tuples for equality, which is true if
// - they are of the same length AND
// - each element of one matches the others.
func (t Tuple) IsMatching(other Tuple) bool {
	tSize := len(t.elements)
	otherSize := len(other.elements)

	// Check length of tuples first
	if tSize != otherSize {
		return false
	}

	// Check each element for equality.
	for i := range tSize {
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
	for i := range shorterSize {
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

// TupleOrder is a comparator function for ordering tuples, used for tidwall/btree.
// Returns `true` if t1 is considered _less than_ t2, `false` otherwise.
func TupleOrder(t1, t2 Tuple) bool {
	return t1.order(t2) == LT
}
