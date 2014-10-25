package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/ostlerc/voter/election"
)

// Voter returns an vote with n candidates
type Voter interface {
	Vote(n int) *election.Vote
}

type PreferVoter struct {
	f    *election.IntPair
	peak bool
}
type RandVoter struct{}
type PeakVoter struct{}

func (*RandVoter) Vote(n int) *election.Vote {
	return vote(n)
}

func (*PeakVoter) Vote(n int) *election.Vote {
	p := election.R.Intn(n)
	e := peakvote(p, n)
	e.Peak = &p
	return e
}

func vote(n int) *election.Vote {
	vote := election.NewVote()
	for i := 0; i < n; i++ {
		r := election.Strn(n)
		for vote.Contains(r) {
			r = election.Strn(n)
		}
		vote.C[r] = i
	}
	return vote
}

func peakvote(p, n int) *election.Vote {
	res := election.NewVote()
	res.C["0"] = p

	l := p - 1
	r := p + 1

	for i := 1; l != -1 || r != n; i++ {
		istr := strconv.Itoa(i)
		if l != -1 && r != n {
			if election.R.Intn(2) == 0 {
				res.C[istr] = l
				l--
			} else {
				res.C[istr] = r
				r++
			}
		} else if l == -1 {
			for ; r != n; r++ {
				res.C[istr] = r
				i++
				istr = strconv.Itoa(i)
			}
			return res
		} else {
			for ; l != -1; l-- {
				res.C[istr] = l
				i++
				istr = strconv.Itoa(i)
			}
			return res
		}
	}
	return res
}

func (p *PreferVoter) Vote(n int) *election.Vote {
	var res *election.Vote

	if p.peak {
		pn := election.R.Intn(n)
		for pn == p.f.Second { //anything but that...
			pn = election.R.Intn(n)
		}
		res = peakvote(pn, n)
		res.Peak = &pn
	} else {
		res = vote(n)
	}

	res.Prefer(p.f)
	aidx := res.Rank(p.f.First)
	bidx := res.Rank(p.f.Second)
	if aidx > bidx {
		res.C[aidx] = p.f.Second
		res.C[bidx] = p.f.First
	}
	return res
}

type csvVoter struct {
	r       *bufio.Reader
	weights []int
	M       map[string]string
	count   int
}

func csvElection(e *election.Election, r io.Reader) {
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	file := strings.Split(string(dat), "\n")
	weights := make([]int, 0)

	*Candidates = len(file) - 1
	*Votes = 0
	e.V = make([]*election.Vote, 0)

	for y, row := range file {
		csv := strings.Split(row, ",")
		if len(csv) < 2 {
			*Candidates--
			continue
		}
		for x, col := range csv {
			if y == 0 { //handle weights
				if x == 0 {
					continue //empty anyways
				}
				iw, err := strconv.Atoi(col)
				if err != nil {
					log.Fatal(err)
				}
				*Votes += iw
				weights = append(weights, iw)
				e.V = append(e.V, &election.Vote{C: make(map[string]int)})
				continue
			}
			if x == 0 { //handle name map
				e.M[strconv.Itoa(y-1)] = col
				continue
			}

			icol, err := strconv.Atoi(col)
			if err != nil {
				log.Fatal(err)
			}
			e.V[x-1].C[strconv.Itoa(y-1)] = icol - 1

		}
	}

	for j, w := range weights {
		for i := 1; i < w; i++ { //add vote w-1 times
			e.V = append(e.V, e.V[j])
		}
	}

	e.N = *Candidates
	if *pref {
		e.F = election.NewPref(*Candidates)
	}
}
