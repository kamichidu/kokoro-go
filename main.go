package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
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

func run(in io.Reader, out io.Writer, errOut io.Writer, args []string) int {
	// parse flags
	var (
		showVersion bool
		config      = new(appConfig)
	)
	flags := flag.NewFlagSet("kokoro-go", flag.ExitOnError)
	flags.BoolVar(&showVersion, "v", false, "Show version")
	flags.StringVar(&config.Host, "H", "kokoro.io", "`HOST` for REST API/WebSocket API")
	flags.BoolVar(&config.Insecure, "k", false, "Use insecure connection")
	flags.StringVar(&config.LogLevel, "l", "info", "Log `LEVEL`")
	flags.StringVar(&config.LogFile, "f", "-", "Log `FILE`")
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(errOut, "Can't parse arguments: %s\n", err)
		return 128
	}

	if showVersion {
		fmt.Fprintln(out, appVersion)
		return 0
	}

	config.AccessToken = flags.Arg(0)
	if config.AccessToken == "" {
		flags.Usage()
		return 128
	}

	// init logger
	logger := log.StandardLogger()
	logger.Formatter = &log.TextFormatter{}
	if config.LogFile == "-" {
		logger.Out = errOut
	} else {
		if err := os.MkdirAll(filepath.Dir(config.LogFile), 0755); err != nil {
			fmt.Fprintf(errOut, "Can't create directory: %s\n", err)
			return 128
		}
		if fw, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
			fmt.Fprintf(errOut, "Can't create logfile: %s\n", err)
			return 128
		} else {
			defer fw.Close()
			logger.Out = fw
		}
	}
	if lvl, err := log.ParseLevel(config.LogLevel); err != nil {
		fmt.Fprintf(errOut, "Can't recognize loglevel: %s\n", err)
		return 128
	} else {
		logger.Level = lvl
	}

	// build app
	app := &kokoro{
		mutex:       new(sync.Mutex),
		w:           out,
		Insecure:    config.Insecure,
		Host:        config.Host,
		AccessToken: config.AccessToken,
		Logger:      logger,
	}
	if err := app.Start(context.Background(), in); err != nil {
		logger.Errorf("Error: %s", err)
		return 1
	}
	return 0

	// ws, _, err := websocket.DefaultDialer.Dial(proxyUrl, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// defer ws.Close()
	//
	// go io.Copy(&wsWriter{ws}, r)
	// _, err = io.Copy(w, &wsReader{ws})
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
