package election

import (
	"flag"
	"math/rand"
	"time"
)

var (
	random     = flag.Bool("rand", false, "Activate random generation")
	fix        = flag.Bool("fix", false, "Use fixed random seed")
	Votes      = flag.Int("Votes", 6, "Number of voters in election")
	Candidates = flag.Int("cand", 3, "Number of candidates in election")

	R = rand.New(rand.NewSource(time.Now().Unix()))
)

type Vote struct {
	C map[string]int `json:"vote"`
}

type Election struct {
	N int    `json:"candidates"`
	F int    `json:"favorite,omitempty"`
	P int    `json:"peak,omitempty"`
	C bool   `json:"condorcet"`
	V []Vote `json:"votes"`
}

func init() {
	flag.Parse()
	if *random {
		*Votes = max(R.Intn(*Votes), 3)           //minimum of 3 votes
		*Candidates = max(R.Intn(*Candidates), 3) //minimum of 3 votes
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RandCand(l []int) int {
	return R.Intn(len(l))
}

func RemoveAt(i int, l []int) []int {
	return append(l[:i], l[i+1:]...)
}

func RemoveValue(v int, l []int) []int {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return RemoveAt(i, l)
		}
	}
	panic("No value")
}

func Contains(v int, l []int) bool {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return true
		}
	}
	return false
}

func New() *Election {
	return &Election{
		V: make([]Vote, *Votes),
		N: *Candidates,
	}
}

func NewVote() Vote {
	return Vote{
		C: make(map[string]int),
	}
}
