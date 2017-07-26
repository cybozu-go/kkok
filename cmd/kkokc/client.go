package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

var (
	client = &http.Client{}

	// ErrNotFound indicates kkok server returned 404.
	ErrNotFound = errors.New("not found")
)

// Call makes a REST API call to kkok server.
// method is HTTP method.  api is API request path such as /version.
// If j is []byte, it will be used as the request body.
// If j is not []byte nor nil, it will be encoded into JSON and
// be used as request body.
//
// If server returns 200, this returns the response body and nil.
// If server returns 404, this returns nil and ErrNotFound.
// For other status codes, this will return non-nil errors.
func Call(ctx context.Context, method, api string, j interface{}) ([]byte, error) {
	u := *kkokURL
	// to allow reverse proxies adding a path prefix.
	u.Path = path.Join(u.Path, api)
	header := make(http.Header)
	if method == "PUT" || method == "POST" {
		header.Set("Content-Type", "application/json")
	}
	if len(*flgToken) != 0 {
		header.Set("Authorization", "Bearer "+*flgToken)
	}
	var body io.ReadCloser
	var contentLength int64
	if j != nil {
		data, ok := j.([]byte)
		if !ok {
			var err error
			data, err = json.Marshal(j)
			if err != nil {
				return nil, errors.Wrap(err, "call "+api)
			}
		}
		body = ioutil.NopCloser(bytes.NewReader(data))
	}
	req := &http.Request{
		Method:        method,
		Host:          u.Host,
		URL:           &u,
		Header:        header,
		Body:          body,
		ContentLength: contentLength,
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "call "+api)
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "call "+api)
	}

	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("call %s: status %d: error %s",
			api, resp.StatusCode, string(data))
	}

	return data, nil
}
