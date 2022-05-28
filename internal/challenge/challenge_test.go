package challenge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name      string
		challenge Challenge
		want      string
		wantErr   bool
	}{
		{
			name: "simple",
			challenge: Challenge{
				Type:       "SHA256",
				Difficulty: 1,
				Seed:       "NOT_SO_RANDOM",
			},
			want: "244",
		},
		{
			name: "medium",
			challenge: Challenge{
				Type:       "SHA256",
				Difficulty: 3,
				Seed:       "NOT_SO_RANDOM",
			},
			want: "6020355",
		},
		{
			name:      "unsupported type",
			challenge: Challenge{Type: "NOT_SUPPORTED"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Solve(context.Background(), tt.challenge)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestSolveDeadline checks that Solve respects context deadline.
func TestSolveDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	challenge := Challenge{
		Type:       "SHA256",
		Difficulty: 5,
		Seed:       "NOT_SO_RANDOM",
	}
	answer, err := Solve(ctx, challenge)
	require.Error(t, err)
	require.Equal(t, context.DeadlineExceeded, err)
	require.Empty(t, answer)
}

// run: go test -v -fuzz=FuzzSolve -fuzztime=30s ./internal/challenge
func FuzzSolve(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	ctx := context.Background()
	f.Fuzz(func(t *testing.T, orig string) {
		challenge := Challenge{
			Type:       "SHA256",
			Difficulty: 2,
			Seed:       orig,
		}
		answer, err := Solve(ctx, challenge)
		require.NoError(t, err)
		require.True(t, IsSolved(challenge, answer))
	})
}
