package main

import (
	"fmt"
	"text/tabwriter"
	"text/template"

	"github.com/kamichidu/kokoro-go"
	"gopkg.in/urfave/cli.v1"
)

var (
	tplChannels = template.Must(template.New("channels").Parse(string(MustAsset("templates/channels.tpl"))))
	tplChannel  = template.Must(template.New("channel").Parse(string(MustAsset("templates/channel.tpl"))))
)

type cmdChannel struct{}

func (self *cmdChannel) ListAction(c *cli.Context) error {
	kcli := c.App.Metadata[metaKokoroClient].(*kokoro.Client)

	opts := kokoro.QueryOptions{}
	if c.Bool("archived") {
		opts = append(opts, kokoro.Archived)
	}
	if c.Bool("no-archived") {
		opts = append(opts, kokoro.NotArchived)
	}
	channels, err := kcli.ListChannels(opts...)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(c.App.Writer, 1, 0, 1, ' ', 0)
	if err := tplChannels.Execute(tw, channels); err != nil {
		return err
	}
	return tw.Flush()
}

func (self *cmdChannel) ShowAction(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.NewExitError("no {channelId} given", 128)
	}

	kcli := c.App.Metadata[metaKokoroClient].(*kokoro.Client)

	channel, err := kcli.GetChannel(c.Args().First())
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(c.App.Writer, 1, 0, 1, ' ', 0)
	if err := tplChannel.Execute(tw, channel); err != nil {
		return err
	}
	return tw.Flush()
}

func init() {
	cmd := &cmdChannel{}
	commands = append(commands, cli.Command{
		Name:   "channel",
		Usage:  "Manage channel",
		Before: initKokoroClient,
		Subcommands: cli.Commands{
			cli.Command{
				Name:  "list",
				Usage: "List channels",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "archived",
						Usage: "List archived channels",
					},
					&cli.BoolFlag{
						Name:  "no-archived",
						Usage: "List not archived channels",
					},
				},
				Action: cmd.ListAction,
			},
			cli.Command{
				Name:      "show",
				Usage:     "Show channel",
				UsageText: fmt.Sprintf("%s channel show [command options] {channelId}", appName),
				Action:    cmd.ShowAction,
			},
		},
	})
}
