package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/kamichidu/kokoro-go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

//go:generate go-bindata ./templates/...

const appName = "kokoro-go"

var appVersion string

var commands = []cli.Command{}

const (
	metaReader       = "reader"
	metaKokoroClient = "kokoroClient"
)

func initLogger(c *cli.Context) error {
	l := log.StandardLogger()
	if lvl, err := log.ParseLevel(c.String("loglevel")); err != nil {
		return cli.NewExitError(err, 128)
	} else {
		l.Level = lvl
	}
	if filename := c.String("logfile"); filename != "-" {
		filename = filepath.Clean(filepath.FromSlash(filename))

		if info, err := os.Stat(filepath.Dir(filename)); err != nil {
			// Output directory is not exists, create it
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				return err
			}
		} else if !info.IsDir() {
			return cli.NewExitError(fmt.Sprintf("Invalid logfile: %s", filename), 128)
		}

		wc, err := os.Create(filename)
		if err != nil {
			return err
		}
		l.Out = wc
	}
	return nil
}

func initKokoroClient(c *cli.Context) error {
	config, err := newAppConfigFrom(c.GlobalString("c"))
	if err != nil {
		return err
	}

	u := &url.URL{}
	u.Scheme = "https"
	if c.GlobalBool("insecure") {
		u.Scheme = "http"
	}
	u.Host = c.GlobalString("host")

	c.App.Metadata[metaKokoroClient] = kokoro.NewClientWithToken(u.String(), config.Token)

	return nil
}

func run(in io.Reader, out io.Writer, errOut io.Writer, args []string) int {
	// init global logger
	logger := log.StandardLogger()
	logger.Formatter = &log.TextFormatter{}
	logger.Out = errOut

	app := cli.NewApp()
	app.Metadata = map[string]interface{}{
		metaReader: in,
	}
	app.Writer = out
	app.ErrWriter = errOut
	app.Name = appName
	app.Version = appVersion
	app.Usage = "kokoro-io client tool"
	app.Commands = commands
	app.Before = initLogger
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "loglevel",
			Value: "info",
			Usage: "Log threshold level (Possible values: debug, info, warn, error, fatal, panic)",
		},
		&cli.StringFlag{
			Name:  "logfile",
			Value: "-",
			Usage: "Write log messages to `PATH`",
		},
		&cli.StringFlag{
			Name:  "c",
			Value: defaultConfigFilename,
			Usage: "Config file `PATH`",
		},
		&cli.BoolFlag{
			Name:  "insecure",
			Usage: "Use http/ws instead of https/wss",
		},
		&cli.StringFlag{
			Name:  "host",
			Value: "kokoro.io",
			Usage: "kokoro.io hostname",
		},
	}
	if err := app.Run(args); err != nil {
		// err is already reported by app.Run() inside
		if eerr, ok := err.(*cli.ExitError); ok {
			return eerr.ExitCode()
		} else {
			fmt.Fprintln(errOut, err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
