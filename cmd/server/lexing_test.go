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
	}
	var expected = [...][]ts.Tuple{
		{},
		{ts.MakeTuple(ts.I(1))},
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
				isDifferent = true
				break
			}
		}

		if isDifferent {
			t.Errorf("output and expected are not equal!\noutput: %v\nexpected: %v", output, expected)
		}
	}
}
