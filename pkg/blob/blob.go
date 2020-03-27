package blob

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	"gocloud.dev/blob"
)

type BlobConnect interface {
	Upload(string, []byte) error
	Download(string) (bytes.Buffer, error)
	Delete(string) error
	IndexFileExists() bool
}

type Connect struct {
	Bucket *blob.Bucket
	ctx    context.Context
}

func NewBlobConnect(bucketUrl string) (BlobConnect, error) {

	ctx := context.Background()

	bucket, err := blob.OpenBucket(ctx, bucketUrl)
	if err != nil {
		return Connect{}, err
	}

	if path := hasPath(bucketUrl); path != "" {
		bucket = blob.PrefixedBucket(bucket, path)
	}
	return Connect{
		Bucket: bucket,
		ctx:    ctx,
	}, nil
}

func (c Connect) Upload(filename string, data []byte) error {

	w, err := c.Bucket.NewWriter(c.ctx, filename, nil)
	if err != nil {
		return fmt.Errorf("Failed to obtain writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("Failed to write to blob: %v", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("Failed to close writer: %s", err)
	}

	return nil
}

func (c Connect) Download(filename string) (bytes.Buffer, error) {

	var bf bytes.Buffer

	r, err := c.Bucket.NewReader(c.ctx, filename, nil)
	if err != nil {
		if errorContains(err, "code=NotFound") {
			return bf, fmt.Errorf("Failed to download %s: specified file does not exist.", filename)
		}
		return bf, fmt.Errorf("Failed to download %s: %v", filename, err)
	}

	_, err = r.WriteTo(&bf)
	if err != nil {
		return bf, fmt.Errorf("Failed to download %s: %v", filename, err)
	}

	if err = r.Close(); err != nil {
		return bf, fmt.Errorf("Failed to close reader: %s", err)
	}
	return bf, nil
}

func (c Connect) Delete(filename string) error {

	err := c.Bucket.Delete(c.ctx, filename)
	if err != nil {
		if errorContains(err, "code=NotFound") {
			return fmt.Errorf("Failed to delete %s: specified file does not exist.", filename)
		}
		return fmt.Errorf("Failed to delete %s: %v", filename, err)
	}
	return nil
}

func (c Connect) IndexFileExists() bool {

	exists, err := c.Bucket.Exists(c.ctx, "index.yaml")
	if err != nil {
		fmt.Printf("Unable to verify helm repository: %v\n", err)
	}
	return exists
}

func hasPath(bucketUrl string) string {

	var path string
	url, err := url.Parse(bucketUrl)
	if err != nil {
		fmt.Printf("Failed to parse %s\n", bucketUrl)
	}
	if len(url.Path) != 0 {
		path = url.Path[1:]
		if path != "" {
			// append / at the end of path if its missing
			if !(strings.HasSuffix(path, "/")) {
				return fmt.Sprintf("%s/", path)
			}
		}
	}
	return path
}

// errorContains check if error contains specific string
func errorContains(err error, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(err.Error(), sub) {
			return true
		}
	}
	return false
}
