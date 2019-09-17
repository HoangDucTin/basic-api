package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/tinwoan-go/basic-api/logger"
)

var (
	client          *http.Client
	contentType     = "Content-Type"
	jsonContentType = "application/json"
	xmlContentType  = "text/xml;charset=UTF-8"
	// ErrUnsupportedContentType points out the
	// unsupported Content-Type of the request.
	ErrUnsupportedContentType = errors.New("unsupported Content-Type")
)

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
func PostJSON(url, username, password string, request, response interface{}) error {
	return post(url, jsonContentType, username, password, request, response)
}

// GetJSON sends a get request
// to the server with given url.
// Both request and response will be
// in JSON format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func GetJSON(url, username, password string, response interface{}) error {
	return get(url, jsonContentType, username, password, response)
}

// PostXML sends a post request
// to the server with given url.
// Both request and response will be
// in XML format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func PostXML(url, username, password string, request, response interface{}) error {
	return post(url, xmlContentType, username, password, request, response)
}

// GetXML sends a get request
// to the server with given url.
// Both request and response will be
// in XML format.
// If username or password is empty,
// this function will send a simple
// request with no authentication.
func GetXML(url, username, password string, response interface{}) error {
	return get(url, xmlContentType, username, password, response)
}

// post creates a post request
// base on the Content-Type and
// uses 'client' to do the request
// and bases on the Content-Type
// to parse result in to the response.
func post(url, ct, user, pass string, request, response interface{}) error {
	var (
		b   []byte
		err error
	)
	switch ct {
	case jsonContentType:
		b, err = json.Marshal(request)
		if err != nil {
			return err
		}
	case xmlContentType:
		b, err = xml.Marshal(request)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set(contentType, ct)
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("Error when close response body, error: %v", err)
		}
	}()
	switch ct {
	case jsonContentType:
		return json.NewDecoder(res.Body).Decode(&response)
	case xmlContentType:
		return xml.NewDecoder(res.Body).Decode(&response)
	default:
		return ErrUnsupportedContentType
	}
}

// get creates a get request
// base on the Content-Type and
// uses 'client' to do the request
// and bases on the Content-Type
// to parse result in to the response.
func get(url, ct, user, pass string, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set(contentType, ct)
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("Error when close response body, error: %v", err)
		}
	}()
	switch ct {
	case jsonContentType:
		return json.NewDecoder(res.Body).Decode(&response)
	case xmlContentType:
		return xml.NewDecoder(res.Body).Decode(&response)
	default:
		return ErrUnsupportedContentType
	}
}

// End-of-file
