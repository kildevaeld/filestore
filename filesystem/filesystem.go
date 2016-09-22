package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kildevaeld/filestore"
	"github.com/mitchellh/mapstructure"
)

type Options struct {
	Path string
}

type fs_impl struct {
	path string
}

func (self *fs_impl) Set(key []byte, reader io.Reader, o *filestore.SetOptions) error {

	fp := filepath.Join(self.path, string(key))
	dir := filepath.Dir(fp)
	if dir != self.path {
		if err := os.MkdirAll(dir, 0755); err != nil {

			return err
		}
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	//defer file.Close()

	_, err = io.Copy(file, reader)
	file.Close()
	return err

}

func (self *fs_impl) Get(key []byte) (filestore.File, error) {
	fp := filepath.Join(self.path, string(key))

	file, err := os.Open(fp)

	if err == os.ErrNotExist {
		return nil, filestore.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return file, nil
}

func (self *fs_impl) Remove(key []byte) error {
	if err := os.Remove(filepath.Join(self.path, string(key))); err != nil {
		if err == os.ErrNotExist {
			return filestore.ErrNotFound
		}
		return err
	}
	return nil
}

func init() {

	filestore.Register("filesystem", func(o interface{}) (filestore.Store, error) {
		var options Options
		var err error
		switch m := o.(type) {
		case map[string]interface{}:
			err = mapstructure.Decode(m, &options)
		case Options:
			options = m
		case string:
			options.Path = m
		}

		if err != nil {
			return nil, err
		}

		if options.Path == "" {
			return nil, errors.New("no path")
		}

		return New(options)
	})

}

func New(options Options) (filestore.Store, error) {
	var err error
	var stat os.FileInfo

	path := options.Path

	if !filepath.IsAbs(path) {

		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(cwd, path)
	}

	stat, err = os.Stat(options.Path)
	if err != nil {
		if err = os.MkdirAll(path, 0766); err != nil {
			return nil, err
		}
	}

	if stat != nil && !stat.IsDir() {
		return nil, fmt.Errorf("path '%s' exists, but is not a directory", options.Path)
	}

	fs := &fs_impl{path}

	return fs, nil
}
