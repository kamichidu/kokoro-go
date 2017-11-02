package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

var appVersion string

var commands = []cli.Command{}

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

func run(in io.Reader, out io.Writer, errOut io.Writer, args []string) int {
	// init global logger
	logger := log.StandardLogger()
	logger.Formatter = &log.TextFormatter{}
	logger.Out = errOut

	app := cli.NewApp()
	app.Writer = out
	app.ErrWriter = errOut
	app.Name = "kokoro-go"
	app.Version = appVersion
	app.Commands = commands
	app.Before = initLogger
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "loglevel",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  "logfile",
			Value: "-",
		},
		&cli.StringFlag{
			Name:  "c",
			Value: defaultConfigFilename,
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
