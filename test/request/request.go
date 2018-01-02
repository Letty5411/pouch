package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/alibaba/pouch/client"
	"github.com/alibaba/pouch/pkg/utils"
	"github.com/alibaba/pouch/test/environment"
)

// Option defines a type used to update http.Request.
type Option func(*http.Request) error

// WithHeader sets the Header of http.Request.
func WithHeader(key string, value string) Option {
	return func(r *http.Request) error {
		r.Header.Add(key, value)
		return nil
	}
}

// WithQuery sets the query field in URL.
func WithQuery(query url.Values) Option {
	return func(r *http.Request) error {
		r.URL.RawQuery = query.Encode()
		return nil
	}
}

// WithJSONBody encodes the input data to JSON and sets it to the body in http.Request
func WithJSONBody(obj interface{}) Option {
	return func(r *http.Request) error {
		b := bytes.NewBuffer([]byte{})

		if obj != nil {
			err := json.NewEncoder(b).Encode(obj)

			if err != nil {
				return err
			}
		}
		r.Body = ioutil.NopCloser(b)
		r.Header.Set("Content-Type", "application/json")
		return nil
	}
}

// DecodeBody decodes body to obj.
func DecodeBody(obj interface{}, body io.ReadCloser) error {
	// TODO: this fuction could only be called once
	defer body.Close()
	return json.NewDecoder(body).Decode(obj)
}

// Delete sends request to the default pouchd server with custom request options.
func Delete(endpoint string, opts ...Option) (*http.Response, error) {
	apiClient, err := newAPIClient(environment.PouchdAddress, environment.TLSConfig)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(http.MethodDelete, apiClient.BaseURL()+endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return apiClient.HTTPCli.Do(req)
}

// Get sends request to the default pouchd server with custom request options.
func Get(endpoint string, opts ...Option) (*http.Response, error) {
	apiClient, err := newAPIClient(environment.PouchdAddress, environment.TLSConfig)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(http.MethodGet, apiClient.BaseURL()+endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return apiClient.HTTPCli.Do(req)
}

// Post sends post request to pouchd.
func Post(endpoint string, opts ...Option) (*http.Response, error) {
	apiClient, err := newAPIClient(environment.PouchdAddress, environment.TLSConfig)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(http.MethodPost, apiClient.BaseURL()+endpoint, opts...)
	if err != nil {
		return nil, err
	}

	// By default, if Content-Type in header is not set, set it to application/json
	if req.Header.Get("Content-Type") == "" {
		WithHeader("Content-Type", "application/json")(req)
	}
	return apiClient.HTTPCli.Do(req)
}

// newAPIClient return new HTTP client with tls.
//
// FIXME: Could we make some functions exported in alibaba/pouch/client?
func newAPIClient(host string, tls utils.TLSConfig) (*client.APIClient, error) {
	return client.NewAPIClient(host, tls)
}

// newRequest creates request targeting on specific host/path by method.
func newRequest(method, url string, opts ...Option) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err := opt(req)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

//func Hijack(endpoint string, opts ...Option) (net.Conn, *bufio.Reader, error) {
//	req, err := newRequest(http.MethodPost, environment.PouchdAddress+endpoint, opts...)
//	if err != nil {
//		return nil,nil, err
//	}
//	req.Header.Set("Connection", "Upgrade")
//	req.Header.Set("Upgrade", "tcp")
//
//	req.Host = environment.PouchdAddress
//	defaultTimeout := time.Second * 10
//	conn, err := net.DialTimeout("unix", environment.PouchdAddress, defaultTimeout)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	clientconn := httputil.NewClientConn(conn, nil)
//	defer clientconn.Close()
//
//	if _, err := clientconn.Do(req); err != nil {
//		return nil, nil, err
//	}
//
//	rwc, br := clientconn.Hijack()
//
//	return rwc, br, nil
//}
