package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/tinwoan-go/basic-api/tlog"
)

var (
	client          *http.Client
	contentType     = "Content-Type"
	jsonContentType = "application/json"
	xmlContentType  = "text/xml;charset=UTF-8"
	// ErrUnsupportedContentType points out the
	// unsupported Content-Type of the request.
	ErrUnsupportedContentType = errors.New("unsupported Content-Type")
	log                       tlog.Logger
)

type (
	// PostInfo contains the
	// information for doing
	// a POST request.
	PostInfo struct {
		Ctx      context.Context
		URL      string
		Username string
		Password string
		Request  interface{}
		Response interface{}
	}

	// GetInfo contains the
	// information for doing
	// a GET request.
	GetInfo struct {
		Ctx      context.Context
		URL      string
		Username string
		Password string
		Response interface{}
	}
)

func init() {
	log = tlog.WithPrefix("client")
}

// NewHTTPClient will create an
// instance of HTTP client base
// on the proxyURL and timeout
// to help your server call to
// other servers in HTTP or HTTPS.
func NewHTTPClient(proxyURL string, timeout time.Duration) error {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxyURL != "" {
		if proxy, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		} else {
			return err
		}
	}

	client = &http.Client{
		Timeout:   timeout,
		Transport: &transport,
	}

	return nil
}

// Close disconnects the
// HTTP instance and
// closes the connection.
func Close() {
	client.Transport = nil
	client.CloseIdleConnections()
}

// PostJSON sends a post request
// to the server with given url.
// Both request and response will be
// in JSON format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func PostJSON(postInfo *PostInfo) error {
	return post(postInfo, jsonContentType)
}

// GetJSON sends a get request
// to the server with given url.
// Both request and response will be
// in JSON format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func GetJSON(getInfo *GetInfo) error {
	return get(getInfo, jsonContentType)
}

// PostXML sends a post request
// to the server with given url.
// Both request and response will be
// in XML format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func PostXML(postInfo *PostInfo) error {
	return post(postInfo, xmlContentType)
}

// GetXML sends a get request
// to the server with given url.
// Both request and response will be
// in XML format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func GetXML(getInfo *GetInfo) error {
	return get(getInfo, xmlContentType)
}

// post creates a post request
// base on the Content-Type and
// uses 'client' to do the request
// and bases on the Content-Type
// to parse result in to the response.
func post(postInfo *PostInfo, ct string) error {
	var (
		b   []byte
		err error
	)
	switch ct {
	case jsonContentType:
		b, err = json.Marshal(postInfo.Request)
		if err != nil {
			return err
		}
	case xmlContentType:
		b, err = xml.Marshal(postInfo.Request)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}
	req, err := http.NewRequest(http.MethodPost, postInfo.URL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set(contentType, ct)
	if user, pass := postInfo.Username, postInfo.Password; user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.TErrorf(postInfo.Ctx, "Error when close response body, error: %v", err)
		}
	}()
	switch ct {
	case jsonContentType:
		return json.NewDecoder(res.Body).Decode(&postInfo.Response)
	case xmlContentType:
		return xml.NewDecoder(res.Body).Decode(&postInfo.Response)
	default:
		return ErrUnsupportedContentType
	}
}

// get creates a get request
// base on the Content-Type and
// uses 'client' to do the request
// and bases on the Content-Type
// to parse result in to the response.
func get(getInfo *GetInfo, ct string) error {
	req, err := http.NewRequest(http.MethodGet, getInfo.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set(contentType, ct)
	if user, pass := getInfo.Username, getInfo.Password; user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.TErrorf(getInfo.Ctx, "Error when close response body, error: %v", err)
		}
	}()
	switch ct {
	case jsonContentType:
		return json.NewDecoder(res.Body).Decode(&getInfo.Response)
	case xmlContentType:
		return xml.NewDecoder(res.Body).Decode(&getInfo.Response)
	default:
		return ErrUnsupportedContentType
	}
}

// End-of-file
