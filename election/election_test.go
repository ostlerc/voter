package election

import "testing"

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
