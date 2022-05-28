package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"word-of-wisdom/internal/challenge"
)

const (
	separator = "\t"
	endMsg    = '\n'
)

type Command string

const (
	OK        Command = "OK"
	Err       Command = "ERR"
	Challenge Command = "CHALLENGE"
)

// Msg represents message that client and server communicate with
// Example: "OK\tpayload"
type Msg struct {
	Cmd     Command
	Payload string
}

func ReadMsg(conn io.Reader) (Msg, error) {
	message, err := bufio.NewReader(conn).ReadString(endMsg)
	if err != nil {
		return Msg{}, fmt.Errorf("failed to read message: %w", err)
	}

	message = strings.TrimRight(message, string(endMsg))
	cmd, payload, _ := strings.Cut(message, separator)

	if Command(cmd) == Err {
		return Msg{}, fmt.Errorf("error response: %s", payload)
	}

	return Msg{Cmd: Command(cmd), Payload: payload}, nil
}

func WriteErr(conn io.Writer, msg string) error {
	return writeMsg(conn, Err, msg)
}

func WriteChallenge(conn io.Writer, challenge challenge.Challenge) error {
	payload := fmt.Sprintf("%s:%d:%s", challenge.Type, challenge.Difficulty, challenge.Seed)
	return writeMsg(conn, Challenge, payload)
}

func ParseChallenge(payload string) (challenge.Challenge, error) {
	parts := strings.SplitN(payload, ":", 3)
	if len(parts) != 3 {
		return challenge.Challenge{}, fmt.Errorf("invalid challenge payload")
	}

	difficulty, _ := strconv.Atoi(parts[1])
	return challenge.Challenge{
		Type:       parts[0],
		Difficulty: difficulty,
		Seed:       parts[2],
	}, nil
}

func WriteOK(conn io.Writer, msg string) error {
	return writeMsg(conn, OK, msg)
}

func writeMsg(conn io.Writer, cmd Command, msg string) error {
	if err := checkAllowedToSend(msg); err != nil {
		return err
	}

	_, err := bytes.NewBufferString(string(cmd) + separator + msg + string(endMsg)).WriteTo(conn)
	return err
}

func checkAllowedToSend(msg string) error {
	if strings.Contains(msg, string(endMsg)) {
		return fmt.Errorf("message contains invalid character: %s", string(endMsg))
	}

	return nil
}
