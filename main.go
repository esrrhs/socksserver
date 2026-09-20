package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/esrrhs/gohome/common"
	"github.com/esrrhs/gohome/loggo"
	"github.com/esrrhs/gohome/network"
)

var (
	version   = "0.3"
	buildDate = "unknown"
)

func main() {
	defer common.CrashLog()

	listen := flag.String("l", "", "listen addr")
	nolog := flag.Int("nolog", 0, "write log file")
	noprint := flag.Int("noprint", 0, "print stdout")
	loglevel := flag.String("loglevel", "info", "log level")
	user := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	showVersion := flag.Bool("v", false, "show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("socksserver %s (built %s)\n", version, buildDate)
		return
	}

	if *listen == "" {
		flag.Usage()
		return
	}

	level := loggo.LEVEL_INFO
	if loggo.NameToLevel(*loglevel) >= 0 {
		level = loggo.NameToLevel(*loglevel)
	}
	loggo.Ini(loggo.Config{
		Level:     level,
		Prefix:    "socksserver",
		MaxDay:    3,
		NoLogFile: *nolog > 0,
		NoPrint:   *noprint > 0,
	})
	loggo.Info("socksserver %s starting...", version)

	tcpaddr, err := net.ResolveTCPAddr("tcp", *listen)
	if err != nil {
		loggo.Error("listen fail %s", err)
		return
	}

	tcplistenConn, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		loggo.Error("Error listening for tcp packets: %s", err)
		return
	}
	loggo.Info("listen ok %s", tcpaddr.String())

	// Graceful shutdown on SIGINT / SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		loggo.Info("received signal %s, shutting down listener...", sig)
		_ = tcplistenConn.Close()
	}()

	for {
		conn, err := tcplistenConn.AcceptTCP()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				loggo.Info("listener closed, exiting accept loop")
				break
			}
			loggo.Info("Error accept tcp %s", err)
			continue
		}

		go process(conn, *user, *password)
	}
	loggo.Info("socksserver stopped")
}

func process(conn *net.TCPConn, user string, password string) {
	defer common.CrashLog()

	var err error
	if err = network.Sock5HandshakeBy(conn, user, password); err != nil {
		loggo.Error("socks handshake: %s", err)
		_ = conn.Close()
		return
	}
	_, targetAddr, err := network.Sock5GetRequest(conn)
	if err != nil {
		loggo.Error("error getting request: %s", err)
		_ = conn.Close()
		return
	}
	// Sending connection established message immediately to client.
	// This saves some round trip time for creating socks connection with the client.
	// But if connection failed, the client will get connection reset error.
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		loggo.Error("send connection confirmation: %s", err)
		_ = conn.Close()
		return
	}

	loggo.Info("accept new sock5 conn: %s", targetAddr)

	tcpsrcaddr := conn.RemoteAddr().String()
	loggo.Info("client accept new direct local tcp %s -> %s", tcpsrcaddr, targetAddr)

	targetconn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		loggo.Info("direct local tcp dial fail: %s %s", targetAddr, err.Error())
		_ = conn.Close()
		return
	}

	go transfer(conn, targetconn, tcpsrcaddr, targetconn.RemoteAddr().String())
	go transfer(targetconn, conn, targetconn.RemoteAddr().String(), tcpsrcaddr)

	loggo.Info("client accept new direct local tcp ok %s -> %s", tcpsrcaddr, targetAddr)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, dst string, src string) {
	defer common.CrashLog()

	defer destination.Close()
	defer source.Close()
	loggo.Info("client begin transfer from %s -> %s", src, dst)
	_, _ = io.Copy(destination, source)
	loggo.Info("client end transfer from %s -> %s", src, dst)
}
