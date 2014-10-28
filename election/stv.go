package election

import "strconv"

type TallySTV struct{}

func init() {
	RegisterTally(&TallySTV{})
}

func (*TallySTV) Tally(e *Election) []int {
	return e.STV()
}

func (*TallySTV) Key() string {
	return "stv"
}

// Returns STV result
func (e *Election) STV() []int {
	res := make([]int, e.N)
	ignore := make([]int, 0)
	for i := 0; i < e.N; i++ {
		score := e.stv(ignore)
		min := -1
		mini := 0
		for j := 0; j < e.N; j++ {
			if Contains(j, ignore) {
				continue
			}
			if min == -1 || score[j] < min {
				min = score[j]
				mini = j
			}
		}
		res[e.N-i-1] = mini
		ignore = append(ignore, mini)
	}
	return res
}

func (e *Election) stv(ignore []int) []int {
	res := make([]int, e.N)
	for _, v := range e.V {
		for i := 0; i < len(v.C); i++ {
			cand := v.C[strconv.Itoa(i)]
			if !Contains(cand, ignore) {
				res[cand]++
				break
			}
		}
	}
	return res
}
