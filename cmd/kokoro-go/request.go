package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jmespath/go-jmespath"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
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
	req.Header.Set("Content-Type", "application/json")

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

	// exit code depends on http response status code
	// - 2xx   => 0
	// - other => same to status code
	// 2xx => 0
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return self.printResponse(c.App.Writer, resp, c.String("query"))
	} else {
		io.Copy(c.App.Writer, resp.Body)
		return cli.NewExitError(resp.Status, resp.StatusCode)
	}
}

func (self *cmdReq) printResponse(w io.Writer, resp *http.Response, query string) error {
	if query == "" {
		log.Debug("no jmespath expr, just copy it")
		_, err := io.Copy(w, resp.Body)
		return err
	}

	log.Debugf("jmespath expr (%s), apply it", query)
	buf := json.RawMessage{}
	if err := json.NewDecoder(resp.Body).Decode(&buf); err != nil {
		return err
	}
	var data interface{}
	if len(buf) == 0 {
		data = nil
	} else {
		switch rune(buf[0]) {
		case '[':
			data = []interface{}{}
		default:
			data = map[string]interface{}{}
		}
		if err := json.Unmarshal(buf, &data); err != nil {
			return err
		}
	}

	v, err := jmespath.Search(query, data)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(v)
}

func (self *cmdReq) createUrl(c *cli.Context) (*url.URL, error) {
	uriPath := c.Args().First()
	if uriPath == "" {
		return nil, cli.NewExitError("{uriPath} must be specified", 128)
	}
	u, err := url.ParseRequestURI(uriPath)
	if err != nil {
		return nil, cli.NewExitError("not a valid {uriPath}", 128)
	}
	u.Scheme = "https"
	if c.GlobalBool("insecure") {
		u.Scheme = "http"
	}
	u.Host = c.GlobalString("host")
	return u, nil
}

func init() {
	cmd := &cmdReq{}
	newSub := func(name string) cli.Command {
		return cli.Command{
			Name: name,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "query",
					Usage: "JMESPath query string. See http://jmespath.org/ for more information and examples.",
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
