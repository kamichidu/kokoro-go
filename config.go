package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
)

var defaultConfigFilename string

func init() {
	const filename = "kokoro-go/config.json"

	if runtime.GOOS == "windows" {
		defaultConfigFilename = filepath.Join(os.Getenv("AppData"), filename)
	} else {
		defaultConfigFilename = filepath.Join(os.Getenv("HOME"), filename)
	}
}

type appConfig struct {
	Token string `json:"token"`
}

func (self *appConfig) WriteFile(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	wc, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer wc.Close()

	return json.NewEncoder(wc).Encode(self)
}

func newAppConfigFrom(filename string) (*appConfig, error) {
	log.Debugf("loading config from %s", filename)
	v := &appConfig{}
	rc, err := os.Open(filename)
	if os.IsNotExist(err) {
		log.Debugf("config file is not exists (%s)", filename)
		return v, nil
	} else if err != nil {
		log.Debugf("can't open file (%s): %s", filename, err)
		return nil, err
	}
	defer rc.Close()

	return v, json.NewDecoder(rc).Decode(v)
}

type cmdConfig struct{}

func (self *cmdConfig) Action(c *cli.Context) error {
	configFilename := c.GlobalString("c")
	config, err := newAppConfigFrom(configFilename)
	if err != nil {
		return err
	}

	switch c.NArg() {
	case 0:
		return self.list(c.App.Writer, config)
	case 1:
		return self.get(c.App.Writer, config, c.Args().First())
	case 2:
		if err := self.set(config, c.Args().Get(0), c.Args().Get(1)); err != nil {
			return err
		}

		return config.WriteFile(configFilename)
	default:
		return cli.NewExitError("See help", 128)
	}
}

func (self *cmdConfig) list(w io.Writer, config *appConfig) error {
	typ := reflect.TypeOf(config).Elem()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name := f.Tag.Get("json")
		if err := self.get(w, config, name); err != nil {
			return err
		}
	}
	return nil
}

func (self *cmdConfig) get(w io.Writer, config *appConfig, name string) error {
	rv := reflect.ValueOf(config).Elem()
	typ := rv.Type()
	for i := 0; i < typ.NumField(); i++ {
		ftyp := typ.Field(i)
		if ftyp.Tag.Get("json") == name {
			fv := rv.Field(i)
			fmt.Fprintf(w, "%s=%v\n", name, fv.Interface())
			return nil
		}
	}
	return cli.NewExitError(fmt.Sprintf("No such config: %s", name), 128)
}

func (self *cmdConfig) set(config *appConfig, name string, value string) error {
	rv := reflect.ValueOf(config).Elem()
	typ := rv.Type()
	for i := 0; i < typ.NumField(); i++ {
		ftyp := typ.Field(i)
		if ftyp.Tag.Get("json") == name {
			fv := rv.Field(i)
			switch fv.Kind() {
			case reflect.String:
				fv.SetString(value)
			default:
				log.Panicf("no implementation for %s", fv.Kind())
			}
			return nil
		}
	}
	return cli.NewExitError(fmt.Sprintf("No such config: %s", name), 128)
}

func init() {
	cmd := &cmdConfig{}
	commands = append(commands, cli.Command{
		Name:   "config",
		Action: cmd.Action,
	})
}
