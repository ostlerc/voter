package election

import (
	"fmt"
	"testing"
)

func TestAllPermutations(t *testing.T) {
	p := Perms(3)
	if len(p) != 6 {
		t.Fatal("incorrect length")
	}

	//easier to test eq with strings. yeah I could write it but easier this way
	s := fmt.Sprintf("%v", p)
	if s != "[[0 1 2] [0 2 1] [1 0 2] [1 2 0] [2 0 1] [2 1 0]]" {
		t.Fatal("Incorrect perms", s)
	}
}
