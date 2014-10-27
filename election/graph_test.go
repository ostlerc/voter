package election

import (
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
