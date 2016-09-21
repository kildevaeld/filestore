package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kildevaeld/filestore"
)

func TestS3(t *testing.T) {

	/*creds := credentials.NewEnvCredentials()

	cfg := aws.NewConfig().WithRegion("eu-west-1").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)*/

	/*client := s3_impl{svc, Options{
		Bucket: "boellefesten",
		ACL:    s3.BucketCannedACLPublicRead,
	}}*/
	client, _ := New(Options{
		Bucket: "boellefesten",
		Region: s3.BucketLocationConstraintEuWest1,
	})

	client.Set([]byte("/test/mig.txt"), bytes.NewReader([]byte("Hello, world")), &filestore.SetOptions{
		MimeType: "text/plain",
	})

	file, er := client.Get([]byte("/test/mig.txt"))
	if er != nil {
		t.Fatal(er)
	}

	b, _ := ioutil.ReadAll(file)

	client.Remove([]byte("/test/mig.txt"))
	fmt.Printf("%s\n", b)
}
