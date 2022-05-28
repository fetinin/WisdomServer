package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteOK(t *testing.T) {
	conn := bytes.Buffer{}

	err := WriteOK(&conn, "test message")

	require.NoError(t, err)
	require.Equal(t, "OK\ttest message\n", conn.String())
}

func TestReadMsg(t *testing.T) {
	conn := bytes.NewBufferString("OK\ttest message\n")

	msg, err := ReadMsg(conn)

	require.NoError(t, err)
	require.Equal(t, Msg{Cmd: "OK", Payload: "test message"}, msg)
}

func TestReadWriteMsg(t *testing.T) {
	testcases := []string{"\n", " ", "!12345", "Hello, world"}
	for _, tc := range testcases {
		t.Run(tc, func(t *testing.T) {
			conn := bytes.Buffer{}

			err := WriteOK(&conn, tc)
			if err := checkAllowedToSend(tc); err != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			msg, err := ReadMsg(&conn)
			require.NoError(t, err)
			require.Equal(t, Msg{Cmd: "OK", Payload: tc}, msg)
		})
	}
}

func FuzzWriteOK(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345", "\n", "\t"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, msg string) {
		conn := bytes.Buffer{}

		if checkAllowedToSend(msg) != nil {
			t.Skip("Skipping invalid input")
		}

		err := WriteOK(&conn, msg)
		require.NoError(t, err)

		msgBack, err := ReadMsg(&conn)
		require.NoError(t, err)

		require.Equal(t, Msg{Cmd: OK, Payload: msg}, msgBack)
	})
}
