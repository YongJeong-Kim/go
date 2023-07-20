package gapi

import (
	"os"
	"testing"
)

func NewTestServer() *Server {
	return NewServer()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
