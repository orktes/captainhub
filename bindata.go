package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
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

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _pr_review_js = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5a\xdd\x73\xdb\xc6\x11\x7f\xd7\x5f\x71\x8a\x67\x02\xd0\xa1\x41\xdb\x8f\xe2\xa8\x4d\x27\xa9\x67\xfa\x61\x4b\x13\x25\xd3\x07\x86\xa5\x20\xe2\x28\x22\x02\x71\x2c\x0e\x10\xcb\x3a\xfc\xdf\xbb\x7b\xdf\x77\x38\x52\x94\xa5\xb6\x49\x63\x3c\x58\x24\xee\xf6\xf3\x76\x7f\xbb\xb7\xf4\xa2\xab\xe7\x6d\xc9\x6a\x72\x4b\xdb\xef\xe8\x7d\x49\x37\xb4\xe1\xef\x58\xf3\xae\xac\x68\xba\x80\x7f\xf8\x90\xac\xbb\xaa\xfa\x8e\xfe\xa3\xa3\xbc\xfd\x81\xd3\x06\x5e\x34\x74\x41\x9b\x86\x16\x86\x62\x40\x3e\x9e\x10\x32\x67\x35\x67\x15\xcd\x2a\x76\x9b\x26\xef\xca\xba\x28\xeb\x5b\xd2\xe8\x3d\xe4\xa6\x6b\x49\x79\x5b\xb3\x06\x5f\x27\xe4\xab\x90\xf1\x60\x7c\x02\x4c\xee\xf3\xc6\xd0\xa0\x16\x9c\x9c\x93\x8f\xbb\xb1\x5a\x41\x95\x8c\x54\xb5\x02\x4b\xb3\x8c\xe6\xf3\x65\xda\x57\x6c\x48\x8c\x85\xda\x1c\xcd\x5c\xea\x6c\x68\xd5\xaa\xbf\x5d\xef\x21\xbe\xe0\x09\x7e\x9b\x82\xf8\xd8\xdb\x9f\x7f\x26\x93\xe9\x78\x3f\x59\xb6\xee\xf8\x32\x35\x5a\xc8\x9d\x3b\xf1\x77\x37\x70\xad\x39\xa8\x91\xeb\x26\x74\x44\xdf\x74\x21\x2d\xc3\x7f\xea\x7c\x65\xf5\x12\xd4\xe5\x82\x18\x0d\x78\x06\x3b\x6e\xdb\x25\x39\x3f\x3f\x27\xaf\xad\xc5\x28\x61\x9d\xb7\x2d\x6d\x6a\xe0\x3f\x03\x4e\x75\x91\xc2\x11\x2f\xca\xdb\x4c\xbd\xef\x69\xf7\x0d\x2c\x5b\x0e\x1e\x0f\xae\xdc\x85\x5b\x34\xfd\xd8\xdb\xd8\x41\x0c\xa0\xaa\xde\x4e\xa3\xa5\xd2\x5c\x6f\x5e\xe5\xed\x7c\x69\xf5\xd2\x42\x32\xbe\xae\xca\x36\x4d\x86\xc9\xc0\xd5\x4d\x2d\xbb\xaa\x11\x32\x1a\xb9\x01\xab\xb6\x0c\x89\xe7\xb6\xc1\xd8\x21\x68\x68\xdb\x81\x33\x84\x68\x0c\xcd\xcb\xbc\x5d\x3e\x4c\xa7\x0e\x55\xcb\x24\xfa\x80\xc8\x82\xc9\x78\x26\xf3\xbc\x6e\xc9\x0d\x25\xed\x92\x12\x0e\xc4\x24\xe7\xe2\x33\xe6\x07\x08\x15\x09\x42\x38\xad\x0b\xda\x18\x46\x78\x82\xd2\x07\x5f\x7e\x49\x52\xe3\xba\x0c\x9c\x41\xff\x79\xb1\x48\xc3\xdc\x12\xa7\xfb\xea\x0d\x46\x81\xdd\xac\x0e\xfe\x77\xe4\xcd\xc0\x77\x8d\xb2\x74\x96\x6d\xca\x76\xc9\xba\xd6\x0a\xe8\xc1\xc1\xc0\x72\x79\xed\x18\x7d\xd2\x33\x1e\x35\x3e\x8d\x1c\x84\x94\xa5\x49\x77\x7a\xbb\x1b\xdc\x56\x0f\x45\x6f\xc3\xa2\xaf\x8f\x11\x88\x61\x92\x57\x0d\xcd\x8b\xad\xf4\x39\xe2\x0e\x32\x2b\x6b\xe0\xc1\xa9\x88\x8d\xd4\xd1\x43\xb3\x34\xaf\x66\xd9\x1d\xdd\xf2\xd4\x83\xa3\x81\x5a\xf5\x0d\x0b\xe5\x38\x4e\xf1\x6d\xb5\x46\x85\x24\xd6\x03\x8e\x1f\x44\x9a\xe7\x75\xc1\x56\x26\x6c\xce\x2d\x97\xc9\x7b\x88\xc0\x6c\x51\x31\xd6\xa4\xe2\xa3\xdc\x9a\x0e\xc8\x4b\x12\xa6\xf7\x40\xe7\xbe\x67\xcd\xc4\x67\x3e\x75\xb8\xc7\xd7\x1d\x74\x3b\xb8\x51\xa2\x5c\x3f\x27\x76\x43\x08\xed\x92\x4b\xef\xa9\x30\xf3\x18\x8d\x4f\xc0\x74\xb7\x2c\xe1\xdb\xab\x65\xfe\x3e\x5f\x4b\x44\x94\xfe\xd4\xc5\x40\xae\x98\x1a\x71\x0c\x74\x5a\xb2\x10\x22\x25\xee\x64\x7c\x99\x6b\x34\x36\x3a\x5a\x22\x5f\x41\x27\xfa\x2e\xd6\xb4\xa6\x45\x4a\xef\x69\xdd\x7e\x9b\xb7\xb9\x55\x34\x08\x51\x90\x63\x36\x65\xb8\x36\x53\x49\x9e\x61\x9a\x21\x1a\x95\xb5\xa9\x85\x0b\x55\x03\xc1\x13\x97\x96\x8d\x70\x95\x15\x95\xd5\xdd\xea\x46\x15\x93\x58\x01\x7d\x4c\x75\xff\xb8\x1b\x8c\x63\x0e\x8e\x1d\x85\xd0\x12\xe3\x3f\x9a\x29\x91\x24\x10\xd0\x4d\x39\xcf\x6f\x29\xb0\x4c\xfe\xb0\xa1\x9c\x01\xe0\x6d\x58\x73\x77\x4a\x3e\xb0\x0d\xf9\xa9\x43\xac\x2b\x01\x0f\xf3\xf9\x1d\x81\x98\x22\x9b\x1c\xbe\x21\x56\x7e\xdf\xe4\xf7\x25\x20\x23\x83\x92\xc2\xb9\x58\x63\x00\x93\x8d\x78\x25\x25\x93\x2d\xeb\x1a\xe8\x43\x0a\x9a\xfd\x58\xff\x58\x27\x32\x50\xb5\xc0\xaf\x40\xe2\x8b\x17\x06\x80\xb9\xd8\xe0\xf6\x00\x9e\xfa\x47\xb4\x0e\x3e\xeb\xaf\xb1\xa7\xd1\x5b\xe0\x63\x62\x14\x38\xb6\xc9\xf0\xf9\xbd\x22\xd7\xc8\x51\x94\x08\xe0\x76\xed\xb0\xdb\x99\xfa\xe2\x12\x98\x0d\x06\x75\xfd\xd5\x17\xc6\x78\x70\xd1\x6a\x05\x0e\xe4\x51\x1f\xbd\x22\xf9\x7c\x4e\xd7\xed\x19\xb9\x5e\x37\x33\xe5\xd9\x8b\xbf\x5c\xef\xdb\x5c\x14\xc6\x6c\x8f\x04\x17\x26\xba\x72\x4c\xc9\x44\x81\xf7\x74\x1f\xa3\xf9\x32\xaf\xe1\x4b\x94\x97\x5a\x9b\xb0\xaa\x98\x39\x2c\x6b\xba\xb1\x5f\xf7\xf1\x5d\x43\xc7\x09\x41\xd5\xe6\x6d\xc7\x3d\xa6\xf2\xd5\xb5\x8d\x83\x39\x40\x72\x4b\xff\xc4\x79\x07\xdd\xc7\x6a\x05\xd9\xd5\xcb\xb1\xa1\x66\xae\x8e\x40\x92\x5c\x09\x4e\xba\x98\xec\xc9\xee\x25\xe0\x3d\x82\x8b\x2e\x30\x09\x00\x06\xf6\xc8\xc9\xf0\x30\x5d\xd7\x54\x86\xe4\xd2\xed\x09\xa4\x19\x67\xa2\x99\x8e\x67\xe0\x4f\xac\xac\xa1\x1b\x22\xd0\x0e\x19\xa1\xda\xfc\xe4\x44\xd5\x31\xfc\xc3\xf3\x7b\x8a\x92\x7b\xf6\x62\xe8\x9d\x61\x0e\x82\xa2\x8a\x70\xb6\x68\xd8\x0a\x98\xfe\xf9\xea\xe2\x43\xc6\x5b\xec\xe7\xcb\xc5\x36\x10\x2d\xd1\x1e\xdd\x7a\x98\xb5\x48\x89\x3e\x33\x0b\x3d\xc8\x69\x0f\xe2\xfe\xb0\x2e\xc0\xf7\xcf\x07\xb9\x4f\x40\x5c\xab\xef\x55\x8b\xd2\x2a\x96\x17\x0f\x19\x6d\x88\x21\xa8\x0d\x28\x1d\x43\x1e\x3b\x8e\x38\x68\xdb\x9b\x93\x2b\xc3\x5e\x9c\xda\x66\xdb\xab\x8a\xb0\x2a\x0e\x63\x9d\x43\x93\x94\x7a\x86\xc9\x43\x85\x86\x15\x1a\xcf\x14\x81\x6b\xe7\x33\x09\x84\x38\x6c\x02\x13\xe3\x8c\xb4\xe7\xb1\xa9\xaf\x00\x2a\x0e\xe3\x25\xd6\x9d\xfd\xb5\xfc\xf4\xdc\x56\x73\xff\x1a\x77\x4c\xed\x97\xfd\x8d\xa8\xfc\x6d\xd3\x51\x05\xab\xf6\x8a\xf6\xd4\x22\xeb\xba\xe3\xc9\x55\xf4\x0a\x4b\x28\x5b\x88\x3b\x83\x74\x21\x10\x53\xd2\xc9\xec\x38\xba\x12\xfe\x3f\x17\x42\xb8\x77\xbd\xa7\x0d\x2c\x83\xe3\xb9\x72\x12\x76\x18\xae\xd1\xee\x99\x1c\x65\x73\xd0\x00\xef\xef\xa1\xcd\xca\xef\x41\x56\x57\xe3\xb5\x23\x60\x1a\x6c\x1c\x90\x33\xa9\x64\x60\xc6\x23\xeb\xd3\xe7\x02\xf5\x3f\x2e\x50\xbd\x43\xb2\x05\xaa\x2c\xbc\x9a\x54\xe2\x99\x2a\xb1\x1a\xb2\xe5\x95\xff\x03\x4e\x04\xdc\xad\xf2\xb5\x5f\xb3\x6e\x58\xb1\xf5\x36\xcd\xa5\xe4\x0c\x17\x32\x50\x1a\x2e\x85\x7a\xef\x7c\x55\xb8\x63\x26\xb1\x43\x8d\x4d\x20\x71\xbc\xb9\x89\x0e\x25\x07\x71\xd5\x2b\x33\x68\x70\x8e\x6b\x10\x8e\x90\x40\x92\x8b\xe6\xce\xbd\x7a\x05\x57\x08\xcd\x89\x77\x37\xd2\xad\xe9\x9b\xd7\x03\xad\x09\x09\x07\x38\x4d\x1b\x1d\x51\xe0\x82\xb2\x2f\x6b\xe8\xba\xca\xe7\x34\x1d\xfd\xfd\xeb\x11\x44\x54\xe2\x0d\x61\xec\xad\xde\x72\x05\xfd\xa2\x4c\x4f\x4f\x85\x8f\x0c\x71\xef\xca\xef\x4d\x3a\xd1\x4a\x0c\x70\xf8\xab\xc2\x19\x94\xb7\x7b\xf9\xa6\xc4\x22\x07\xab\x93\xd7\xd3\xac\x65\x7f\x65\x10\x82\xdf\xe4\x50\x13\xbd\xa1\xcb\x1c\xde\x90\x84\xdd\x25\x67\x8e\x3e\x22\x0a\xda\x06\x4f\xd4\x6d\x0a\x20\x76\x0e\xb6\x01\xfa\xc1\xf3\x52\xe4\xbe\x9d\x6e\xcd\xf6\xa5\x85\x45\xcd\x29\xe1\x9a\xd3\x38\x20\x2b\x68\x45\x5b\x1a\x40\x99\x8d\xdd\x69\xb8\xdf\xa4\xdb\x01\x33\x1e\xce\xdd\x88\xee\xeb\xe6\x5b\xda\xe6\x65\xd5\x6f\xde\xd4\x7b\x90\xd8\xd3\xfe\xc1\xba\x1b\xc4\xb4\x7d\x3c\x60\x35\xb2\x2d\x88\x92\x84\x77\x70\x8f\xe2\x08\x1e\x76\x19\xb1\x32\x8a\x92\xa4\x60\x35\x45\x24\x74\x12\x2a\x54\x36\x5a\x01\xca\x62\x48\xbe\xb0\xf7\xba\x35\x9e\x46\x71\xfa\x45\x8f\x78\x47\x68\x05\x31\xf6\x80\x21\xbd\x55\x42\x22\xb6\x45\x76\xf5\x8a\x45\x9c\x85\x53\x2b\x3c\xea\x27\xd7\x0d\x5f\x17\xbf\x86\xb8\x4f\xdf\x2f\x7e\x2c\xf9\x5d\xa9\xbb\x4d\x7f\xba\x01\x7f\xdd\x8d\x83\xcc\x85\xcb\xae\x97\xba\x18\x57\x98\xf3\x6f\xa6\x61\xec\x60\xa8\xc2\xbd\x55\x67\xd8\x24\x48\x90\x27\xe6\xfc\xc1\xac\x8f\xe7\xfd\x23\x32\xbf\x47\xa9\xcd\x7c\xdb\x33\x53\x3e\x8e\xa1\x33\x04\x7d\x53\x0c\x62\xb7\x2b\x48\xcf\x43\x3d\xa1\xfb\x44\x67\xf1\x52\x91\x43\x23\x7c\xfd\xec\x1e\x2b\xc8\x63\x19\xe7\xd8\x7f\xbb\xeb\xfb\xcb\x47\x49\x19\x20\x53\x51\x24\x65\x5f\x18\x5f\x17\x53\xd7\xa1\x71\x66\x44\xd4\xf3\x80\x6a\xc8\xf5\x93\x51\xf5\x37\x0f\x2a\x7b\x91\x5a\xde\x8e\xe4\xc1\xe2\x49\xe1\x94\x8c\x16\x7d\xa8\xdf\x07\x43\x2e\x10\x1d\x00\xa5\xd1\xcb\x00\x9f\x1a\xba\x62\xf7\xf4\x11\x10\xf5\x8b\xc5\xa0\x90\x30\xda\x7f\xa8\xcc\xf9\xef\x24\xca\x13\xcf\xea\xe5\x28\x38\x2b\x39\xec\xdc\x73\x56\xf8\x83\x5f\x1c\x71\x7f\xd1\xa7\x16\xad\x1c\xd1\x33\x8b\x23\x71\x7f\xeb\xdb\x07\x80\xf3\xad\x05\xce\xb8\x9c\x18\x8c\x3f\x2e\x98\xfe\x53\xb8\x7b\x3c\x7a\xc8\x50\x29\xf0\x27\x10\xb3\xf4\x36\x66\xda\xee\x79\x63\x56\x83\x0a\x03\x80\x76\x02\x55\x0c\x20\xf3\x3b\xfa\x47\xbc\x88\xe2\x98\x51\xde\x68\xcf\x08\x2a\xef\xce\x0c\xce\xf6\xd7\x92\x9d\x95\xd2\xff\x75\xcd\x30\x77\x4c\x8c\xaa\x26\x67\xfb\xbf\xbe\xdb\x54\x30\xdb\xfb\x9b\xd4\x4a\xdf\x11\x50\x2b\x67\x9e\x67\x9f\x4f\x1e\xdc\xd9\xe7\xd8\x11\x5e\x4f\xe8\x83\xc3\xbc\xb8\x8c\x57\xc4\x9d\xea\x45\x05\xc4\x9a\xaa\xe8\xa0\xcf\x27\x09\x5e\xed\xc9\xa5\x60\x4e\x66\x9f\x4f\xbf\x12\x3e\x57\xbf\xf4\xf9\x36\xf9\xab\x6d\xfc\x9e\x76\x9b\xd4\xff\xf9\x43\xc6\x31\xf0\x1a\x8d\xc8\x65\xc3\xf0\xdc\xc9\x92\xb1\x3b\x52\x00\x66\x9c\x60\x80\x8a\x51\xdf\xf7\xdb\x35\x15\xd1\x98\xb8\xd8\x9a\x60\x8b\x60\x47\x81\xb9\xcc\x4a\xb1\x8d\x09\x14\x4d\x64\xec\x1e\xfa\xbf\x0b\x20\x5c\x1e\x70\x44\x96\x98\x55\xce\xd4\x84\xf1\x80\x30\x19\x09\x45\xb0\x43\x4e\x3a\x5d\x7d\x7b\xda\xf4\xc7\xa6\x87\xd4\x39\xd2\x74\xbe\xad\xe7\xcb\x86\xd5\xe5\xbf\x68\xdf\xfe\xfe\x2f\x89\xe8\xfd\x7f\x07\x00\x00\xff\xff\x16\xf9\xdb\xcc\x1e\x29\x00\x00")

func pr_review_js_bytes() ([]byte, error) {
	return bindata_read(
		_pr_review_js,
		"pr_review.js",
	)
}

func pr_review_js() (*asset, error) {
	bytes, err := pr_review_js_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "pr_review.js", size: 10526, mode: os.FileMode(420), modTime: time.Unix(1441709124, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	"pr_review.js": pr_review_js,
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

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"pr_review.js": &_bintree_t{pr_review_js, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

