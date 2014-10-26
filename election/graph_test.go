package election

import (
	"encoding/json"
	"testing"
)

const testGraph = `{
	"nodes": {
		"2": {
			"edges": {
				"1": 10
			}
		},
		"1": {},
		"0": {
			"edges": {
				"1": 2
			}
		}
	}
}`

func TestSlater(t *testing.T) {
	j := &jsonegraph{}
	err := json.Unmarshal([]byte(testGraph), j)
	if err != nil {
		t.Fatal(err)
	}
	g := j.egraph()
	if d, w := g.slater([]int{0, 1, 2}); d != 1 && w != 10 {
		t.Fatal("Incorrect slater results", d, w)
	}
	if d, w := g.slater([]int{0, 2, 1}); d != 0 && w != 0 {
		t.Fatal("Incorrect slater results", d, w)
	}
	if d, w := g.slater([]int{1, 0, 2}); d != 2 && w != 12 {
		t.Fatal("Incorrect slater results", d, w)
	}
	if d, w := g.slater([]int{1, 2, 0}); d != 2 && w != 12 {
		t.Fatal("Incorrect slater results", d, w)
	}
	if d, w := g.slater([]int{2, 0, 1}); d != 0 && w != 0 {
		t.Fatal("Incorrect slater results", d, w)
	}
	if d, w := g.slater([]int{2, 1, 0}); d != 1 && w != 2 {
		t.Fatal("Incorrect slater results", d, w)
	}

	e, w := g.Slater()
	if e[0] != 0 || e[1] != 2 || e[2] != 1 || w[0] != 0 || w[1] != 2 || w[2] != 1 {
		t.Fatal("Incorrect slater result", e, w)
	}
}
