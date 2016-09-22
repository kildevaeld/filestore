package bolt

import (
	"io"

	"github.com/boltdb/bolt"
	"github.com/kildevaeld/filestore"
)

type Options struct {
	Bolt   *bolt.DB
	Bucket []byte
}

type bolt_impl struct {
	bolt *bolt.DB
}

func (self *bolt_impl) Set(key []byte, reader io.Reader, o *filestore.SetOptions) error {

	return nil
}

func (self *bolt_impl) Get(key []byte) (filestore.File, error) {

	return nil, filestore.ErrNotFound
}

func (self *bolt_impl) Remove(key []byte) error {

	return filestore.ErrNotFound

}

func New() filestore.Store {
	return &bolt_impl{files: make(map[string][]byte)}
}

func init() {

	filestore.Register("memory", func(o interface{}) (filestore.Store, error) {
		return &bolt_impl{files: make(map[string][]byte)}, nil
	})

}
