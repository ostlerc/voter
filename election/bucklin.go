package election

import "strconv"

type TallyBucklin struct{}

func init() {
	RegisterTally(&TallyBucklin{})
}

func (*TallyBucklin) Tally(e *Election) []int {
	res := make([]int, e.N)
	for i := 0; i < len(res); i++ {
		res[i] = -1
	}
	for i := 1; i < len(res)+1; i++ {
		if b := e.Bucklin(i); b != -1 {
			res[0] = b
			res[1] = i - 1
			break
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

	for i := 0; i < e.N; i++ {
		if res[i] > m {
			return i
		}
	}
	return -1
}
