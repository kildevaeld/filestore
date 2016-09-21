package filesystem

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFs(t *testing.T) {

	fss, err := New(Options{
		Path: "./filesystem_test.go",
	})

	assert.NotNil(t, err)
	assert.Nil(t, fss)

	fss, err = New(Options{
		Path: "./path",
	})

	assert.Nil(t, err)
	assert.NotNil(t, fss)

	os.RemoveAll(fss.(*fs_impl).path)

}

func TestSet(t *testing.T) {

	fss, err := New(Options{
		Path: "./path",
	})

	assert.Nil(t, err)

	defer os.RemoveAll(fss.(*fs_impl).path)

	read := bytes.NewBuffer(nil)
	key := []byte("test/test.txt")

	read.WriteString("Hello, World")
	err = fss.Set(key, read, nil)

	assert.Nil(t, err)

	bs, berr := ioutil.ReadFile(filepath.Join(fss.(*fs_impl).path, string(key)))

	assert.Nil(t, berr)
	assert.Equal(t, "Hello, World", string(bs))

}

func TestGet(t *testing.T) {

	fss, err := New(Options{
		Path: "./path",
	})

	assert.Nil(t, err)

	defer os.RemoveAll(fss.(*fs_impl).path)

	writer := bytes.NewBuffer(nil)
	key := []byte("test/test.txt")

	writer.WriteString("Hello, World")
	err = fss.Set(key, writer, nil)

	assert.Nil(t, err)

	read, rerr := fss.Get(key)

	assert.Nil(t, rerr)

	b, e := ioutil.ReadAll(read)
	assert.Nil(t, e)
	assert.Equal(t, "Hello, World", string(b))

}
