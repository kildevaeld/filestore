package memory

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/kildevaeld/filestore"
)

type Options struct {
	Path string
}

type memory_impl struct {
	files map[string][]byte
	lock  sync.RWMutex
}

func (self *memory_impl) Set(key []byte, reader io.Reader, o *filestore.SetOptions) error {
	b, e := ioutil.ReadAll(reader)
	if e != nil {
		return e
	}
	self.lock.Lock()
	self.files[string(key)] = b
	self.lock.Unlock()
	return nil
}

func (self *memory_impl) Get(key []byte) (filestore.File, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if b, ok := self.files[string(key)]; ok {
		return ioutil.NopCloser(bytes.NewReader(b)), nil
	}

	return nil, filestore.ErrNotFound
}

func (self *memory_impl) Remove(key []byte) error {
	if _, ok := self.files[string(key)]; ok {
		delete(self.files, string(key))
		return nil
	}

	return filestore.ErrNotFound

}

func New() filestore.Store {
	return &memory_impl{files: make(map[string][]byte)}
}

func init() {

	filestore.Register("memory", func(o interface{}) (filestore.Store, error) {
		return &memory_impl{files: make(map[string][]byte)}, nil
	})

}
