package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/gorilla/websocket"
)

var appVersion string

type wsReader struct {
	*websocket.Conn
}

func (self *wsReader) Read(p []byte) (int, error) {
	mt, b, err := self.Conn.ReadMessage()
	if err != nil {
		return len(b), err
	}
	if mt != websocket.TextMessage {
		return 0, errors.New("Only supported TextMessage")
	}
	return copy(p, b), nil
}

var _ io.Reader = (*wsReader)(nil)

type wsWriter struct {
	*websocket.Conn
}

func (self *wsWriter) Write(data []byte) (int, error) {
	return len(data), self.Conn.WriteMessage(websocket.TextMessage, data)
}

var _ io.Writer = (*wsWriter)(nil)

func run() int {
	var err error
	var (
		showVersion bool
		proxyUrl    string
	)
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&proxyUrl, "u", "ws://kokoro.io", "Proxy target url")
	flag.Parse()

	if showVersion {
		println(appVersion)
		return 0
	}

	r := io.Reader(os.Stdin)
	w := io.Writer(os.Stdout)

	ws, _, err := websocket.DefaultDialer.Dial(proxyUrl, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	go io.Copy(&wsWriter{ws}, r)
	_, err = io.Copy(w, &wsReader{ws})
	if err != nil {
		panic(err)
	}
	return 0
}

func main() {
	os.Exit(run())
}
