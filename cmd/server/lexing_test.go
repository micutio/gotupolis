package gotupolis

import (
	"fmt"
	"testing"

	ts "github.com/micutio/gotupolis/pkg/tuplespace"
)

func TestLexer(t *testing.T) {
	var inputs = [...]string{
		"",
		"(1)",
		"(2.0,3.14,\"steeze\")",
		"(1,2.0,\"3.14\",_,666),(\"foo\",\"bar\")",
		"(1,2.0,\"3.14\",_,666), (\"foo\",\"bar\"), ((1,-2,-3),(-4.1,5.1,\"6.1\"))",
	}
	var expected = [...][]ts.Tuple{
		{},
		{ts.MakeTuple(ts.I(1))},
		{ts.MakeTuple(ts.F(2.0), ts.F(3.14), ts.S("steeze"))},
		{ts.MakeTuple(ts.I(1), ts.F(2.0), ts.S("3.14"), ts.Any(), ts.I(666)),
			ts.MakeTuple(ts.S("foo"), ts.S("bar"))},
		{ts.MakeTuple(ts.I(1), ts.F(2.0), ts.S("3.14"), ts.Any(), ts.I(666)),
			ts.MakeTuple(ts.S("foo"), ts.S("bar")),
			ts.MakeTuple(
				ts.T(ts.MakeTuple(ts.I(1), ts.I(-2), ts.I(-3))),
				ts.T(ts.MakeTuple(ts.F(-4.1), ts.F(5.1), ts.S("6.1"))),
			)},
	}

	for i := range inputs {
		checkOutput(t, inputs[i], expected[i])
	}
}

func checkOutput(t *testing.T, input string, expected []ts.Tuple) {
	lexer := NewLexer(input)
	output, err := lexer.IntoTuples()
	if err != nil {
		fmt.Println(err)
	}

	if len(output) != len(expected) {
		t.Errorf("output and expected are of diffent length!\noutput: %v\nexpected: %v", output, expected)
	} else {
		isDifferent := false
		for i := range output {
			if !output[i].IsMatching(expected[i]) {
				fmt.Printf("element %v not matching %v", output[i], expected[i])
				isDifferent = true
				break
			}
		}

		if isDifferent {
			t.Errorf("output and expected are not equal!\noutput: %v\nexpected: %v", output, expected)
		}
	}
}
