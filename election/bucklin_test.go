package election

import "testing"

func TestBucklin(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 5, "1": 2, "2": 4, "3": 3, "4": 1}, W: 1},
		&Vote{C: map[string]int{"0": 3, "1": 4, "2": 2, "3": 5, "4": 1}, W: 3},
		&Vote{C: map[string]int{"0": 4, "1": 2, "2": 3, "3": 5, "4": 1}, W: 5},
		&Vote{C: map[string]int{"0": 4, "1": 1, "2": 2, "3": 5, "4": 3}, W: 4},
		&Vote{C: map[string]int{"0": 4, "1": 5, "2": 3, "3": 1, "4": 2}, W: 5}},
		N: 5,
	}

	b := e.Bucklin(1)
	if b != 4 {
		t.Fatal("Incorrect bucklin value", b)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}

	b = e.Bucklin(1)
	if b != -1 {
		t.Fatal("Incorrect bucklin value", b)
	}
	b = e.Bucklin(2)
	if b != 1 {
		t.Fatal("Incorrect bucklin value", b)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 2},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 5}},
		N: 3,
	}

	b = e.Bucklin(1)
	if b != -1 {
		t.Fatal("Incorrect bucklin value", b)
	}
	b = e.Bucklin(2)
	if b != 1 {
		t.Fatal("Incorrect bucklin value", b)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 3, "1": 5, "2": 2, "3": 1, "4": 4}, W: 2},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 3, "3": 5, "4": 4}, W: 4},
		&Vote{C: map[string]int{"0": 3, "1": 2, "2": 4, "3": 1, "4": 5}, W: 3},
		&Vote{C: map[string]int{"0": 2, "1": 4, "2": 1, "3": 5, "4": 3}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 5, "3": 4, "4": 3}, W: 2}},
		N: 5,
	}

	tally := GetTally("bucklin")
	res := tally.Tally(e)
	if res[0] != 2 || res[1] != 0 {
		t.Fatal(res)
	}
}
