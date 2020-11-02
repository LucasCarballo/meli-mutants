package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Interface provides exported methods signatures
type Interface interface {
	Initialize()
	SetDNA([]string, bool)
	GetDNACounts() (int64, int64)
}

// Repository struct to inject redis client
type Repository struct {
	Client *redis.Client
}

var ctx = context.Background()

// Initialize counters in case they doesn't exist
func (repo Repository) Initialize() {
	if repo.Client == nil {
		repo.Client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}

	err := repo.Client.SetNX(ctx, "mutantCounter", 0, 0).Err()
	if err != nil {
		panic(err)
	}

	err = repo.Client.SetNX(ctx, "humanCounter", 0, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (repo Repository) setMutant(dna []byte) {
	counter, err := repo.Client.Get(ctx, "mutantCounter").Result()
	if err != nil {
		panic(err)
	}

	err = repo.Client.Set(ctx, "mutant:"+counter, dna, 0).Err()
	if err != nil {
		panic(err)
	}

	err = repo.Client.Incr(ctx, "mutantCounter").Err()
	if err != nil {
		panic(err)
	}
}

func (repo Repository) setHuman(dna []byte) {
	counter, err := repo.Client.Get(ctx, "humanCounter").Result()
	if err != nil {
		panic(err)
	}

	err = repo.Client.Set(ctx, "human:"+counter, dna, 0).Err()
	if err != nil {
		panic(err)
	}

	err = repo.Client.Incr(ctx, "humanCounter").Err()
	if err != nil {
		panic(err)
	}
}

// SetDNA if it does not exist.
func (repo Repository) SetDNA(dna []string, isMutant bool) {
	dnaString := strings.Join(dna, "")
	set, err := repo.Client.SetNX(ctx, dnaString, 0, 0).Result()
	if err != nil {
		panic(err)
	}

	if set == false {
		return
	}

	dnaJSON, err := json.Marshal(dna)
	if err != nil {
		panic(err)
	}

	if isMutant {
		repo.setMutant(dnaJSON)
	} else {
		repo.setHuman(dnaJSON)
	}
}

// GetDNACounts returns mutants and humans dna stored
func (repo Repository) GetDNACounts() (mutantCount int64, humanCount int64) {
	mutantCountString, err := repo.Client.Get(ctx, "mutantCounter").Result()
	if err != nil {
		panic(err)
	}

	mutantCount, err = strconv.ParseInt(mutantCountString, 10, 64)
	if err != nil {
		panic(err)
	}

	humanCountString, err := repo.Client.Get(ctx, "humanCounter").Result()
	if err != nil {
		panic(err)
	}

	humanCount, err = strconv.ParseInt(humanCountString, 10, 64)
	if err != nil {
		panic(err)
	}

	return mutantCount, humanCount
}
