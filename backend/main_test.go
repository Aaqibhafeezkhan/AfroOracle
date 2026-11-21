package main

import (
	"errors"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestHTTPServerConfiguresRequestTimeouts(t *testing.T) {
	handler := http.NewServeMux()
	srv := newHTTPServer(":8000", handler)

	if srv.Handler != handler {
		t.Fatal("newHTTPServer did not assign the supplied handler")
	}
	if srv.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout = %s, want 5s", srv.ReadHeaderTimeout)
	}
	if srv.ReadTimeout != 15*time.Second {
		t.Errorf("ReadTimeout = %s, want 15s", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 30*time.Second {
		t.Errorf("WriteTimeout = %s, want 30s", srv.WriteTimeout)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout = %s, want 60s", srv.IdleTimeout)
	}
}

func TestHTTPServerFailsWhenPortIsBusy(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve test port: %v", err)
	}
	defer listener.Close()

	srv := newHTTPServer(listener.Addr().String(), http.NewServeMux())
	err = srv.ListenAndServe()
	if err == nil {
		t.Fatal("ListenAndServe succeeded on an occupied port")
	}
	if !errors.Is(err, syscall.EADDRINUSE) {
		t.Fatalf("ListenAndServe error = %v, want address already in use", err)
	}
}
