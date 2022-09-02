package gotupolis

import (
	"testing"

	opt "github.com/micutio/goptional"
	ts "github.com/micutio/gotupolis/pkg/tuplespace"
)

func TestLexer(t *testing.T) {
	var inputs = [...]string{
		"",
		"(1)",
	}
	var expected = [...][]opt.Maybe[ts.Tuple]{
		{opt.NewNothing[ts.Tuple]()},
		{opt.NewJust[ts.Tuple](ts.MakeTuple(ts.I(1)))},
	}

	for i := range inputs {
		checkOutput(t, inputs[i], expected[i])
	}
}

func checkOutput(t *testing.T, input string, expected []opt.Maybe[ts.Tuple]) {
	var output []opt.Maybe[ts.Tuple] = NewLexer(input).IntoTuples()

	if len(output) != len(expected) {
		t.Errorf("output and expected are of diffent length!\noutput: %v\nexpected: %v", output, expected)
	} else {
		isDifferent := false
		for i := range output {
			if output[i] != expected[i] {
				isDifferent = true
				break
			}
		}

		if isDifferent {
			t.Errorf("output and expected are not equal!\noutput: %v\nexpected: %v", output, expected)
		}
	}
}
