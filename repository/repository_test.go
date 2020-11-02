package repository

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

var s *miniredis.Miniredis
var repo Interface

var repositoryTests = []struct {
	dna      []string
	isMutant bool
}{
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCATA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTGC", "TTATGT", "AGATGG", "CCCATA", "TCACTG"}, true},
	{[]string{"ATGCGA", "CAGTCC", "TTATGT", "AGATGG", "CCCATA", "TCACTG"}, false},
	{[]string{"ATGCGA", "CAGTCC", "TGATGT", "AGATGG", "CACATA", "TCACTG"}, false},
}

func TestIsMutant(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := Repository{
		Client: redis.NewClient(&redis.Options{
			Addr:     s.Addr(),
			Password: "",
			DB:       0,
		}),
	}
	repo.Initialize()

	mutantCounter, err := s.Get("mutantCounter")
	if err != nil {
		panic(err)
	}

	humanCounter, err := s.Get("humanCounter")
	if err != nil {
		panic(err)
	}

	if mutantCounter != "0" {
		t.Errorf("Expected: %v, got: %v", "0", mutantCounter)
	}

	if humanCounter != "0" {
		t.Errorf("Expected: %v, got: %v", "0", humanCounter)
	}
}

func TestSetDNA(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := Repository{
		Client: redis.NewClient(&redis.Options{
			Addr:     s.Addr(),
			Password: "",
			DB:       0,
		}),
	}
	repo.Initialize()

	for _, tt := range repositoryTests {
		t.Run(strings.Join(tt.dna, ""), func(t *testing.T) {
			mutantCount, humanCount := repo.GetDNACounts()

			repo.SetDNA(tt.dna, tt.isMutant)

			updatedMutantCount, updatedHumanCount := repo.GetDNACounts()

			var dnaKey string
			if tt.isMutant {
				dnaKey = "mutant:" + strconv.FormatInt(mutantCount, 10)
			} else {
				dnaKey = "human:" + strconv.FormatInt(humanCount, 10)
			}
			savedDna, err := s.Get(dnaKey)
			if err != nil {
				panic(err)
			}

			dnaJSON, err := json.Marshal(tt.dna)
			if err != nil {
				panic(err)
			}

			if tt.isMutant {
				if updatedMutantCount != mutantCount+1 {
					t.Errorf("Expected: %v, got: %v", mutantCount+1, updatedMutantCount)
				}

				if updatedHumanCount != humanCount {
					t.Errorf("Expected: %v, got: %v", humanCount, updatedHumanCount)
				}
			} else {
				if updatedHumanCount != humanCount+1 {
					t.Errorf("Expected: %v, got: %v", humanCount+1, updatedHumanCount)
				}

				if updatedMutantCount != mutantCount {
					t.Errorf("Expected: %v, got: %v", mutantCount, updatedMutantCount)
				}
			}

			if savedDna != string(dnaJSON) {
				t.Errorf("Expected: %v, got: %v", string(dnaJSON), savedDna)
			}
		})
	}
}

func TestPreventDuplication(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := Repository{
		Client: redis.NewClient(&redis.Options{
			Addr:     s.Addr(),
			Password: "",
			DB:       0,
		}),
	}
	repo.Initialize()

	mutantCount, humanCount := repo.GetDNACounts()

	dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
	repo.SetDNA(dna, true)
	repo.SetDNA(dna, true)

	updatedMutantCount, updatedHumanCount := repo.GetDNACounts()

	if mutantCount+1 != updatedMutantCount {
		t.Errorf("Expected: %v, got: %v", mutantCount+1, updatedMutantCount)
	}

	if humanCount != updatedHumanCount {
		t.Errorf("Expected: %v, got: %v", humanCount, updatedHumanCount)
	}
}
