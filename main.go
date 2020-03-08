package main

import (
	"flag"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/network"
	"io"
	"net"
)

func main() {

	defer common.CrashLog()

	listen := flag.String("l", "", "listen addr")
	nolog := flag.Int("nolog", 0, "write log file")
	noprint := flag.Int("noprint", 0, "print stdout")
	loglevel := flag.String("loglevel", "info", "log level")

	flag.Parse()

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
	loggo.Info("start...")

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

	for {
		conn, err := tcplistenConn.AcceptTCP()
		if err != nil {
			loggo.Info("Error accept tcp %s", err)
			continue
		}

		go process(conn)
	}
}

func process(conn *net.TCPConn) {

	defer common.CrashLog()

	var err error = nil
	if err = network.Sock5HandshakeBy(conn); err != nil {
		loggo.Error("socks handshake: %s", err)
		conn.Close()
		return
	}
	_, targetAddr, err := network.Sock5GetRequest(conn)
	if err != nil {
		loggo.Error("error getting request: %s", err)
		conn.Close()
		return
	}
	// Sending connection established message immediately to client.
	// This some round trip time for creating socks connection with the client.
	// But if connection failed, the client will get connection reset error.
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		loggo.Error("send connection confirmation: %s", err)
		conn.Close()
		return
	}

	loggo.Info("accept new sock5 conn: %s", targetAddr)

	tcpsrcaddr := conn.RemoteAddr().(*net.TCPAddr)

	loggo.Info("client accept new direct local tcp %s %s", tcpsrcaddr.String(), targetAddr)

	tcpaddrTarget, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		loggo.Info("direct local tcp ResolveTCPAddr fail: %s %s", targetAddr, err.Error())
		return
	}

	targetconn, err := net.DialTCP("tcp", nil, tcpaddrTarget)
	if err != nil {
		loggo.Info("direct local tcp DialTCP fail: %s %s", targetAddr, err.Error())
		return
	}

	go transfer(conn, targetconn, conn.RemoteAddr().String(), targetconn.RemoteAddr().String())
	go transfer(targetconn, conn, targetconn.RemoteAddr().String(), conn.RemoteAddr().String())

	loggo.Info("client accept new direct local tcp ok %s %s", tcpsrcaddr.String(), targetAddr)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, dst string, src string) {

	defer common.CrashLog()

	defer destination.Close()
	defer source.Close()
	loggo.Info("client begin transfer from %s -> %s", src, dst)
	io.Copy(destination, source)
	loggo.Info("client end transfer from %s -> %s", src, dst)
}
