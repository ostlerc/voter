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
	R        = rand.New(rand.NewSource(time.Now().UnixNano()))
	talliers = make(map[string]Tallier)
)

type Election struct {
	N int               `json:"candidates"`
	F *IntPair          `json:"pref,omitempty"`
	P *int              `json:"peak,omitempty"`
	C *int              `json:"condorcet,omitempty"`
	R []int             `json:"ranks"`
	M map[string]string `json:"names,omitempty"`
	V []*Vote           `json:"votes"`
}

type Vote struct {
	C    map[string]int `json:"vote"`
	W    int            `json:"weight"`
	D    bool           `json:"diminish"`
	Peak *int           `json:"peak,omitempty"`
}

type IrrelevantCand struct {
	Candidate     string `json:"candidate"`
	ChangesWinner bool   `json:"causes_change"`
}

type TallyResult struct {
	Results    map[string][]int           `json:"results"`
	Names      map[string]string          `json:"names,omitempty"`
	Condorcet  *int                       `json:"condorcet,omitempty"`
	Election   *Election                  `json:"election,omitempty"`
	M          map[string]*Manipulation   `json:"manipulations,omitempty"`
	PrefIntact map[string]bool            `json:"pref_intact,omitempty"`
	Irrelevant map[string]*IrrelevantCand `json:"irr_cand,omitempty"`
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

var DumbScore = false

// Returns a score of how similar the votes are. 0 is exact match
func (v *Vote) Score(r []int) int {
	if DumbScore {
		return v.DumbScore(r)
	}
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
		if r[i] == -1 {
			continue
		}
		w := (r[i] - len(r)/2)
		if w < 0 {
			w = -w
			w += len(r) - r[i]
		} else {
			w++
		}
		rx := (i - h[r[i]]) * w
		if rx < 0 {
			res += -rx
		} else {
			res += rx
		}
	}
	return res
}

func (v *Vote) DumbScore(r []int) int {
	if v.C["0"] == r[0] {
		return 0
	}
	return 1
}

func (v *Vote) PeakValue() int {
	a := v.C["0"] - 1
	b := v.C["0"] + 1
	for i := 1; i < len(v.C); i++ {
		is := strconv.Itoa(i)
		val := v.C[is]
		if a == val {
			a--
		} else if b == val {
			b++
		} else {
			return -1
		}
	}
	p := v.C["0"]
	v.Peak = &p
	return p
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

// wins returns list of candidates that 'a' beats in a 1-1
func (e *Election) wins(a int) []int {
	res := make([]int, 0)

	for i := 0; i < e.N; i++ {
		if i == a {
			continue
		}
		if e.cmp(a, i) < 0 {
			res = append(res, i)
		}
	}

	return res
}

// Peak returns the peak of an election
func (e *Election) Peak() int {
	res := -1
	for i, v := range e.V {
		if i == 0 {
			res = v.PeakValue()
		}
		if res != v.PeakValue() {
			return -1
		}
	}
	return res
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
	res += fmt.Sprint("weight,")
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

func (e *Election) Plurality() []int {
	score := make([]int, e.N)
	for _, v := range e.V {
		for j := 0; j < len(v.C); j++ {
			idx := v.C[strconv.Itoa(j)]
			for len(score) <= idx {
				score = append(score, 0)
			}
			score[idx] += len(v.C) - j
		}
	}
	return score
}

type Manipulation struct {
	VoteIndex int    `json:"vote_index"`
	OrigVote  *Vote  `json:"orig_vote"`
	NewVote   *Vote  `json:"new_vote"`
	OrigScore int    `json:"orig_score"`
	NewScore  int    `json:"new_score"`
	OrigTally []int  `json:"orig_tally"`
	NewTally  []int  `json:"new_tally"`
	TallyType string `json:"type"`
}

// FindManipulation returns an available manipulation if one is found. The first manipulation found is returned
func (e *Election) FindManipulation(t Tallier) *Manipulation {
	perms := Perms(e.N)
	origTally := t.Tally(e)
	for i, v := range e.V {
		origVote := make(map[string]int)
		for k, x := range v.C {
			origVote[k] = x
		}
		origScore := v.Score(origTally)
		for _, p := range perms {
			for j := 0; j < e.N; j++ {
				v.C[strconv.Itoa(j)] = p[j]
			}
			newTally := t.Tally(e)
			newVotes := make(map[string]int)
			for k, x := range v.C {
				newVotes[k] = x
			}
			for k, x := range origVote {
				v.C[k] = x
			}
			newScore := v.Score(newTally)

			if newScore < origScore {
				return &Manipulation{
					VoteIndex: i,
					OrigVote:  v,
					OrigTally: origTally,
					NewTally:  newTally,
					NewVote: &Vote{
						C:    newVotes,
						W:    v.W,
						Peak: v.Peak,
					},
					OrigScore: origScore,
					NewScore:  newScore,
					TallyType: t.Key(),
				}
			}
		}
	}

	return nil
}

func (e *Election) Copy() *Election {
	m := make(map[string]string)
	for k, v := range e.M {
		m[k] = v
	}
	v := make([]*Vote, len(e.V))
	for i := 0; i < len(e.V); i++ {
		v[i] = e.V[i]
	}
	r := make([]int, len(e.R))
	for i := 0; i < len(e.R); i++ {
		r[i] = e.R[i]
	}
	return &Election{
		N: e.N,
		F: e.F,
		P: e.P,
		C: e.C,
		R: r,
		M: m,
		V: v,
	}
}

func (v *Vote) Copy() *Vote {
	c := make(map[string]int)
	for k, x := range v.C {
		c[k] = x
	}
	return &Vote{
		C:    c,
		W:    v.W,
		D:    v.D,
		Peak: v.Peak,
	}
}

func (e *Election) RemoveCandidate(k int) *Election {
	res := e.Copy()
	res.N--
	//TODO: care about pref, peak and condorcet values?
	res.F = nil
	res.C = nil
	res.P = nil
	if len(res.R) != 0 {
		res.R = RemoveAt(k, res.R)
	}
	for x, v := range res.M {
		if v == strconv.Itoa(k) {
			delete(res.M, x)
			break
		}
	}
	for i := 0; i < len(res.V); i++ {
		res.V[i] = res.V[i].RemoveCandidate(k)
	}

	return res
}

func (e *Election) Votes() int {
	res := 0
	for _, v := range e.V {
		res += v.W
	}
	return res
}

func (v *Vote) RemoveCandidate(cand int) *Vote {
	res := v.Copy()
	k := v.Ranki(cand) //get rank of candidate
	for i := k + 1; i < len(res.C); i++ {
		res.C[strconv.Itoa(i-1)] = res.C[strconv.Itoa(i)]
	}
	delete(res.C, strconv.Itoa(len(res.C)-1)) //remove last

	for _k, v := range res.C {
		if v > cand {
			res.C[_k]--
		}
	}
	return res
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

func (v *Vote) Ranki(k int) int {
	for key, v := range v.C {
		if v == k {
			ret, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}
			return ret
		}
	}
	log.Println(k, v.C)
	panic("Not Found")
}

//who comes after candidate k? -1 for none
func (v *Vote) after(k int) int {
	r := v.Ranki(k)
	if r+1 == len(v.C) {
		return -1
	}
	return v.C[strconv.Itoa(r+1)]
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

func (e *Election) CandVotes() VoteM {
	res := make(VoteM, e.N, e.N)
	for i := 0; i < e.N; i++ {
		res[i] = make(VoteAr, 0)
	}

	for _, v := range e.V {
		idx := v.C["0"]
		for x := 0; x < v.W; x++ {
			//fmt.Println(idx, e.N, len(v.C))
			res[idx] = append(res[idx], v.Copy())
		}
	}
	return res
}
