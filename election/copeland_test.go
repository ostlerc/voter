package election

import "testing"

func TestCopeland(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}

	v := e.Copeland()[0]
	if v != -1 {
		t.Fatal("Invalid copeland result", v)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 1}},
		N: 3,
	}

	v = e.Copeland()[0]
	if v != 1 {
		t.Fatal("Invalid copeland result", v)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 3}},
		N: 3,
	}

	v = e.Copeland()[0]
	if v != 2 {
		t.Fatal("Invalid copeland result", v)
	}
}
