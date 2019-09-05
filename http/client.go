package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"github.com/tinwoan-go/basic-api/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	client          *http.Client
	contentType     = "Content-Type"
	jsonContentType = "application/json"
	xmlContentType  = "text/xml;charset=UTF-8"
)

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

func PostJSON(url, username, password string, request, response interface{}) error {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}
	req.Header.Set(contentType, jsonContentType)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &response)
}

func GetJSON(url, username, password string, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set(contentType, jsonContentType)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
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

	return json.NewDecoder(res.Body).Decode(&response)
}

func PostXML(url, username, password string, request, response interface{}) error {
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqXML))
	if err != nil {
		return err
	}
	req.Header.Set(contentType, xmlContentType)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, &response)
}

func GetXML(url, username, password string, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set(contentType, xmlContentType)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
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

	return json.NewDecoder(res.Body).Decode(&response)
}

// End-of-file
