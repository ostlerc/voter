package election

import "strconv"

type TallyBucklin struct{}

func init() {
	RegisterTally(&TallyBucklin{})
}

func (*TallyBucklin) Tally(e *Election) []int {
	res := make([]int, len(e.V[0].C))
	for i := 0; i < len(res); i++ {
		res[i] = -1
	}
	for i := 0; i < len(res); i++ {
		if b := e.Bucklin(i); b != -1 {
			res[0] = b
			res[1] = i + 1
		}
	}
	return res
}

func (*TallyBucklin) Key() string {
	return "bucklin"
}

// Returns bucklin result with value k
func (e *Election) Bucklin(k int) int {
	res := make(map[int]int)
	total := 0
	for _, v := range e.V {
		for i := 0; i < k; i++ {
			res[v.C[strconv.Itoa(i)]] += v.W
		}
		total += v.W
	}
	m := total / 2
	if total%2 != 0 {
		m++
	}

	for i := 0; i < k; i++ {
		if res[i] >= m {
			return i
		}
	}
	return -1
}
