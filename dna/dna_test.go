package dna

import (
	"strings"
	"testing"
)

var isMutantTests = []struct {
	dna      []string
	expected bool
}{
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCATA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGATGG", "CCCATA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTCC", "TTATGT", "AGATGG", "CCCATA", "TCACTG"}, false},
	{[]string{"ATGCGA", "CAGTCC", "TGATGT", "AGATGG", "CACATA", "TCACTG"}, false},
}

var dnaService Interface

func init() {
	dnaService = Service{}
}

func TestIsMutant(t *testing.T) {

	for _, tt := range isMutantTests {
		t.Run(strings.Join(tt.dna, ""), func(t *testing.T) {

			got := dnaService.IsMutant(tt.dna)
			if got != tt.expected {
				t.Errorf("Expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}
