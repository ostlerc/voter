package election

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type TallySTV struct {
	Droop  int
	e      *Election
	Score  Ints
	Ignore Ints
	Votes  []Ints
	Start  int
	End    int
	res    Ints
}

func tojson(i interface{}) string {
	dat, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func init() {
	RegisterTally(&TallySTV{})
}

func (t *TallySTV) Tally(e *Election) []int {
	t.e = e
	return t.STV()
}

func (*TallySTV) Key() string {
	return "stv"
}

// Returns STV result
func (t *TallySTV) STV() []int {
	e := t.e
	t.Droop = e.Votes()/e.N + 1
	fmt.Println("droop=", t.Droop, e.Votes(), e.N)
	t.res = make([]int, e.N)
	t.Ignore = make([]int, 0)
	t.Start = 0
	t.End = e.N
	t.Votes = e.CandVotes()
	t.Score = make([]int, e.N)
	copy(t.Score, t.Votes[0])
	for t.Start != t.End {
		t.step()
	}
	return t.res
}

func (e *Election) CandVotes() []Ints {
	res := make([]Ints, e.N, e.N)
	for i := 0; i < e.N; i++ {
		res[i] = make([]int, e.N, e.N)
		for _, v := range e.V {
			res[i][v.C[strconv.Itoa(i)]] += v.W
		}
	}
	return res
}

//distributes from i to scores.
func (t *TallySTV) step() {
	fmt.Println(tojson(t))
	c := t.Score.Maxi(t.Ignore)
	if t.Score[c] >= t.Droop {
		fmt.Println("Winner", c, t.Score[c])
		t.Start++
	} else {
		c = t.Score.Mini(t.Ignore)
		fmt.Println("loser", c, t.Score[c])
		t.End--
	}
	t.distribute(c)
	t.Ignore = append(t.Ignore, c)
	t.res = append(t.res, c)
}

func (t *TallySTV) distribute(i int) {
	ratio := float64(1)
	if t.Score[i] >= t.Droop {
		ratio = float64(t.Score[i]-t.Droop) / float64(t.Score[i])
		fmt.Println("Ratio:", ratio, t.Score[i], t.Droop)
	}

	a := t.after(i)
	for k := 0; k < len(t.Score); k++ {
		t.Score[k] += int(ratio * float64(a[k]))
	}
}

func (t *TallySTV) after(k int) Ints {
	res := make(Ints, t.e.N, t.e.N)
	for _, v := range t.e.V {
		if !t.first(v, k) {
			continue
		}

		if x := v.after(k); x != -1 {
			res[x] += v.W
		}
	}
	return res
}

func (t *TallySTV) first(v *Vote, k int) bool {
	for i := 0; i < len(v.C); i++ {
		x := v.C[strconv.Itoa(i)]
		if !Contains(x, t.Ignore) {
			return x == k
		}
	}

	return false
}
