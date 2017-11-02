package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type cmdReq struct{}

func (self *cmdReq) Action(c *cli.Context) error {
	config, err := newAppConfigFrom(c.GlobalString("c"))
	if err != nil {
		return err
	}

	methods := map[string]string{
		"get":    http.MethodGet,
		"post":   http.MethodPost,
		"put":    http.MethodPut,
		"delete": http.MethodDelete,
	}
	method, ok := methods[c.Command.Name]
	if !ok {
		log.Panicf("%s is not supported http method", c.Command.Name)
	}

	u, err := self.createUrl(c)
	if err != nil {
		return err
	}

	var body io.Reader
	if c.Bool("read") {
		body = c.App.Metadata[metaReader].(io.Reader)
	}

	log.Debugf("%s %s", method, u.String())
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("X-Access-Token", config.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.GlobalBool("debug") {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		fmt.Fprintln(c.App.Writer, string(b))
		return nil
	}

	// print response payload
	_, err = io.Copy(c.App.Writer, resp.Body)
	if err != nil {
		return err
	}

	// exit code depends on http response status code
	// - 2xx   => 0
	// - other => same to status code
	// 2xx => 0
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return cli.NewExitError(resp.Status, resp.StatusCode)
	}
}

func (self *cmdReq) createUrl(c *cli.Context) (*url.URL, error) {
	scheme := "https"
	if c.GlobalBool("insecure") {
		scheme = "http"
	}
	host := c.GlobalString("host")
	uriPath := c.Args().First()
	if uriPath == "" {
		return nil, cli.NewExitError("{uriPath} must be specified", 128)
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   uriPath,
	}
	if queries := c.StringSlice("query"); queries != nil {
		q := u.Query()
		for _, query := range queries {
			kv := strings.SplitN(query, "=", 2)
			if len(kv) == 1 {
				q.Add(kv[0], "")
			} else {
				q.Add(kv[0], kv[1])
			}
		}
		u.RawQuery = q.Encode()
	}
	return u, nil
}

func init() {
	cmd := &cmdReq{}
	newSub := func(name string) cli.Command {
		return cli.Command{
			Name: name,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name: "query",
				},
			},
			Action: cmd.Action,
		}
	}
	withBody := func(c cli.Command) cli.Command {
		c.Flags = append(c.Flags, &cli.BoolFlag{
			Name: "read",
		})
		return c
	}
	commands = append(commands, cli.Command{
		Name: "request",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "insecure",
			},
			&cli.StringFlag{
				Name:  "host",
				Value: "kokoro.io",
			},
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		Subcommands: cli.Commands{
			newSub("get"),
			withBody(newSub("post")),
			withBody(newSub("put")),
			newSub("delete"),
		},
	})
}
