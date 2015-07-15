package s3

import (
	"errors"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"net/url"
	"time"
)

var S3 *s3.S3
var awsProxyURL *url.URL

func init() {

	auth, err := aws.GetAuth(env.Get("AWS_ACCESS_KEY_ID"), env.Get("AWS_SECRET_ACCESS_KEY"))
	if err != nil {
		log.Criticalln(err)
	}

	regionName := env.Get("AWS_DEFAULT_REGION")
	S3 = s3.New(auth, aws.Regions[regionName])

	if awsURL := env.Get("AWS_PROXY_URL"); len(awsURL) > 0 {

		var err error
		awsProxyURL, err = url.Parse(awsURL)
		if err != nil {
			log.Criticalln("happyngine.s3.init():", err)
		}
	}
}

const (
	Private           = s3.ACL("private")
	PublicRead        = s3.ACL("public-read")
	PublicReadWrite   = s3.ACL("public-read-write")
	AuthenticatedRead = s3.ACL("authenticated-read")
	BucketOwnerRead   = s3.ACL("bucket-owner-read")
	BucketOwnerFull   = s3.ACL("bucket-owner-full-control")
)

func Put(bucket, path string, data []byte, contentType string, perm s3.ACL) error {

	b := S3.Bucket(bucket)
	if b == nil {
		return errors.New("Unknown bucket: " + bucket)
	}

	return b.Put(path, data, contentType, perm)
}

func SignedURL(bucket, path string, expires time.Time) string {

	b := S3.Bucket(bucket)
	return b.SignedURL(path, expires)
}

func Get(bucket, path string) ([]byte, error) {

	b := S3.Bucket(bucket)
	if b == nil {
		return nil, errors.New("Unknown bucket: " + bucket)
	}
	return b.Get(path)
}

func Del(bucket, path string) error {

	b := S3.Bucket(bucket)
	if b == nil {
		return errors.New("Unknown bucket: " + bucket)
	}
	return b.Del(path)
}

func Url(bucket, path string) string {

	u := new(url.URL)

	if awsProxyURL != nil {

		*u = *awsProxyURL

		u.Path = fmt.Sprintf("%s/%s", bucket, url.QueryEscape(path))

	} else {

		u.Scheme = "https"
		u.Host = bucket
		u.Path = url.QueryEscape(path)
	}

	return u.String()
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
