package s3

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/awslabs/aws-sdk-go/aws"
)

// needsMD5 returns true for any request that requires the Content-MD5 header
// be included with the request (and signed as a canonical header).
func needsMD5(r *aws.Request) bool {
	return r.Operation.Name == "DeleteObjects"
}

// contentMD5 computes and sets the HTTP Content-MD5 header for requests that
// require it.
func contentMD5(r *aws.Request) {
	if !needsMD5(r) {
		return
	}

	h := md5.New()

	// hash the body.  seek back to the first position after reading to reset
	// the body for transmission.  copy errors may be assumed to be from the
	// body.
	_, err := io.Copy(h, r.Body)
	if err != nil {
		r.Error = fmt.Errorf("content-md5: read: %v", err)
		return
	}
	_, err = r.Body.Seek(0, 0)
	if err != nil {
		r.Error = fmt.Errorf("content-md5: seek: %v", err)
		return
	}

	// encode the md5 checksum in base64 and set the request header.
	sum := h.Sum(nil)
	sum64 := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
	base64.StdEncoding.Encode(sum64, sum)
	r.HTTPRequest.Header.Set("Content-MD5", string(sum64))
}
