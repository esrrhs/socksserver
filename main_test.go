package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/esrrhs/gohome/network"
)

// startEchoServer starts a simple TCP echo server for testing targets.
func startEchoServer(t *testing.T) (net.Listener, string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start echo server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = io.Copy(c, c)
			}(conn)
		}
	}()

	return ln, ln.Addr().String()
}

// startTestSocksServer starts socksserver on a random port.
func startTestSocksServer(t *testing.T, user, password string) (*net.TCPListener, string) {
	t.Helper()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}

	go func() {
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				return
			}
			go process(conn, user, password)
		}
	}()

	return ln, ln.Addr().String()
}

func socks5Connect(proxyAddr, targetAddr, user, password string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", proxyAddr, 3*time.Second)
	if err != nil {
		return nil, err
	}

	// SOCKS5 Handshake client-side
	if user != "" || password != "" {
		// Method 0x02: USERNAME/PASSWORD
		if _, err := conn.Write([]byte{0x05, 0x01, 0x02}); err != nil {
			conn.Close()
			return nil, err
		}
		resp := make([]byte, 2)
		if _, err := io.ReadFull(conn, resp); err != nil {
			conn.Close()
			return nil, err
		}
		if resp[0] != 0x05 || resp[1] != 0x02 {
			conn.Close()
			return nil, fmt.Errorf("unexpected auth method response: %v", resp)
		}

		// Auth subnegotiation (RFC 1929)
		// VER(1) | ULEN(1) | UNAME | PLEN(1) | PASSWD
		buf := []byte{0x01, byte(len(user))}
		buf = append(buf, []byte(user)...)
		buf = append(buf, byte(len(password)))
		buf = append(buf, []byte(password)...)
		if _, err := conn.Write(buf); err != nil {
			conn.Close()
			return nil, err
		}
		authResp := make([]byte, 2)
		if _, err := io.ReadFull(conn, authResp); err != nil {
			conn.Close()
			return nil, err
		}
		if authResp[1] != 0x00 {
			conn.Close()
			return nil, fmt.Errorf("authentication failed: status %d", authResp[1])
		}
	} else {
		// Method 0x00: NO AUTHENTICATION REQUIRED
		if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
			conn.Close()
			return nil, err
		}
		resp := make([]byte, 2)
		if _, err := io.ReadFull(conn, resp); err != nil {
			conn.Close()
			return nil, err
		}
		if resp[0] != 0x05 || resp[1] != 0x00 {
			conn.Close()
			return nil, fmt.Errorf("unexpected method response: %v", resp)
		}
	}

	// SOCKS5 CONNECT Request
	host, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		conn.Close()
		return nil, err
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	ip := net.ParseIP(host)
	var req []byte
	if ip4 := ip.To4(); ip4 != nil {
		req = []byte{0x05, 0x01, 0x00, network.Socks5AtypIP4}
		req = append(req, ip4...)
	} else {
		req = []byte{0x05, 0x01, 0x00, network.Socks5AtypDomain, byte(len(host))}
		req = append(req, []byte(host)...)
	}
	req = append(req, byte(port>>8), byte(port&0xFF))

	if _, err := conn.Write(req); err != nil {
		conn.Close()
		return nil, err
	}

	// Read connect reply (10 bytes for IPv4)
	reply := make([]byte, 10)
	if _, err := io.ReadFull(conn, reply); err != nil {
		conn.Close()
		return nil, fmt.Errorf("read connect reply failed: %w", err)
	}
	if reply[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("connect failed with code %d", reply[1])
	}

	return conn, nil
}

func TestSocksServerNoAuth(t *testing.T) {
	echoLn, echoAddr := startEchoServer(t)
	defer echoLn.Close()

	proxyLn, proxyAddr := startTestSocksServer(t, "", "")
	defer proxyLn.Close()

	client, err := socks5Connect(proxyAddr, echoAddr, "", "")
	if err != nil {
		t.Fatalf("socks5Connect failed: %v", err)
	}
	defer client.Close()

	msg := []byte("hello modern socksserver")
	if _, err := client.Write(msg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(client, buf); err != nil {
		t.Fatalf("failed to read echoed message: %v", err)
	}

	if !bytes.Equal(buf, msg) {
		t.Fatalf("expected %q, got %q", string(msg), string(buf))
	}
}

func TestSocksServerWithAuthSuccess(t *testing.T) {
	echoLn, echoAddr := startEchoServer(t)
	defer echoLn.Close()

	proxyLn, proxyAddr := startTestSocksServer(t, "alice", "secret123")
	defer proxyLn.Close()

	client, err := socks5Connect(proxyAddr, echoAddr, "alice", "secret123")
	if err != nil {
		t.Fatalf("socks5Connect with auth failed: %v", err)
	}
	defer client.Close()

	msg := []byte("authenticated traffic check")
	if _, err := client.Write(msg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(client, buf); err != nil {
		t.Fatalf("failed to read echoed message: %v", err)
	}

	if !bytes.Equal(buf, msg) {
		t.Fatalf("expected %q, got %q", string(msg), string(buf))
	}
}

func TestSocksServerWithAuthFailure(t *testing.T) {
	echoLn, echoAddr := startEchoServer(t)
	defer echoLn.Close()

	proxyLn, proxyAddr := startTestSocksServer(t, "alice", "secret123")
	defer proxyLn.Close()

	_, err := socks5Connect(proxyAddr, echoAddr, "alice", "wrongpass")
	if err == nil {
		t.Fatal("expected error with wrong password, but succeeded")
	}
}

func TestSocksServerConcurrent(t *testing.T) {
	echoLn, echoAddr := startEchoServer(t)
	defer echoLn.Close()

	proxyLn, proxyAddr := startTestSocksServer(t, "", "")
	defer proxyLn.Close()

	var wg sync.WaitGroup
	concurrentClients := 10

	for i := 0; i < concurrentClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client, err := socks5Connect(proxyAddr, echoAddr, "", "")
			if err != nil {
				t.Errorf("client %d connect failed: %v", id, err)
				return
			}
			defer client.Close()

			msg := []byte(fmt.Sprintf("concurrent message from %d", id))
			if _, err := client.Write(msg); err != nil {
				t.Errorf("client %d write failed: %v", id, err)
				return
			}

			buf := make([]byte, len(msg))
			if _, err := io.ReadFull(client, buf); err != nil {
				t.Errorf("client %d read failed: %v", id, err)
				return
			}

			if !bytes.Equal(buf, msg) {
				t.Errorf("client %d expected %q, got %q", id, string(msg), string(buf))
			}
		}(i)
	}

	wg.Wait()
}
