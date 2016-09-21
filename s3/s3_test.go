package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestS3(t *testing.T) {

	creds := credentials.NewEnvCredentials()

	cfg := aws.NewConfig().WithRegion("eu-west-1").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	client := s3_impl{svc, Options{
		Bucket: "boellefesten",
		ACL:    s3.BucketCannedACLPublicRead,
	}}

	client.Set([]byte("/test/mig.txt"), bytes.NewReader([]byte("Hello, world")), nil)

	file, er := client.Get([]byte("/test/mig.txt"))
	if er != nil {
		t.Fatal(er)
	}

	b, _ := ioutil.ReadAll(file)
	fmt.Printf("%s\n", b)
}
