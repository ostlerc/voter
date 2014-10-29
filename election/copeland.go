package election

type TallyCopeland struct{}

func init() {
	RegisterTally(&TallyCopeland{})
}

func (*TallyCopeland) Tally(e *Election) []int {
	return e.Copeland()
}

func (*TallyCopeland) Key() string {
	return "copeland"
}

// Returns Copeland result
func (e *Election) Copeland() []int {
	//gather first order scores
	fo := make([]int, e.N)
	for i := 0; i < e.N; i++ {
		for j := i + 1; j < e.N; j++ {
			r := e.cmp(i, j)
			if r < 0 {
				fo[i]++
				fo[j]--
			} else if r > 0 {
				fo[j]++
				fo[i]--
			}
		}
	}

	max := 0
	idxs := make([]int, 0)
	for i := 0; i < len(fo); i++ {
		if fo[i] > max {
			max = fo[i]
			idxs = []int{i}
		} else if fo[i] == max {
			idxs = append(idxs, i)
		}
	}
	if len(idxs) == 1 {
		for len(idxs) != len(e.V[0].C) {
			idxs = append(idxs, -1)
		}
		return idxs
	}

	//tie breaker
	res := make([]int, len(idxs))
	for i := 0; i < len(idxs); i++ {
		for _, v := range e.wins(idxs[i]) {
			res[i] += fo[v]
		}
	}

	max = 0
	idxs2 := make([]int, 0)
	for i := 0; i < len(res); i++ {
		if res[i] > max {
			max = res[i]
			idxs2 = []int{i}
		} else if res[i] == max {
			idxs2 = append(idxs2, i)
		}
	}
	for len(idxs2) != len(e.V[0].C) {
		idxs2 = append(idxs2, -1)
	}
	return idxs2
}
