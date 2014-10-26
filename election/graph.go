package election

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var (
	talliers = make(map[string]Tallier)
)

type Egraph struct {
	// Majority graph: [parent node][child node]edge weight
	Nodes map[int]map[int]int
	names map[string]string
}

type Tallier interface {
	Tally(*Election) []int
}

func RegisterTally(s string, t Tallier) {
	talliers[s] = t
}

func GetTally(s string) Tallier {
	return talliers[s]
}

func (e *Election) Graph() *Egraph {
	g := &Egraph{
		Nodes: make(map[int]map[int]int),
		names: e.M,
	}
	for i := 0; i < e.N; i++ {
		g.Nodes[i] = make(map[int]int)
	}
	for i := 0; i < e.N; i++ {
		for j := i + 1; j < e.N; j++ {
			r := e.cmp(i, j)
			if r < 0 { //i wins
				g.Nodes[i][j] += -r
			} else if r > 0 { //j wins
				g.Nodes[j][i] += r
			}
		}
	}
	return g
}

func (e *Egraph) Edges() *Egraph {
	g := &Egraph{
		Nodes: make(map[int]map[int]int),
		names: e.names,
	}
	for i := 0; i < len(e.Nodes); i++ {
		g.Nodes[i] = make(map[int]int)
	}

	for a, m := range e.Nodes {
		for b, w := range m {
			if w != e.Nodes[b][a] {
				g.Nodes[a][b] += (w - e.Nodes[b][a])
			}
		}
	}

	return g
}

// Dot returns a dot file
func (e *Egraph) Dot() string {
	res := "digraph G {\n"
	for a, m := range e.Nodes {
		for b, w := range m {
			if w != 0 {
				res += fmt.Sprintf("\t"+`%s -> %s [label="%d"];`+"\n", GetName(e.names, strconv.Itoa(a)), GetName(e.names, strconv.Itoa(b)), w)
			}
		}
	}
	return res + "}"
}

//the soul purpose of this struct is to allow json marshalling
type jsonegraph struct {
	Nodes map[string]node `json:"nodes"`
}

//the soul purpose of this struct is to allow json marshalling
type node struct {
	Edges map[string]int `json:"edges,omitempty"`
}

func (j *jsonegraph) egraph() *Egraph {
	res := &Egraph{Nodes: make(map[int]map[int]int)}
	for a, m := range j.Nodes {
		ai, _ := strconv.Atoi(a)
		res.Nodes[ai] = make(map[int]int)
		for b, w := range m.Edges {
			if w != 0 {
				bi, _ := strconv.Atoi(b)
				res.Nodes[ai][bi] = w
			}
		}
	}
	return res
}

func (e *Egraph) JSON() string {
	res := &jsonegraph{Nodes: make(map[string]node)}
	for a, m := range e.Nodes {
		astr := strconv.Itoa(a)
		res.Nodes[astr] = node{Edges: make(map[string]int)}
		for b, w := range m {
			if w != 0 {
				bstr := strconv.Itoa(b)
				res.Nodes[astr].Edges[bstr] = w
			}
		}
	}

	v, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	return string(v)
}
