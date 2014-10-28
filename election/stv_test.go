package election

import "testing"

func TestSTV(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}

	v := e.STV()
	if !ArEq(v, []int{1, 2, 0}) {
		t.Fatal("Invalid stv result", v)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}

	v = e.STV()
	if !ArEq(v, []int{2, 1, 0}) {
		t.Fatal("Invalid stv result", v)
	}
}
