package s3

/// A Buffered ReadWriter
import (
	"errors"
	"io"
	"sync"
)

type writereader struct {
	buffer    []byte
	lock      sync.Mutex
	closed    bool
	cond      *sync.Cond
	readReady bool
}

func (self *writereader) Write(bs []byte) (int, error) {
	self.lock.Lock()
	self.buffer = append(self.buffer, bs...)
	self.lock.Unlock()

	self.cond.L.Lock()
	self.readReady = true
	self.cond.L.Unlock()
	self.cond.Signal()

	return len(bs), nil
}

func (self *writereader) Read(bs []byte) (int, error) {
	l := len(bs)

	if l == 0 {
		return 0, errors.New("")
	}

	if self.Length() > 0 {
		return self.read(bs)
	}

	self.cond.L.Lock()
	for !self.closed && !self.readReady {
		self.cond.Wait()
	}
	var err error
	if self.readReady {
		self.cond.L.Unlock()
		return self.read(bs)
	} else if self.closed {
		err = io.EOF
	}

	self.cond.L.Unlock()

	return 0, err
}

func (self *writereader) read(bs []byte) (n int, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	n = copy(bs, self.buffer[0:])
	self.buffer = self.buffer[n:]

	if len(self.buffer) == 0 {
		self.cond.L.Lock()
		self.readReady = false
		self.cond.L.Unlock()
		self.cond.Broadcast()
	}

	return n, nil
}

func (self *writereader) Length() int {
	self.lock.Lock()
	defer self.lock.Unlock()
	return len(self.buffer)
}

func (self *writereader) CloseWriter() error {
	self.cond.L.Lock()
	self.closed = true
	self.cond.L.Unlock()
	self.cond.Broadcast()
	return nil
}

func (self *writereader) Close() error {
	self.CloseWriter()

	self.cond.L.Lock()
	for self.readReady {
		self.cond.Wait()
	}
	self.cond.L.Unlock()

	return nil
}

func NewWriteReader() io.ReadWriteCloser {
	c := sync.NewCond(&sync.Mutex{})
	return &writereader{cond: c}
}

func newWriteReader() *writereader {
	c := sync.NewCond(&sync.Mutex{})
	return &writereader{cond: c}
}
