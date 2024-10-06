package internalhttp

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelloHandler(t *testing.T) {
	server := NewServer(nil)
	writer := httptest.NewRecorder()
	server.helloHandler(writer, nil)

	response := writer.Result()
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	expected := string(data)
	actual := "Hello, World!"
	require.Equal(t, expected, actual)
}
