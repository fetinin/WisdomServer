package challenge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Challenge struct {
	Type       string
	Difficulty int
	Seed       string
}

const SHA256 = "SHA256"

func GenerateSHA256(difficulty int) Challenge {
	seed := rand.Int63()
	return Challenge{Type: SHA256, Difficulty: difficulty, Seed: fmt.Sprintf("%d", seed)}
}

func IsSolved(challenge Challenge, answer string) bool {
	expectedPrefix := make([]byte, challenge.Difficulty)
	for i := range expectedPrefix {
		expectedPrefix[i] = '0'
	}
	res := sha256.Sum256([]byte(challenge.Seed + answer))
	return bytes.Compare(res[:challenge.Difficulty], expectedPrefix) == 0
}

func Solve(ctx context.Context, challenge Challenge) (string, error) {
	if challenge.Type != SHA256 {
		return "", fmt.Errorf("unsupported challenge type: %s", challenge.Type)
	}

	expectedPrefix := make([]byte, challenge.Difficulty)
	for i := range expectedPrefix {
		expectedPrefix[i] = '0'
	}

	var nonce int32 = -1
	for nonce < math.MaxInt32 {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		nonce++
		res := sha256.Sum256([]byte(fmt.Sprintf("%s%d", challenge.Seed, nonce)))
		if eq := bytes.Compare(res[:challenge.Difficulty], expectedPrefix); eq == 0 {
			break
		}
	}

	return strconv.Itoa(int(nonce)), nil
}
