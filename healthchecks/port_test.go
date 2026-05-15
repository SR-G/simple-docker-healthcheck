package healthchecks

import (
	"net"
	"testing"
)

func TestPortHealthCheckExecute(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	defer listener.Close()

	t.Run("success", func(t *testing.T) {
		checker := &PortHealthCheck{Hostname: "127.0.0.1", Port: addr.Port}
		ok, err := checker.Execute()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !ok {
			t.Fatal("expected health check to succeed")
		}
	})

	t.Run("failure", func(t *testing.T) {
		closedListener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to start second listener: %v", err)
		}
		closedAddr := closedListener.Addr().(*net.TCPAddr)
		closedListener.Close()

		checker := &PortHealthCheck{Hostname: "127.0.0.1", Port: closedAddr.Port}
		ok, err := checker.Execute()
		if err == nil {
			t.Fatal("expected error when connecting to closed port")
		}
		if ok {
			t.Fatal("expected health check to fail when port is closed")
		}
	})
}
