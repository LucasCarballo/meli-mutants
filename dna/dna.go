package dna

var x []int = []int{-1, -1, -1, 0, 0, 1, 1, 1}
var y []int = []int{-1, 0, 1, -1, 1, -1, 0, 1}

// Interface provides exported methods signatures
type Interface interface {
	IsMutant([]string) bool
}

// Service will provide isMutant to check dna arrays
type Service struct{}

func search2D(grid []string, row int, col int, word string) bool {
	if grid[row][col] != word[0] {
		return false
	}

	length := len(word)
	if len(grid) == 0 {
		panic("grid can't be empty")
	}

	C, R := len(grid), len(grid[0])

	for dir := 0; dir < 8; dir++ {
		k, rd, cd := 1, (row + x[dir]), (col + y[dir])

		for k = 1; k < length; k++ {
			if rd >= R || rd < 0 || cd >= C || cd < 0 {
				break
			}

			if grid[rd][cd] != word[k] {
				break
			}

			rd += x[dir]
			cd += y[dir]
		}

		if k == length {
			return true
		}
	}

	return false
}

func searchPattern(grid []string, word string) bool {
	if len(grid) == 0 {
		panic("grid can't be empty")
	}

	C, R := len(grid), len(grid[0])

	for row := 0; row < R; row++ {
		for col := 0; col < C; col++ {
			if search2D(grid, row, col, word) {
				return true
			}
		}
	}
	return false
}

// IsMutant returns true if the dna contains some of these patterns: "AAAA", "CCCC", "GGGG".
func (s Service) IsMutant(dna []string) bool {
	patterns := []string{"AAAA", "CCCC", "GGGG"}

	for i := 0; i < len(patterns); i++ {
		if searchPattern(dna, patterns[i]) {
			return true
		}
	}

	return false
}
