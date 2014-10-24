package election

import (
	"flag"
	"math/rand"
	"strconv"
	"time"
)

var (
	random     = flag.Bool("rand", false, "Use a random vote/cand count")
	fix        = flag.Bool("fix", false, "Use fixed random seed")
	Votes      = flag.Int("vote", 6, "Number of voters in election")
	Candidates = flag.Int("cand", 3, "Number of candidates in election")

	R = rand.New(rand.NewSource(time.Now().Unix()))
)

type Vote struct {
	C    map[string]int `json:"vote"`
	Peak *int           `json:"peak,omitempty"`
}

type IntPair struct {
	First  int `json:"a"`
	Second int `json:"b"`
}

type Election struct {
	N int      `json:"candidates"`
	F *IntPair `json:"pref,omitempty"`
	P *bool    `json:"peak,omitempty"`
	C *int     `json:"condorcet,omitempty"`
	R []int    `json:"rank"`
	V []*Vote  `json:"votes"`
}

func Setup() { //Requires flag.parse to have been called
	if *fix {
		R = rand.New(rand.NewSource(99))
	}

	if *random {
		*Votes = max(R.Intn(*Votes), 3)           //minimum of 3 votes
		*Candidates = max(R.Intn(*Candidates), 3) //minimum of 3 candidates
	}
}

//cmp returns -1 if a beats b, 1 if b beats a and 0 if a tie
func (e *Election) cmp(a, b int) int {
	acnt := 0
	bcnt := 0

	for _, v := range e.V {
		for i := 0; i < len(v.C); i++ {
			x := v.C[strconv.Itoa(i)]
			if x == a {
				acnt++
				break
			}
			if x == b {
				bcnt++
				break
			}
		}
	}
	if acnt == bcnt {
		return 0
	}
	if acnt > bcnt {
		return -1
	}
	return 1
}

//Rank returns int array of how many individual wins each candidate has
func (e *Election) Rank() []int {
	res := make([]int, len(e.V))
	for i := 0; i < len(e.V); i++ {
		for j := i + 1; j < len(e.V); j++ {
			r := e.cmp(i, j)
			if r == -1 {
				res[i]++
			} else if r == 1 {
				res[j]++
			}
		}
	}
	return res
}

//Condorcet returns the condorcet winner as an int. -1 is returned if no winner is found
func (e *Election) Condorcet() int {
	max := 0
	imax := 0
	res := e.Rank()
	for i, v := range res {
		if v > max {
			imax = i
			max = v
		}
	}

	f := false
	for _, v := range res {
		if v == max {
			if f {
				return -1
			}
			f = true
		}
	}

	return imax
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

func Index(v int, l []int) int {
	for i := 0; i < len(l); i++ {
		if v == l[i] {
			return i
		}
	}
	return -1
}

func (v *Vote) Contains(k string) bool {
	_, ok := v.C[k]
	return ok
}

func (v *Vote) Rank(k int) string {
	for key, v := range v.C {
		if v == k {
			return key
		}
	}
	panic("Not found")
}

func Strn(n int) string {
	return strconv.Itoa(R.Intn(n))
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
		V: make([]*Vote, *Votes),
		N: *Candidates,
	}
}

func NewVote() *Vote {
	return &Vote{
		C: make(map[string]int),
	}
}
