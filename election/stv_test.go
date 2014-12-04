package election

import (
	"fmt"
	"testing"
)

func TestSTV(t *testing.T) {

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

	v := e.CandVotes()
	if !ArEq(v[0], []int{6, 12, 2}) {
		fmt.Println("Incorrect", v[0])
	}
	if !ArEq(v[1], []int{9, 5, 6}) {
		fmt.Println("Incorrect", v[1])
	}
	if !ArEq(v[2], []int{5, 3, 12}) {
		fmt.Println("Incorrect", v[2])
	}
}

func TestDistribute(t *testing.T) {
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
	_t.Tally(e)
}
