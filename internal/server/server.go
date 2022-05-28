package server

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"word-of-wisdom/internal/challenge"
	"word-of-wisdom/internal/protocol"
)

// difficulty is hardcoded, but can be made to vary based on the server load, or set higher for fraud clients.
const challengeDifficulty = 2

type handler struct {
	quotes []Quote
}

func Run(ctx context.Context, addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listen.Close()

	fmt.Printf("Server started on %s\n", addr)

	quotes, err := loadQuotes(quotesJson)
	if err != nil {
		return fmt.Errorf("failed to load quotes: %w", err)
	}

	handler := handler{quotes: quotes}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}

			if err := conn.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}

			fmt.Printf("Handling new connection: %s\n", conn.RemoteAddr())
			go func() {
				clientAddr := conn.RemoteAddr()
				defer fmt.Printf("Closing connection: %s\n", clientAddr)
				defer conn.Close()

				if shouldGiveClientChallenge() {
					solved, err := giveClientChallenge(conn, challengeDifficulty)
					if err != nil {
						fmt.Printf("Client solve challenge error %s: %s\n", clientAddr, err)
						return
					}
					if !solved {
						fmt.Printf("Client not solved challenge %s\n", clientAddr)
						return
					}
				}

				if err := handler.Handle(conn); err != nil {
					fmt.Printf("Error handling %s: %s\n", conn.RemoteAddr(), err)
				}
			}()
		}

	}()

	<-ctx.Done()
	return nil
}

func (s handler) Handle(conn io.ReadWriter) error {
	quote := s.quotes[randRange(0, len(s.quotes))]
	if err := protocol.WriteOK(conn, fmt.Sprintf("%s", quote)); err != nil {
		return fmt.Errorf("failed to write quote: %w", err)
	}

	return nil
}

// shouldGiveClientChallenge decides if client should be given a challenge
//
// This is a simple implementation that returns true with a probability of 0.5,
// but it can be extended to be smarter and decide based on the client's IP address or current server load.
//
func shouldGiveClientChallenge() bool {
	return rand.Intn(100) > 50
}

func giveClientChallenge(conn io.ReadWriter, difficulty int) (bool, error) {
	challengeToSolve := challenge.GenerateSHA256(difficulty)
	if err := protocol.WriteChallenge(conn, challengeToSolve); err != nil {
		return false, fmt.Errorf("failed to write task: %w", err)
	}

	result, err := protocol.ReadMsg(conn)
	if err != nil {
		return false, fmt.Errorf("failed to read answer: %w", err)
	}

	if !challenge.IsSolved(challengeToSolve, result.Payload) {
		_ = protocol.WriteErr(conn, "task not solved")
		return false, nil
	}

	return true, nil
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}
