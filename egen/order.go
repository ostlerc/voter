package main

import "github.com/ostlerc/voter/election"

// Orderer returns an ordering of candidates
type Orderer interface {
	Order(int) []int
}

type PreferOrder struct {
	f    *election.Favorite
	peak bool
}
type RandOrder struct{}

func (*RandOrder) Order(n int) []int {
	return order(n)
}

func order(n int) []int {
	res := make([]int, 0)
	for i := 0; i < n; i++ {
		v := election.R.Intn(n)
		for election.Contains(v, res) {
			v = election.R.Intn(n)
		}
		res = append(res, v)
	}
	return res
}

func peakorder(p, n int) []int {
	res := make([]int, 0)
	l := p - 1
	r := p + 1
	res = append(res, p)
	for l != -1 || r != n {
		if l != -1 && r != n {
			if election.R.Intn(2) == 0 {
				res = append(res, l)
				l--
			} else {
				res = append(res, r)
				r++
			}
		} else if l == -1 {
			for ; r != n; r++ {
				res = append(res, r)
			}
			return res
		} else {
			for ; l != -1; l-- {
				res = append(res, l)
			}
			return res
		}
	}
	return res
}

func (p *PreferOrder) Order(n int) []int {
	var res []int

	if p.peak {
		pn := election.R.Intn(n)
		for pn == p.f.B { //anything but that...
			pn = election.R.Intn(n)
		}
		res = peakorder(pn, n)
	} else {
		res = order(n)
	}

	aidx := election.Index(p.f.A, res)
	bidx := election.Index(p.f.B, res)
	if aidx > bidx {
		res[aidx] = p.f.B
		res[bidx] = p.f.A
	}
	return res
}
