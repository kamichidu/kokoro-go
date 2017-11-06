package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

type cmdWS struct{}

func (self *cmdWS) Action(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.NewExitError("no channelIds", 128)
	}
	channelIds := make([]string, c.NArg())
	for i := 0; i < c.NArg(); i++ {
		channelIds[i] = c.Args().Get(i)
	}
	log.Debugf("Subscribing %v", channelIds)

	config, err := newAppConfigFrom(c.GlobalString("c"))
	if err != nil {
		return err
	}

	u := self.createUrl(c)
	log.Debugf("Connecting %s", u.String())

	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"X-Access-Token": []string{config.Token},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer ws.Close()

	if c.Bool("debug") {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		fmt.Fprintln(c.App.ErrWriter, string(b))
	}

	return self.loop(c.App.Writer, ws, channelIds)
}

type cableEvent struct {
	Type string `json:"type,omitempty"`

	Message *json.RawMessage `json:"message,omitempty"`

	Identifier string `json:"identifier,omitempty"`
}

type cableCommand struct {
	Command string `json:"command"`

	Data string `json:"data,omitempty"`

	Identifier string `json:"identifier"`
}

func (self *cmdWS) debugDumpJson(tag string, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	log.Debugf("%s %s", tag, string(b))
}

func (self *cmdWS) loop(w io.Writer, ws *websocket.Conn, channelIds []string) error {
	for {
		mt, r, err := ws.NextReader()
		if err != nil {
			log.Infof("connection was closed: %s", err)
			return nil
		} else if mt != websocket.TextMessage {
			log.Errorf("received frame is not a text frame (%v)", mt)
			continue
		}

		evt := &cableEvent{}
		if err := json.NewDecoder(r).Decode(evt); err != nil {
			log.Errorf("received a invalid message: %s", err)
			continue
		}
		self.debugDumpJson("IN", evt)

		switch evt.Type {
		case "welcome":
			msg := &cableCommand{
				Command:    "subscribe",
				Identifier: `{"channel":"ChatChannel"}`,
			}
			self.debugDumpJson("OUT", msg)
			err := ws.WriteJSON(msg)
			if err != nil {
				return err
			}
		case "confirm_subscription":
			// subscribe channels
			msg := &cableCommand{
				Command:    "message",
				Identifier: `{"channel":"ChatChannel"}`,
			}
			dataBytes, err := json.Marshal(map[string]interface{}{
				"channels": channelIds,
				"action":   "subscribe",
			})
			if err != nil {
				log.Panic(err)
			}
			msg.Data = string(dataBytes)
			self.debugDumpJson("OUT", msg)
			err = ws.WriteJSON(msg)
			if err != nil {
				return err
			}
		case "":
			if evt.Message == nil {
				continue
			}
			if err := json.NewEncoder(w).Encode(evt.Message); err != nil {
				return err
			}
		}
	}
}

func (self *cmdWS) createUrl(c *cli.Context) *url.URL {
	u := &url.URL{
		Host: c.GlobalString("host"),
		Path: c.String("path"),
	}
	u.Scheme = "wss"
	if c.GlobalBool("insecure") {
		u.Scheme = "ws"
	}
	return u
}

func init() {
	cmd := &cmdWS{}
	commands = append(commands, cli.Command{
		Name: "websocket",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "/cable",
				Usage: "ActionCable endpoint path",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Dump HTTP Response",
			},
		},
		Action: cmd.Action,
	})
}
