package election

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	fix      = flag.Bool("fix", false, "Use fixed random seed")
	R        = rand.New(rand.NewSource(time.Now().Unix()))
	talliers = make(map[string]Tallier)
)

type Election struct {
	N int               `json:"candidates"`
	F *IntPair          `json:"pref,omitempty"`
	P *bool             `json:"peak,omitempty"`
	C *int              `json:"condorcet,omitempty"`
	R []int             `json:"ranks"`
	M map[string]string `json:"names,omitempty"`
	V []*Vote           `json:"votes"`
}

type Vote struct {
	C    map[string]int `json:"vote"`
	W    int            `json:"weight"`
	Peak *int           `json:"peak,omitempty"`
}

type Tallier interface {
	Tally(*Election) []int
	Key() string
}

func RegisterTally(t Tallier) {
	talliers[t.Key()] = t
}

func GetTally(s string) Tallier {
	return talliers[s]
}

func TallyKeys() []string {
	res := make([]string, 0)
	for k, _ := range talliers {
		res = append(res, k)
	}
	return res
}

func Setup() { //Requires flag.parse to have been called
	if *fix {
		R = rand.New(rand.NewSource(99))
	}
}

// Returns a score of how similar the votes are. 0 is exact match
func (v *Vote) Score(r []int) int {
	h := make(map[int]int) //map candidates to index
	for k, v := range v.C {
		ki, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		h[v] = ki
	}

	res := 0
	for i := 0; i < len(r); i++ {
		w := (r[i] - len(r)/2)
		if w < 0 {
			w = -w
			w += len(r) - r[i]
		} else {
			w++
		}
		rx := (r[i] - h[i]) * w
		if rx < 0 {
			res += -rx
		} else {
			res += rx
		}
	}
	return res
}

func (v *Vote) PeakValue() int {
	a := v.C["0"] - 1
	b := v.C["0"] + 1
	for i := 1; i < len(v.C); i++ {
		is := strconv.Itoa(i)
		val := v.C[is]
		if (a != -1 && a != val) && (b != len(v.C) && b != val) {
			return -1
		} else if a == val {
			a--
		} else {
			b++
		}
	}
	return v.C["0"]
}

//cmp returns <0 if a beats b, >0 if b beats a and 0 if a tie
func (e *Election) cmp(a, b int) int {
	cnt := 0

	for _, v := range e.V {
		for i := 0; i < len(v.C); i++ {
			x := v.C[strconv.Itoa(i)]
			if x == a {
				cnt -= v.W
				break
			}
			if x == b {
				cnt += v.W
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

func (e *Election) Pref() *IntPair {
	c := 0
	for _, v := range e.V {
		c += v.W
	}
	for i := 0; i < e.N; i++ {
		for j := i + 1; j < e.N; j++ {
			r := e.cmp(i, j)
			if r == c {
				return &IntPair{First: j, Second: i}
			} else if r == -c {
				return &IntPair{First: i, Second: j}
			}
		}
	}

	return nil
}

func (e *Election) CSV() string {
	lines := make(map[int][]int)
	res := ""
	res += fmt.Sprint(",")
	for i, v := range e.V {
		res += fmt.Sprint(v.W)
		if i+1 != len(e.V) {
			res += fmt.Sprint(",")
		}
	}

	for i := 0; i < e.N; i++ {
		lines[i] = make([]int, len(e.V))
	}

	for i, v := range e.V {
		for j := 0; j < e.N; j++ {
			lines[j][i] = v.C[strconv.Itoa(j)] + 1
		}
	}

	res += fmt.Sprint("\n")
	for i := 0; i < e.N; i++ {
		if len(e.M) > 0 {
			res += fmt.Sprint(e.M[strconv.Itoa(i)], ",")
		} else {
			res += fmt.Sprint(i, ",")
		}
		for j := 0; j < len(e.V); j++ {
			res += fmt.Sprint(lines[i][j])
			if j+1 != len(e.V) {
				res += fmt.Sprint(",")
			}
		}
		if i+1 != e.N {
			res += fmt.Sprint("\n")
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

func NewVote(w int) *Vote {
	return &Vote{
		C: make(map[string]int),
		W: w,
	}
}

func GetName(M map[string]string, v string) string {
	if s, ok := M[v]; ok {
		return s
	}
	return v
}
