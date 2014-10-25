package election

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	fix = flag.Bool("fix", false, "Use fixed random seed")

	R = rand.New(rand.NewSource(time.Now().Unix()))
)

type Election struct {
	N int               `json:"candidates"`
	F *IntPair          `json:"pref,omitempty"`
	P *bool             `json:"peak,omitempty"`
	C *int              `json:"condorcet,omitempty"`
	R []int             `json:"ranks"`
	M map[string]string `json:"names"`
	V []*Vote           `json:"votes"`
}

type Vote struct {
	C    map[string]int `json:"vote"`
	Peak *int           `json:"peak,omitempty"`
}

func Setup() { //Requires flag.parse to have been called
	if *fix {
		R = rand.New(rand.NewSource(99))
	}
}

//cmp returns <0 if a beats b, >0 if b beats a and 0 if a tie
func (e *Election) cmp(a, b int) int {
	cnt := 0

	for _, v := range e.V {
		for i := 0; i < len(v.C); i++ {
			x := v.C[strconv.Itoa(i)]
			if x == a {
				cnt--
				break
			}
			if x == b {
				cnt++
				break
			}
		}
	}
	return cnt
}

//Rank returns int array of how many individual wins each candidate has
func (e *Election) Rank() []int {
	res := make([]int, e.N)
	for i := 0; i < e.N; i++ {
		for j := i + 1; j < e.N; j++ {
			r := e.cmp(i, j)
			if r < 0 {
				res[i]++
			} else if r > 0 {
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
	e.R = e.Rank()
	for i, v := range e.R {
		if v > max {
			imax = i
			max = v
		}
	}

	f := false
	for _, v := range e.R {
		if v == max {
			if f {
				return -1
			}
			f = true
		}
	}

	return imax
}

func (v *Vote) Contains(k string) bool {
	_, ok := v.C[k]
	return ok
}

func (v *Vote) Prefer(f *IntPair) {
	aidx := v.Rank(f.First)
	bidx := v.Rank(f.Second)
	if aidx > bidx {
		v.C[aidx] = f.Second
		v.C[bidx] = f.First
	}
}

func (v *Vote) Rank(k int) string {
	for key, v := range v.C {
		if v == k {
			return key
		}
	}
	log.Println(k, v.C)
	panic("Not Found")
}

func NewPref(c int) *IntPair {
	f := &IntPair{ //prefer a over b (always)
		First:  R.Intn(c),
		Second: R.Intn(c),
	}

	for f.First == f.Second { //verify unique random values
		f.Second = R.Intn(c)
	}
	return f
}

func New(v, c int) *Election {
	return &Election{
		V: make([]*Vote, v),
		N: c,
		M: make(map[string]string),
	}
}

func NewVote() *Vote {
	return &Vote{
		C: make(map[string]int),
	}
}
