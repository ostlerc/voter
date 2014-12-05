package election

import "testing"

func TestCandVotes(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 1, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 2},
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 8},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 0, "2": 1}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 1}},
		N: 3,
	}

	v1 := &Vote{C: map[string]int{"0": 1, "1": 0}, W: 1}
	v2 := &Vote{C: map[string]int{"0": 0, "1": 1}, W: 2}
	e = &Election{V: []*Vote{v1, v2}, N: 2}

	v := e.CandVotes()
	if len(v[0]) != 2 {
		t.Fatal("Incorrect Size", len(v[0]), tojson(v[0]))
	}
	if tojson(v[0][0]) != tojson(v2) {
		t.Fatal("Incorrect element", tojson(v[0][0]))
	}
}

func TestSTV(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 1, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 2},
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 8},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 0, "2": 1}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 1}},
		N: 3,
	}

	_t := TallySTV{}
	res := _t.Tally(e)
	if !ArEq(res, []int{1, 0, 2}) {
		t.Fatal("Incorrect stv result", res)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 1, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 0, "2": 1}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}
	c := e.CandVotes()
	if !ArEq(c.Lens(), []int{8, 8, 8}) {
		t.Fatal("Incorrect CandVotes")
	}
	_t = TallySTV{}
	res = _t.Tally(e)
	if !ArEq(res, []int{0, 1, 2}) {
		t.Fatal("Incorrect stv result", res)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 1, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4}},
		N: 3,
	}
	c = e.CandVotes()
	if !ArEq(c.Lens(), []int{8, 8, 8}) {
		t.Fatal("Incorrect CandVotes")
	}
	_t = TallySTV{}
	res = _t.Tally(e)
	if !ArEq(res, []int{1, 0, 2}) {
		t.Fatal("Incorrect stv result", res)
	}
}

func TestFirst(t *testing.T) {
	tally := &TallySTV{
		Ignore: make(Ints, 0),
	}

	v := &Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4}

	f := tally.first(v)
	if f != 1 {
		t.Fatal("Incorrect first", f)
	}

	tally.Ignore = append(tally.Ignore, 1)
	f = tally.first(&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4})
	if f != 2 {
		t.Fatal("Incorrect first", f)
	}

	tally.Ignore = append(tally.Ignore, 2)
	f = tally.first(&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4})
	if f != 0 {
		t.Fatal("Incorrect first", f)
	}
}
