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

var _pr_review_js = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5a\x5b\x73\xdb\xba\x11\x7e\xf7\xaf\x80\x4f\x66\x0e\xa9\x1c\x85\xbe\x3c\x5a\xe3\x5e\xe6\xa4\x99\xe9\x25\xb1\x27\x4e\xa6\x0f\x8a\x46\xa1\x45\xc8\x62\x2c\x11\x2a\x41\x5a\x55\x13\xfd\xf7\xee\xe2\x0e\x10\x92\xe5\xd8\x6d\x93\x9e\xf0\x21\x96\x88\xc5\xde\xb0\xfb\xed\x62\x95\x69\x5b\x4d\x9a\x92\x55\xe4\x86\x36\x6f\xe9\x5d\x49\x57\xb4\xe6\xaf\x58\xfd\xaa\x9c\xd3\x74\x0a\xff\xf0\x3e\x59\xb6\xf3\xf9\x5b\xfa\x8f\x96\xf2\xe6\x3d\xa7\x35\xbc\xa8\xe9\x94\xd6\x35\x2d\xcc\x8e\x1e\xf9\x7c\x40\xc8\x84\x55\x9c\xcd\x69\x36\x67\x37\x69\xf2\xaa\xac\x8a\xb2\xba\x21\xb5\xa6\x21\xd7\x6d\x43\xca\x9b\x8a\xd5\xf8\x3a\x21\xbf\x84\x8c\x7b\x83\x03\x60\x72\x97\xd7\x66\x0f\x6a\xc1\xc9\x39\xf9\xbc\x19\xa8\x15\x54\xc9\x48\x55\x2b\xb0\x34\xce\x68\x3e\x99\xa5\x5d\xc5\xfa\xc4\x58\xa8\xcd\xd1\xcc\xa5\xce\x66\xaf\x5a\xf5\xc9\x35\x0d\xf1\x05\x0f\xf1\xdb\x08\xc4\xc7\xde\x7e\xf9\x42\x86\xa3\xc1\xf6\x6d\xd9\xb2\xe5\xb3\xd4\x68\x21\x29\x37\xe2\xef\xa6\xe7\x5a\xb3\x53\x23\xd7\x4d\xe8\x88\xae\xe9\x42\x5a\x86\xff\x54\xf9\xc2\xea\x25\x76\x97\x53\x62\x34\xe0\x19\x50\xdc\x34\x33\x72\x7e\x7e\x4e\x8e\xad\xc5\x28\x61\x99\x37\x0d\xad\x2b\xe0\x3f\x06\x4e\x55\x91\xc2\x11\x4f\xcb\x9b\x4c\xbd\xef\x68\xf7\x2b\x2c\x5b\x0e\x1e\x0f\xae\xdc\x85\x24\x7a\xff\xc0\x23\x6c\x21\x06\x50\x55\x8f\xd2\x68\xa9\x34\xd7\xc4\x8b\xbc\x99\xcc\xac\x5e\x5a\x48\xc6\x97\xf3\xb2\x49\x93\x7e\xd2\x73\x75\x53\xcb\xae\x6a\x84\x1c\x1d\xb9\x01\xab\x48\xfa\xc4\x73\x5b\x6f\xe0\x6c\xa8\x69\xd3\x82\x33\x84\x68\x0c\xcd\xcb\xbc\x99\xdd\xbf\x4f\x1d\xaa\x96\x49\xf4\x01\x91\x29\x93\xf1\x4c\x26\x79\xd5\x90\x6b\x4a\x9a\x19\x25\x1c\x36\x93\x9c\x8b\xcf\x98\x1f\x20\x54\x24\x08\xe1\xb4\x2a\x68\x6d\x18\xe1\x09\x4a\x1f\xfc\xfc\x33\x49\x8d\xeb\x32\x70\x06\xfd\xe7\xc5\x34\x0d\x73\x4b\x9c\xee\x8b\x13\x8c\x02\x4b\xac\x0e\xfe\x77\xe4\xa4\xe7\xbb\x46\x59\x3a\xce\x56\x65\x33\x63\x6d\x63\x05\x74\xe0\xa0\x67\xb9\x1c\x3b\x46\x1f\x74\x8c\x47\x8d\x0f\x23\x07\x21\x65\xe9\xad\x1b\x4d\xee\x06\xb7\xd5\x43\xed\xb7\x61\xd1\xd5\xc7\x08\xc4\x30\xc9\xe7\x35\xcd\x8b\xb5\xf4\x39\xe2\x0e\x32\x2b\x2b\xe0\xc1\xa9\x88\x8d\xd4\xd1\x43\xb3\x34\xaf\xc6\xd9\x2d\x5d\xf3\xd4\x83\xa3\x9e\x5a\xf5\x0d\x0b\xe5\x38\x4e\xf1\x6d\xb5\x46\x85\x5b\xac\x07\x1c\x3f\x88\x34\xcf\xab\x82\x2d\x4c\xd8\x9c\x5b\x2e\xc3\xd7\x10\x81\xd9\x74\xce\x58\x9d\x8a\x8f\x92\x34\xed\x91\xe7\x24\x4c\xef\x9e\xce\x7d\xcf\x9a\xa1\xcf\x7c\xe4\x70\x8f\xaf\x3b\xe8\xb6\x93\x50\xa2\x5c\x37\x27\x36\x7d\x08\xed\x92\x4b\xef\xa9\x30\xf3\x18\x0d\x0e\xc0\x74\xb7\x2c\xe1\xdb\xab\x59\xfe\x3a\x5f\x4a\x44\x94\xfe\xd4\xc5\x40\xae\x98\x1a\xb1\x0f\x74\xda\x6d\x21\x44\x4a\xdc\xc9\xf8\x2c\xd7\x68\x6c\x74\xb4\x9b\x7c\x05\x9d\xe8\xbb\x58\xd2\x8a\x16\x29\xbd\xa3\x55\xf3\x32\x6f\x72\xab\x68\x10\xa2\x20\xc7\x10\x65\xb8\x36\x56\x49\x9e\x61\x9a\x21\x1a\x95\x95\xa9\x85\x53\x55\x03\xc1\x13\x97\x96\x8d\x70\x95\x15\x95\x55\xed\xe2\x5a\x15\x93\x58\x01\x7d\x48\x75\xff\xbc\xe9\x0d\x62\x0e\x8e\x1d\x85\xd0\x12\xe3\x3f\x9a\x29\x91\x24\x10\xd0\x4d\x39\xcf\x6f\x28\xb0\x4c\xfe\xb8\xa2\x9c\x01\xe0\xad\x58\x7d\x7b\x48\xde\xb0\x15\xf9\xd4\x22\xd6\x95\x80\x87\xf9\xe4\x96\x40\x4c\x91\x55\x0e\xdf\x10\x2b\xdf\xd5\xf9\x5d\x09\xc8\xc8\xa0\xa4\x70\x2e\xd6\x18\xc0\x64\x2d\x5e\x49\xc9\x64\xcd\xda\x1a\xfa\x90\x82\x66\x1f\xaa\x0f\x55\x22\x03\x55\x0b\xfc\x05\x24\x3e\x7b\x66\x00\x98\x0b\x02\xb7\x07\xf0\xd4\xdf\xa3\x75\xf0\x59\xff\x01\x7b\x1a\x4d\x02\x1f\x13\xa3\xc0\xbe\x4d\x86\xcf\xef\x85\x68\x92\x44\x85\xf0\x99\x6d\x4c\x75\x71\xc9\x0d\x81\xc1\x5c\x7f\xf5\x99\x31\x1d\x1c\xb4\x58\x80\xfb\x78\xd4\x43\x2f\x48\x3e\x99\xd0\x65\x73\x46\x3e\x2e\xeb\xb1\xf2\xeb\xc5\x5f\x3f\x6e\x23\x2e\x0a\x63\xb4\xb7\x05\x17\x86\xba\x6e\x8c\xc8\x50\x41\xf7\x68\x1b\xa3\xc9\x2c\xaf\xe0\x4b\x94\x97\x5a\x1b\xb2\x79\x31\x76\x58\x56\x74\x65\xbf\x6e\xe3\xbb\x84\x7e\x13\x42\xaa\xc9\x9b\x96\x7b\x4c\xe5\xab\x8f\x36\x0a\x26\x00\xc8\x0d\xfd\x33\xe7\x2d\xf4\x1e\x8b\x05\xe4\x56\x27\xc3\xfa\x9a\xb9\x3a\x02\xb9\xe5\x4a\x70\xd2\xa5\x64\x4b\x6e\xcf\x00\xed\x11\x5a\x74\x79\x49\x00\x2e\xb0\x43\x4e\xfa\xbb\xf7\xb5\xf5\xdc\x6c\xb9\x74\x3b\x02\x69\xc6\x99\x88\x92\x78\xfe\x7d\x62\x65\x05\xbd\x10\x81\x66\xc8\x08\xd5\xe6\x27\x07\xaa\x8a\xe1\x1f\x9e\xdf\x51\x94\xdc\xb1\x17\x23\xef\x0c\x33\x10\x14\x55\x1b\xc7\xd3\x9a\x2d\x80\xe9\x5f\xae\x2e\xde\x64\xbc\xc1\x6e\xbe\x9c\xae\x03\xd1\x12\xeb\xd1\xad\xbb\x59\x8b\x84\xe8\x32\xb3\xc0\x83\x9c\xb6\xe0\xed\xfb\x65\x01\xbe\x7f\x3a\xc0\x7d\x04\xde\x5a\x7d\xaf\x1a\x94\x36\x67\x79\x71\x9f\xd1\x66\x33\x04\xb5\x81\xa4\x7d\xb6\xc7\x8e\x23\x0e\xd9\xf6\xde\xe4\xca\xb0\xd7\xa6\xa6\x5e\x77\x6a\x22\xac\x8a\xc3\x58\xe6\xd0\x22\xa5\x9e\x61\xf2\x50\xa1\x5d\x85\xb6\x33\x45\xd8\xda\xf8\x4c\x02\x21\x0e\x9b\xc0\xc4\x38\x23\xed\x79\x6c\xe9\xe7\x00\x15\xbb\xd1\x12\xab\xce\xf6\x4a\x7e\x78\x6e\x6b\xb9\x7f\x89\xdb\xa7\xf2\xcb\xee\x46\xd4\xfd\xa6\x6e\xa9\x82\x55\x7b\x41\x7b\x6c\x89\x75\xdd\xf1\xe8\x1a\x7a\x85\x05\x94\x4d\xc5\x8d\x41\xba\x10\x36\x53\xd2\xca\xec\xd8\xbb\x0e\xfe\xff\x96\x41\xb8\x73\xbd\xa6\x35\x2c\x83\xdb\xb9\x72\x11\x76\x17\xae\xc9\xee\x89\xec\x65\x71\xd0\xfc\x6e\xef\x9f\xcd\xca\xef\x41\x56\x5b\xe1\x95\x23\x60\x1a\x10\xf6\xc8\x99\x54\x32\x30\xe3\x81\xd5\xe9\x47\x79\xfa\x1f\x97\xa7\xce\x21\xd9\xf2\x54\x16\x5e\x45\x2a\xf1\x4c\x95\x58\x0d\xd8\xf2\xba\xff\x06\xa7\x01\x2e\xa9\x7c\xed\x57\xac\x6b\x56\xac\x3d\xa2\x89\x94\x9c\xe1\x42\x06\x4a\xc3\x85\x50\xd3\x4e\x16\x85\x3b\x62\x12\x14\x6a\x64\x02\x89\xe3\xcd\x4c\x74\x28\x39\x78\xab\x5e\x99\x21\x83\x73\x5c\xbd\x70\x7c\x04\x92\x5c\x2c\x77\xee\xd4\x0b\xb8\x3e\x68\x4e\xbc\xbd\x96\x6e\x4d\x4f\x8e\x7b\x5a\x13\x12\x0e\x6f\xea\x06\xf8\x6a\x44\xc6\xaf\xda\x2a\xc8\x0e\x7b\x5b\xb7\x3b\x40\x76\x74\x9e\x71\x78\x28\xec\xd7\x6f\x37\x9d\xab\xbc\x37\xc1\x44\x0b\x30\x78\xe1\xaf\x0a\x55\x50\xcc\xd2\xf2\x55\x89\xe5\x0b\x56\x87\xc7\xa3\xac\x61\x7f\x63\x10\x5e\xbf\xe6\x50\xed\xbc\x61\xca\x04\xde\x90\x84\xdd\x26\x67\x8e\x3e\xe2\x84\x9b\x1a\x4f\xcb\x2d\xf7\x10\x17\x3b\x0b\xbc\x7e\xf0\x2c\xd4\x76\xdf\x4e\xb7\x1a\xfb\xd2\xc2\x72\xe5\x14\x67\xcd\x69\x10\x6c\x2b\xe8\x9c\x36\x34\x80\x29\x1b\x97\xa3\x90\xde\xa4\xd2\x0e\x33\xee\xcf\xcb\x88\xee\xcb\xfa\x25\x6d\xf2\x72\xde\x6d\xcb\xd4\x7b\x90\xd8\xd1\xfe\xde\x8a\x1a\xc4\xab\x7d\x3c\xd0\x34\xb2\x2d\x40\x92\x84\xb7\x70\x43\xe2\x08\x0c\x76\x19\x71\x30\x8a\x80\xa4\x60\x15\x45\x94\x73\x92\x25\x54\x36\x8a\xee\x65\xd1\x27\x3f\xd9\x1b\xdb\x12\x4f\xa3\x38\xfc\xa9\xb3\x79\x43\xe8\x1c\x62\xec\x1e\x43\x3a\xab\x84\x44\x6c\x8b\x50\x75\x0a\x41\x9c\x85\x53\x07\xbc\xdd\x8f\xae\x09\xbe\x2e\x7e\x7d\x70\x9f\xae\x5f\xfc\x58\xf2\xfb\x4d\x97\x4c\x7f\xba\x06\x7f\xdd\x0e\x82\xcc\x85\x6b\xac\x97\xba\x18\x57\x98\xf3\x27\xa3\x30\x76\x30\x54\xe1\x46\xaa\x33\x6c\x18\x24\xc8\x23\x73\x7e\x67\xd6\xc7\xf3\xfe\x01\x99\xdf\xd9\xa9\xcd\x3c\xed\x98\x29\x1f\xc7\x50\x83\xf1\xb1\x2b\x13\x64\xe6\xae\x46\xcf\x7d\xa2\xe3\x75\xa9\xc3\xae\xa9\xbc\x7e\x36\x91\xb7\x9b\xae\x61\x3e\x9c\xc9\x93\x1c\x09\x2b\x64\x73\x16\x5f\x17\x63\xcf\xbe\xb1\x3a\x22\xea\x69\xd0\x2f\x76\x80\x5f\x05\x7f\xbf\xf9\xec\xdf\x0a\xa9\xf2\x82\x22\x0f\x16\x4f\x0a\x07\x55\xb4\xe8\x62\xf2\x36\xbc\x70\x11\x63\x07\x7a\x1c\x3d\x0f\x80\xa4\xa6\x0b\x76\x47\x1f\x80\x25\xdf\x2c\x58\x84\x1b\xa3\x8d\x82\xca\x9c\xff\x4e\xa2\x3c\xf2\xac\x9e\x1f\x05\x67\x25\xe7\x8d\x5b\xce\x0a\x7f\x71\x8b\x43\xe3\x37\x7d\x6a\x51\x88\x8f\x9e\x59\x1c\xa0\xbb\xa4\xa7\xf7\x00\xe7\xa9\x05\xce\xb8\x9c\x18\x8c\x3f\x2c\x98\xfe\x53\xb8\xbb\x3f\x7a\xc8\x50\x29\xf0\x37\x08\xb3\x74\x1a\x33\x6d\xf3\xb4\x31\xab\x41\x85\x01\x40\x3b\x81\x2a\x66\x80\xf9\x2d\xfd\x13\xde\x06\x71\xd2\x27\xaf\x95\x67\x04\x95\x77\x2f\xee\x67\xdb\x6b\xc9\xc6\x4a\xe9\xfe\xbc\x65\x98\x3b\x26\x46\x55\x93\xe3\xf5\xef\xef\xda\x13\x8c\xd7\xfe\x2e\xb5\xd2\xcd\x3c\x6a\xe5\x8c\xd4\xec\xf3\xd5\xb3\x33\xfb\xec\x3b\x45\xeb\x08\xbd\x77\x9e\x16\x97\xb1\x75\xb2\x66\x9f\x58\x53\x15\x9d\xb6\xf9\x5b\x82\x57\x5b\x72\x29\x18\x56\xd9\xe7\xeb\xef\x6e\x4f\xd5\x2f\xfd\xb8\xf6\x7d\xb7\x8d\xdf\xe3\xae\x7d\xfa\x7f\x5f\xc8\x38\x06\x5e\x47\x47\xe4\xb2\x66\x78\xee\x64\xc6\xd8\x2d\x29\x00\x33\x0e\x30\x40\xc5\xbc\xed\xdd\x7a\x49\x45\x34\x26\x2e\xb6\x26\xd8\x22\xd8\x79\x5c\x2e\xb3\x52\x90\x31\x81\xa2\x89\x8c\xdd\x5d\xff\x79\x00\x84\xcb\x03\x8e\xc8\x12\x03\xc3\xb1\x1a\xf3\xed\x10\x26\x23\xa1\x08\x28\xe4\xb8\xd1\xd5\xb7\xa3\x4d\x77\x76\xb9\x4b\x9d\x3d\x4d\xe7\xeb\x6a\x32\xab\x59\x55\xfe\x8b\x76\xed\xef\xfe\x98\x87\xde\xff\x77\x00\x00\x00\xff\xff\xc2\x4c\x69\x0d\x9f\x28\x00\x00")

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

	info := bindata_file_info{name: "pr_review.js", size: 10399, mode: os.FileMode(420), modTime: time.Unix(1441706166, 0)}
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

