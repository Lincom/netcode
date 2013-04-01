package netcode

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var stop chan bool = make(chan bool)
var Codes map[int64]Responder = make(map[int64]Responder)

type Responder interface {
	Respond(connection net.Conn, content string) // can just fmt.FPrintf(connection, content, ...)
}

func Stop() {
	stop <- true
}

func Listen(nettype string, loc string) {
	ln, err := net.Listen(nettype, loc)
	if err != nil {
		panic(err)
	}

	nextConn := make(chan net.Conn)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("Error accepting new connection: %s\n", err)
				continue
			}
			nextConn <- conn
		}
	}()

	terminate := false

	for !terminate {
		select {
		case conn := <-nextConn:
			handleConnectionRead(conn)
		case <-stop:
			terminate = true
		}
	}

	return
}

func handleConnectionRead(conn net.Conn) {
	reader := bufio.NewReader(conn)
	stop := false
	for !stop {
		req, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// closed connection; do not read again after processing
				stop = true
			} else {
				fmt.Printf("Error reading from %s: %s\n", conn.RemoteAddr().String(), err)
				continue
			}
		}

		if req == "" { // nothing to do... (useful for EOF)
			continue
		}

		split := strings.SplitN(req, " ", 2)
		detail := ""
		if len(split) < 1 {
			fmt.Printf("Invalid request from %s: %s\n", conn.RemoteAddr().String(), req)
			continue
		} else if len(split) > 1 {
			detail = split[1]
		}
		code, err := strconv.ParseInt(split[0], 10, 0)
		if err != nil {
			fmt.Printf("Invalid code from %s: requested %s: %s\n", conn.RemoteAddr().String(), req, err)
			continue
		}

		res, ok := Codes[code]
		if !ok {
			fmt.Printf("I don't know how to handle code %d from %s.\n", code, conn.RemoteAddr().String())
			continue
		}
		res.Respond(conn, detail)
	}
}
