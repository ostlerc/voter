package election

import "testing"

func TestBucklin(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
	}

	b := e.Bucklin(1)
	if b != -1 {
		t.Fatal("Incorrect bucklin value", b)
	}
	b = e.Bucklin(2)
	if b != 1 {
		t.Fatal("Incorrect bucklin value", b)
	}
}
