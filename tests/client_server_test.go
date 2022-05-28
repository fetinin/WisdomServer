package tests

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"word-of-wisdom/internal/client"
	"word-of-wisdom/internal/server"
)

// TestClientServer quickly allows to quickly check that server and client communicate successfully.
func TestClientServer(t *testing.T) {
	rand.Seed(0)
	ctx := context.Background()
	addr := "localhost:9001"
	go func() {
		if err := server.Run(ctx, addr); err != nil {
			t.Errorf("Server failed: %s", err)
		}
	}()

	reply, err := client.Run(ctx, addr)
	require.NoError(t, err)
	require.Equal(
		t,
		"There are no accidents... there is only some purpose that we haven't yet understood. -- Deepak Chopra",
		reply,
	)
}
