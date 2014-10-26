package election

import (
	"encoding/json"
	"strings"
	"testing"
)

const (
	testGraph = `{
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
)

var (
	d1 = []string{`digraph G {`, `1 -> 2 [label="3"];`, `1 -> 0 [label="3"];`, `2 -> 0 [label="3"];`}
	j1 = []string{`{"nodes":{`, `"0":{}`, `"1":{"edges":{`, `"0":3`, `"2":3`, `"2":{"edges":{"0":3}`}
)

func has(str string, l []string) bool {
	for _, s := range l {
		if strings.Index(str, s) == -1 {
			return false
		}
	}

	return true
}

func TestGraphEdgesDot(t *testing.T) {
	e := &Election{V: []*Vote{
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1},
		&Vote{C: map[string]int{"0": 1, "1": 2, "2": 0}, W: 1}},
		N: 3,
	}
	d := e.Graph().Edges().Dot()
	if !has(d, d1) {
		t.Fatal("Invalid dot graph generated", d)
	}
	j := e.Graph().Edges().JSON()
	if !has(j, j1) {
		t.Fatal("Invalid json", j)
	}
}

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
