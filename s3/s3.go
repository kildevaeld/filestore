package s3

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kildevaeld/filestore"
	"github.com/mitchellh/mapstructure"
)

type Options struct {
	AccessKey       []byte
	SecretAccessKey []byte
	Bucket          string
	Region          string
	ACL             string
	Cache           string
}

type s3_impl struct {
	client *s3.S3
	o      Options
}

func (self *s3_impl) Set(key []byte, reader io.Reader, o *filestore.SetOptions) error {

	uploader := s3manager.NewUploaderWithClient(self.client)

	k := string(key)

	options := s3manager.UploadInput{
		Bucket: &self.o.Bucket,
		Key:    &k,
		Body:   reader,
	}

	if self.o.Cache != "" {
		options.CacheControl = &self.o.Cache
	}
	if self.o.ACL != "" {
		options.ACL = &self.o.ACL
	}

	if o != nil {
		if o.MimeType != "" {
			options.ContentType = &o.MimeType
		}

	}

	_, err := uploader.Upload(&options)

	if err != nil {
		return err
	}

	return nil
}

func (self *s3_impl) Get(key []byte) (filestore.File, error) {

	//i := s3.GetObjectInput{}
	k := string(key)
	out, err := self.client.GetObject(&s3.GetObjectInput{
		Bucket: &self.o.Bucket,
		Key:    &k,
	})

	if err != nil {
		if e, ok := err.(awserr.Error); ok {
			fmt.Printf("%s\n", e)
		}
		return nil, err
	}

	return out.Body, nil
}

func (self *s3_impl) Remove(key []byte) error {

	return filestore.ErrNotFound

}

func init() {

	filestore.Register("s3", func(o interface{}) (filestore.Store, error) {
		var options Options
		var err error

		switch m := o.(type) {
		case map[string]interface{}:
			err = mapstructure.Decode(m, &options)
		case Options:
			options = m
			//case *s3.S3:
			//	return &s3_impl{m, }, nil
		}
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

}
