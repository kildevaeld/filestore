package filestore

import (
	"errors"
	"fmt"
	"io"
	"time"
)

var _stores map[string]func(o interface{}) (Store, error)

type Entry struct {
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

type File interface {
	io.ReadCloser
}

var (
	ErrNotFound = errors.New("Not Found")
)

type SetOptions struct {
	Size     int64
	MimeType string
}

type Options struct {
	Driver  string
	Options interface{}
}

type FileInfo interface {
	Name() string
	Size() int64
	ModTime() time.Time
}

type Store interface {
	Set(key []byte, reader io.Reader, o ...interface{}) error
	Stat(key []byte) (FileInfo, error)
	Get(key []byte) (File, error)
	Remove(key []byte) error
}

func Register(name string, fn func(o interface{}) (Store, error)) {
	_stores[name] = fn
}

func init() {
	_stores = make(map[string]func(o interface{}) (Store, error))
}

func New(o Options) (Store, error) {

	var (
		fn func(o interface{}) (Store, error)
		ok bool
	)

	if fn, ok = _stores[o.Driver]; !ok {
		return nil, fmt.Errorf("No driver named: '%s'. Have you forgotton to import?", o.Driver)
	}

	return fn(o.Options)
}
