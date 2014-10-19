package main

import (
	"github.com/ostlerc/voter/election"
)

// Orderer returns an ordering of candidates
type Orderer interface {
	Order(int) []int
}

type PreferOrder struct {
	f *election.Favorite
}
type RandOrder struct{}

func (*RandOrder) Order(n int) []int {
	return order(n)
}

func order(n int) []int {
	res := make([]int, 0)
	for i := 0; i < n; i++ {
		for {
			v := election.R.Intn(n)
			if !election.Contains(v, res) {
				res = append(res, v)
				break
			}
		}
	}
	return res
}

func (p *PreferOrder) Order(n int) []int {
	res := order(n)
	aidx := election.Index(p.f.A, res)
	bidx := election.Index(p.f.B, res)
	if aidx > bidx {
		res[aidx] = p.f.B
		res[bidx] = p.f.A
	}
	return res
}
