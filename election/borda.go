package election

type TallyBorda struct{}

func init() {
	RegisterTally(&TallyBorda{})
}

func (*TallyBorda) Tally(e *Election) []int {
	return e.Borda()
}

func (*TallyBorda) Key() string {
	return "borda"
}

// Returns Borda result
func (e *Election) Borda() []int {
	score := e.Plurality()
	res := make([]int, 0)
	for i := 0; i < e.N; i++ {
		max := score[0]
		maxi := 0
		for j := 1; j < e.N; j++ {
			if score[j] > max {
				max = score[j]
				maxi = j
			}
		}
		res = append(res, maxi)
		score[maxi] = -1
	}
	return res
}
