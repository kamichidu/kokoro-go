package main

import (
	"bufio"
	"fmt"
	"io"
	"syscall"

	"github.com/kamichidu/kokoro-go"
	"github.com/kamichidu/kokoro-go/types"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultDeviceName = "kokoro-go"
	defaultDeviceId   = "kokoro-go"
)

type cmdLogin struct{}

func (self *cmdLogin) Action(c *cli.Context) (err error) {
	r := c.App.Metadata[metaReader].(io.Reader)

	email := c.String("email")
	if email == "" {
		email, err = self.prompt(c.App.Writer, r, "Email: ", false)
	}
	password := c.String("password")
	if password == "" {
		password, err = self.prompt(c.App.Writer, r, "Password: ", true)
	}

	var baseUrl string
	if c.GlobalBool("insecure") {
		baseUrl += "http://"
	} else {
		baseUrl += "https://"
	}
	baseUrl += c.GlobalString("host")

	device, err := kokoro.NewClient(baseUrl).RegisterDevice(email, password, &types.Device{
		Name:             c.String("device-name"),
		Kind:             types.DeviceUnknown,
		DeviceIdentifier: c.String("device-id"),
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return c.App.Run([]string{c.App.Name, "config", "token", device.AccessToken.Token})
}

// XXX: expose for test
var readPassword = func() ([]byte, error) {
	return terminal.ReadPassword(int(syscall.Stdin))
}

func (self *cmdLogin) prompt(w io.Writer, r io.Reader, msg string, secure bool) (string, error) {
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return "", err
	}

	if secure {
		// XXX: readPassword() will exit without eol
		defer io.WriteString(w, "\n")
		bytes, err := readPassword()
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		br := bufio.NewReader(r)
		bytes, _, err := br.ReadLine()
		return string(bytes), err
	}
}

func init() {
	cmd := &cmdLogin{}
	commands = append(commands, cli.Command{
		Name:  "login",
		Usage: "Logging in to kokoro.io",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "device-name",
				Value: defaultDeviceName,
				Usage: "Device name for register a device",
			},
			&cli.StringFlag{
				Name:  "device-id",
				Value: defaultDeviceId,
				Usage: "Device identifier for register a device",
			},
			&cli.StringFlag{
				Name:  "email",
				Usage: "Your login email",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Your login password",
			},
		},
		Action: cmd.Action,
	})
}
