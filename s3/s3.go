package s3

import (
	"errors"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"io"
	"net"
	"time"
)

var S3 *s3.S3

var attempts = aws.AttemptStrategy{
	Min:   5,
	Total: 5 * time.Second,
	Delay: 200 * time.Millisecond,
}

func init() {

	auth, err := aws.GetAuth(env.Get("AWS_ACCESS_KEY_ID"), env.Get("AWS_SECRET_ACCESS_KEY"))
	if err != nil {
		log.Criticalln(err)
	}

	regionName := env.Get("AWS_DEFAULT_REGION")
	S3 = s3.New(auth, aws.Regions[regionName])
}

const (
	Private           = s3.ACL("private")
	PublicRead        = s3.ACL("public-read")
	PublicReadWrite   = s3.ACL("public-read-write")
	AuthenticatedRead = s3.ACL("authenticated-read")
	BucketOwnerRead   = s3.ACL("bucket-owner-read")
	BucketOwnerFull   = s3.ACL("bucket-owner-full-control")
)

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	switch err {
	case io.ErrUnexpectedEOF, io.EOF:
		return true
	}
	switch e := err.(type) {
	case *net.DNSError:
		return true
	case *net.OpError:
		switch e.Op {
		case "read", "write":
			return true
		}
	case *s3.Error:
		switch e.Code {
		case "InternalError", "NoSuchUpload", "NoSuchBucket":
			return true
		}
	}
	return false
}

func PutHeader(bucket, path string, data []byte, contentType string, headers map[string][]string, perm s3.ACL) error {

	b := S3.Bucket(bucket)
	if b == nil {
		return errors.New("Unknown bucket: " + bucket)
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["Content-Type"] = []string{contentType}

	var err error
	for attempt := attempts.Start(); attempt.Next(); {
		err := b.PutHeader(path, data, headers, perm)
		if !shouldRetry(err) {
			break
		}
	}

	return err
}

func Put(bucket, path string, data []byte, contentType string, perm s3.ACL) error {

	return PutHeader(bucket, path, data, contentType, nil, perm)
}

func SignedURL(bucket, path string, expires time.Time) string {

	b := S3.Bucket(bucket)
	signedUrl := b.SignedURL(path, expires)

	return signedUrl
}

func Get(bucket, path string) ([]byte, error) {

	b := S3.Bucket(bucket)
	if b == nil {
		return nil, errors.New("Unknown bucket: " + bucket)
	}
	var err error
	var data []byte
	for attempt := attempts.Start(); attempt.Next(); {
		data, err = b.Get(path)
		if !shouldRetry(err) {
			break
		}
	}

	return data, err
}

func Del(bucket, path string) error {

	b := S3.Bucket(bucket)
	if b == nil {
		return errors.New("Unknown bucket: " + bucket)
	}
	return b.Del(path)
}

func Url(bucket, path string) string {

	return fmt.Sprintf("%s/%s/%s", S3.Region.S3Endpoint, bucket, path)
}

func GetBucketContents(bucket string) (*map[string]s3.Key, error) {

	b := S3.Bucket(bucket)
	if b == nil {
		return nil, errors.New("Unknown bucket: " + bucket)
	}

	return b.GetBucketContents()
}

func Copy(bucket, oldPath, newPath string, perm s3.ACL) error {

	b := S3.Bucket(bucket)
	if b == nil {
		return errors.New("Unknown bucket: " + bucket)
	}

	return b.Copy(oldPath, newPath, perm)
}
