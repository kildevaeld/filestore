package s3

import "time"

type file_info struct {
	name  string
	size  int64
	mtime time.Time
}

func (self *file_info) Name() string {
	return self.name
}

func (self *file_info) Size() int64 {
	return self.size
}

func (self *file_info) ModTime() time.Time {
	return self.mtime
}
