package s3

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/filestore"
	"github.com/mitchellh/mapstructure"
)

type Options struct {
	AccessKey       string
	SecretAccessKey string
	Token           string
	Bucket          string
	Region          string
	ACL             string
	Cache           string
}

type SetOptions struct {
	ACL          string
	MimeType     string
	CacheControl string
	Size         int64
}

type s3_impl struct {
	client *s3.S3
	o      Options
}

func get_options(v interface{}) *SetOptions {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case dict.Map, map[string]interface{}:
		var out SetOptions
		err := mapstructure.Decode(o, &out)
		if err != nil {
			return &out
		}
	case *SetOptions:
		return o
	case SetOptions:
		return &o
	case filestore.SetOptions:
		return &SetOptions{
			MimeType: o.MimeType,
			Size:     o.Size,
		}
	case *filestore.SetOptions:
		return &SetOptions{
			MimeType: o.MimeType,
			Size:     o.Size,
		}
	}
	return nil
}

func (self *s3_impl) Set(key []byte, reader io.Reader, o ...interface{}) error {

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

	so := get_options(o)
	if so != nil {
		if so.MimeType != "" {
			options.ContentType = &so.MimeType
		}
		if so.ACL != "" {
			options.ACL = &so.ACL
		}
		if so.CacheControl != "" {
			options.CacheControl = &so.CacheControl
		}
	}

	_, err := uploader.Upload(&options)

	if err != nil {
		return err
	}

	return nil
}

func (self *s3_impl) Get(key []byte) (filestore.File, error) {

	k := string(key)
	out, err := self.client.GetObject(&s3.GetObjectInput{
		Bucket: &self.o.Bucket,
		Key:    &k,
	})

	if err != nil {
		if e, ok := err.(awserr.RequestFailure); ok {
			if e.StatusCode() == http.StatusNotFound {
				return nil, filestore.ErrNotFound
			}
		}
		if e, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("%s %s", e.Code(), e.Message())
		}
		return nil, err
	}

	return out.Body, nil
}

func (self *s3_impl) Remove(key []byte) error {

	k := string(key)
	_, e := self.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &self.o.Bucket,
		Key:    &k,
	})

	return e

}

func New(o Options) (filestore.Store, error) {

	if o.Bucket == "" {
		return nil, errors.New("s3: No bucket")
	}

	var creds *credentials.Credentials

	if string(o.AccessKey) == "" || o.SecretAccessKey == "" {
		creds = credentials.NewEnvCredentials()
	} else {
		creds = credentials.NewStaticCredentials(o.AccessKey, o.SecretAccessKey, o.Token)
	}

	if creds == nil {
		return nil, errors.New("s3: no auth")
	}

	cfg := aws.NewConfig().WithCredentials(creds)

	if o.Region != "" {
		cfg = cfg.WithRegion(o.Region)
	}

	client := s3.New(session.New(), cfg)

	c := &s3_impl{client, o}

	return c, nil
}

func init() {

	filestore.Register("s3", func(o interface{}) (filestore.Store, error) {
		var options Options
		var err error

		switch m := o.(type) {
		case map[string]interface{}, dict.Map:
			err = mapstructure.Decode(m, &options)
		case Options:
			options = m
		default:
			return nil, errors.New("s3: No options")

		}
		if err != nil {
			return nil, err
		}

		return New(options)
	})

}
