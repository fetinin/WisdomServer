package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"word-of-wisdom/internal/challenge"
	"word-of-wisdom/internal/protocol"
)

func Run(ctx context.Context, serverAddr string) (string, error) {
	conn, err := net.DialTimeout("tcp", serverAddr, time.Second*5)
	if err != nil {
		return "", fmt.Errorf("failed to connect to server: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(time.Second * 10))

	msg, err := protocol.ReadMsg(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read challenge: %w", err)
	}

	if msg.Cmd == protocol.Challenge {
		challengeMsg, err := protocol.ParseChallenge(msg.Payload)
		if err != nil {
			_ = protocol.WriteErr(conn, fmt.Sprintf("failed to parse challenge: %s", err))
			return "", fmt.Errorf("failed to parse challenge: %w", err)
		}

		result, err := challenge.Solve(ctx, challengeMsg)
		if err != nil {
			return "", fmt.Errorf("failed to solve challenge: %w", err)
		}

		if err := protocol.WriteOK(conn, result); err != nil {
			return "", fmt.Errorf("failed to write challenge result: %w", err)
		}

		msg, err = protocol.ReadMsg(conn)
		if err != nil {
			return "", fmt.Errorf("failed to read challenge result: %w", err)
		}
	}

	return msg.Payload, nil
}
