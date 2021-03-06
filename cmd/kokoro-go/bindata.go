// Code generated by go-bindata.
// sources:
// templates/channel.tpl
// templates/channels.tpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesChannelTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x90\x31\x4b\xc6\x30\x10\x86\xf7\xfb\x15\xa1\x38\xe8\x60\x7e\x40\xc1\xa1\xd8\xa5\x48\x45\xd0\xd1\x25\xb6\xa7\x39\x68\xd3\x90\x44\xa1\x1c\xf7\xdf\x25\x09\x95\x0e\xf2\x7d\x53\xf2\xbc\xf7\x2c\xef\x3b\xf4\xad\x62\xd6\xc3\x2c\x02\xcf\x66\xc5\x42\x8f\xd6\x38\x87\x4b\x66\x11\x78\x22\x37\x97\x38\x7f\x44\xa0\x0b\x93\xa5\x1f\xac\xd9\x01\x22\xd0\x63\x9c\x02\xf9\x44\x9b\x6b\x81\x59\x9f\x58\x04\x60\xc4\xf5\x03\x43\xcc\xa7\x7b\x15\x8c\xfb\x42\x75\xb3\xd6\xcc\x92\x57\xed\x83\xd2\xe3\x1f\x46\x11\x60\x3e\xdd\x75\xf7\x9d\xec\x16\x28\xed\x22\xcc\xcd\x7b\x6a\xf2\x7b\x16\x5e\xc2\xf6\x49\x0b\x96\x2a\x97\x8d\xb7\xdd\xe3\x35\xe7\x75\x0a\x88\xae\x2e\x70\x98\xb7\xff\xab\x3d\x45\xbf\x98\xbd\xba\x77\xa5\x1e\x96\xa1\x7e\x03\x00\x00\xff\xff\x61\x33\x77\x4f\x5c\x01\x00\x00")

func templatesChannelTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesChannelTpl,
		"templates/channel.tpl",
	)
}

func templatesChannelTpl() (*asset, error) {
	bytes, err := templatesChannelTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/channel.tpl", size: 348, mode: os.FileMode(438), modTime: time.Unix(1509960857, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesChannelsTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x74\xa9\xae\x56\x8a\x29\x51\xaa\xad\xf5\x4b\xcc\x4d\x85\xb1\xbd\x33\xf3\x52\x60\x6c\xc7\xa2\xe4\x8c\xcc\xb2\x54\x38\xdf\x25\xb5\x38\xb9\x28\xb3\xa0\x24\x33\x3f\x8f\xab\xba\x5a\x57\xa1\x28\x31\x2f\x3d\x55\x41\x25\x39\x23\x31\x2f\x2f\x35\x47\xc1\xca\x56\x41\xaf\xb6\x96\xab\xba\x1a\x26\xa2\xe7\x99\x52\x5b\x0b\xd3\x8c\x24\xec\x0c\xa1\x41\xd6\x62\x95\x07\xb9\x01\xab\x04\xcc\x41\x58\x25\x91\x5c\x07\x76\x85\xae\x42\x2a\xc8\x18\x2e\x40\x00\x00\x00\xff\xff\x22\xae\xf3\x23\xe8\x00\x00\x00")

func templatesChannelsTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesChannelsTpl,
		"templates/channels.tpl",
	)
}

func templatesChannelsTpl() (*asset, error) {
	bytes, err := templatesChannelsTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/channels.tpl", size: 232, mode: os.FileMode(438), modTime: time.Unix(1509959142, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/channel.tpl": templatesChannelTpl,
	"templates/channels.tpl": templatesChannelsTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"channel.tpl": &bintree{templatesChannelTpl, map[string]*bintree{}},
		"channels.tpl": &bintree{templatesChannelsTpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

