package election

import "strconv"

type TallySTV struct {
	Droop  int
	e      *Election
	Score  VoteM
	Ignore Ints
	Start  int
	End    int
	res    Ints
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
	//fmt.Println("droop=", t.Droop, e.Votes(), e.N)
	t.res = make([]int, t.e.N)
	t.Ignore = make([]int, 0)
	t.Start = 0
	t.End = e.N
	t.Score = e.CandVotes()
	for t.Start != t.End {
		//fmt.Println("Start", t.Start, "End", t.End)
		t.step()
	}
	return t.res
}

//distributes from i to scores.
func (t *TallySTV) step() {
	c := t.Score.Maxi(t.Ignore)
	//fmt.Println("score=", tojson(t.Score.Lens()))
	//fmt.Println("ignore=", tojson(t.Ignore))
	if len(t.Score[c]) >= t.Droop {
		//fmt.Println("Winner", c, t.Score.Lens()[c])
		t.res[t.Start] = c
		t.Start++
	} else {
		c = t.Score.Mini(t.Ignore)
		//fmt.Println("loser", c, t.Score.Lens()[c])
		t.End--
		t.res[t.End] = c
	}
	t.Ignore = append(t.Ignore, c)
	//fmt.Println("before", t.Score.Lens(), t.e.N)
	t.distribute(c)
	//fmt.Println("after", t.Score.Lens(), t.e.N)
}

func (t *TallySTV) distribute(i int) {
	ratio := float64(1)
	if len(t.Score[i]) >= t.Droop {
		ratio = float64(len(t.Score[i])-t.Droop) / float64(len(t.Score[i]))
		//fmt.Println("Ratio:", ratio, len(t.Score[i])-t.Droop, t.Score.Lens()[i])
	}

	afterVotes := t.after(i)
	for k := 0; k < t.e.N; k++ {
		if k == i {
			continue //don't distribute to ourselves
		}

		distributableVotes := afterVotes[k]
		c := int(ratio * float64(len(distributableVotes)))
		//append first c votes
		t.Score[k] = append(t.Score[k], distributableVotes[0:c]...)
	}
}

// returns distribution list of votes for each candidate given ignore set
func (t *TallySTV) after(k int) VoteM {
	res := make(VoteM, t.e.N, t.e.N)
	for i := 0; i < t.e.N; i++ {
		res[i] = make(VoteAr, 0)
	}
	for _, v := range t.Score[k] {
		if i := t.first(v); i != -1 {
			//fmt.Println(i, t.e.N)
			res[i] = append(res[i], v)
		}
	}
	return res
}

// returns first choice taking into account Ignore list and who it is they voted for
func (t *TallySTV) first(v *Vote) int {
	for i := 0; i < len(v.C); i++ {
		x := v.C[strconv.Itoa(i)]
		if !Contains(x, t.Ignore) {
			return x
		}
	}

	return -1
}
