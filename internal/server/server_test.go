package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"word-of-wisdom/internal/protocol"
)

func Test_handlerHandle__Success(t *testing.T) {
	// Setup
	sockClient, sockServer := newFakeSockets()
	s := handler{
		quotes: []Quote{{Text: "test", Author: "test"}},
	}

	// Execute
	go func() {
		err := s.Handle(&sockServer)
		require.NoError(t, err)
	}()

	// Assert
	// on connection, quote should be given
	resp, err := protocol.ReadMsg(&sockClient)
	require.NoError(t, err)
	require.Equal(t, "test -- test", resp.Payload)
}

type fakeTwoWaySocket struct {
	From          *bytes.Buffer
	To            *bytes.Buffer
	RespondedTo   chan struct{}
	RespondedFrom chan struct{}
}

func (f fakeTwoWaySocket) Write(p []byte) (n int, err error) {
	n, err = f.To.Write(p)
	f.RespondedTo <- struct{}{}
	return n, err
}

func (f fakeTwoWaySocket) Read(p []byte) (n int, err error) {
	<-f.RespondedFrom
	return f.From.Read(p)
}

func newFakeSockets() (fakeTwoWaySocket, fakeTwoWaySocket) {
	var buffer1 bytes.Buffer
	var buffer2 bytes.Buffer
	responded1 := make(chan struct{})
	responded2 := make(chan struct{})
	return fakeTwoWaySocket{
			From:          &buffer1,
			To:            &buffer2,
			RespondedTo:   responded2,
			RespondedFrom: responded1,
		}, fakeTwoWaySocket{
			From:          &buffer2,
			To:            &buffer1,
			RespondedTo:   responded1,
			RespondedFrom: responded2,
		}
}
