package election

import (
	"strings"
	"testing"
)

func TestCmp(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
	}

	v := e.cmp(1, 2)
	if v != -3 {
		t.Fatal("Incorrect cmp", v)
	}
	v = e.cmp(2, 1)
	if v != 3 {
		t.Fatal("Incorrect cmp", v)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 1}},
	}
	v = e.cmp(1, 0)
	if v != 0 {
		t.Fatal("Incorrect cmp", v)
	}
	v = e.cmp(0, 1)
	if v != 0 {
		t.Fatal("Incorrect cmp", v)
	}
}

func TestCondorcet(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3}

	if e.Condorcet() != 1 {
		t.Fatal("Incorrect condorcet winner", e.Condorcet())
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 0, "2": 2}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 1}},
		N: 3}

	if e.Condorcet() != -1 {
		t.Fatal("Incorrect condorcet winner", e.Condorcet())
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 0}, W: 1},
		&Vote{C: map[string]int{"0": 0, "1": 1}, W: 1}},
		N: 2}

	if e.Condorcet() != -1 {
		t.Fatal("Incorrect condorcet winner", e.Condorcet())
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 0}, W: 1}},
		N: 2}

	if e.Condorcet() != 1 {
		t.Fatal("Incorrect condorcet winner", e.Condorcet())
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 1}, W: 1},
		&Vote{C: map[string]int{"0": 0, "1": 1}, W: 1}},
		N: 2}

	if e.Condorcet() != 0 {
		t.Fatal("Incorrect condorcet winner", e.Condorcet())
	}
}

func TestScore(t *testing.T) {
	v := &Vote{C: map[string]int{"0": 0, "1": 1, "2": 2}}
	s := v.Score([]int{0, 1, 2})
	if s != 0 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{0, 2, 1})
	if s != 3 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{1, 0, 2})
	if s != 5 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{1, 2, 0})
	if s != 11 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{2, 1, 0})
	if s != 12 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{2, 0, 1})
	if s != 9 {
		t.Fatal("Incorrect Score", s)
	}

	v = &Vote{C: map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4}}
	s = v.Score([]int{0, 1, 2, 3, 4})
	if s != 0 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{4, 3, 2, 1, 0})
	if s != 54 {
		t.Fatal("Incorrect Score", s)
	}

	v = &Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}}
	s = v.Score([]int{2, 0, 1})
	if s != 6 {
		t.Fatal("Incorrect Score", s)
	}

	s = v.Score([]int{2, 1, 0})
	if s != 11 {
		t.Fatal("Incorrect Score", s)
	}
}

var csvelection = `weight,5,4,3,6
Alex,5,1,6,4
Bart,1,6,5,5
Cindy,2,3,7,3
David,4,4,1,2
Erik,6,5,3,1
Frank,3,2,2,6
Greg,7,7,4,7`

func TestCSV(t *testing.T) {
	r := strings.NewReader(csvelection)
	e := ParseFrom("csv", r)
	if e.CSV() != csvelection {
		t.Fatal("Invalid csv conversion\n'", e.CSV(), "'\n'", csvelection, "'")
	}
}

func TestPref(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	p := e.Pref()
	if p.First != 1 || p.Second != 0 {
		t.Fatal("Invalid Preference ", p)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 2, "1": 0, "2": 1}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	p = e.Pref()
	if p.First != 2 || p.Second != 0 {
		t.Fatal("Invalid Preference ", p)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 0, "1": 2, "2": 1}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	p = e.Pref()
	if p != nil {
		t.Fatal("Invalid Preference ", p)
	}
}

func TestPeak(t *testing.T) {
	v := &Vote{C: map[string]int{"0": 0, "1": 1, "2": 2, "3": 3}}
	if p := v.PeakValue(); p != 0 {
		t.Fatal("Incorrect peak value", p)
	}

	v = &Vote{C: map[string]int{"0": 1, "1": 0, "2": 2, "3": 3}}
	if p := v.PeakValue(); p != 1 {
		t.Fatal("Incorrect peak value", p)
	}

	v = &Vote{C: map[string]int{"0": 2, "1": 0, "2": 1, "3": 3}}
	if p := v.PeakValue(); p != -1 {
		t.Fatal("Incorrect peak value", p)
	}

	v = &Vote{C: map[string]int{"0": 2, "1": 1, "2": 0, "3": 3}}
	if p := v.PeakValue(); p != 2 {
		t.Fatal("Incorrect peak value", p)
	}

	v = &Vote{C: map[string]int{"0": 2, "1": 1, "2": 3, "3": 0}}
	if p := v.PeakValue(); p != 2 {
		t.Fatal("Incorrect peak value", p)
	}

	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	if p := e.Peak(); p != 1 {
		t.Fatal("Invalid election peak", p)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	if p := e.Peak(); p != -1 {
		t.Fatal("Invalid election peak", p)
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 2, "1": 1, "2": 0}, W: 4},
		&Vote{C: map[string]int{"0": 2, "1": 0, "2": 1}, W: 1}},
		N: 3,
	}
	if p := e.Peak(); p != -1 {
		t.Fatal("Invalid election peak", p)
	}
}

func TestManipulation(t *testing.T) {
	/*e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	tally := GetTally("stv")
	m := e.FindManipulation(tally)
	if m == nil {
		t.Fatal("Invalid, should have found a manipulation")
	}*/
}

func TestVotes(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}

	if e.Votes() != 3 {
		t.Fatal("Incorrect # of votes")
	}

	e = &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 14},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 25}},
		N: 3,
	}

	if e.Votes() != 40 {
		t.Fatal("Incorrect # of votes", e.Votes())
	}
}

func TestAfter(t *testing.T) {
	v := &Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 25}
	if v.after(1) != 2 {
		t.Fatal("Incorrect after")
	}
	if v.after(2) != 0 {
		t.Fatal("Incorrect after")
	}
	if v.after(0) != -1 {
		t.Fatal("Incorrect after")
	}
}

func TestRemoveCandidate(t *testing.T) {
	v := &Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 4}
	v2 := v.RemoveCandidate(1)
	if len(v2.C) != 2 || len(v.C) != 3 {
		t.Fatal("Incorrect candidate removal")
	}
	if v2.C["0"] != 1 || v2.C["1"] != 0 || v.C["0"] != 1 || v.C["1"] != 2 || v.C["2"] != 0 {
		t.Fatal("Incorrect candidate removal", v2.C)
	}
}

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
