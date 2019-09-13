package handler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// This middleware will set the header
// of each and every request to store
// no cache.
func SetNoCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache")                                 // HTTP 1.0
		w.Header().Set("Expires", "0")                                       // Proxies
		next.ServeHTTP(w, r)
	})
}

// This middleware will print out
// the request and response of each
// and every request in JSON format
// (suitable for elastic search).
func NewLogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logFields := logrus.Fields{}
		start := time.Now()
		logFields["Start"] = start
		if reqID := middleware.GetReqID(r.Context()); reqID != "" {
			logFields["RequestID"] = reqID
		}
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		logFields["HttpMethod"] = r.Method
		logFields["RemoteAddr"] = r.RemoteAddr
		logFields["UserAgent"] = r.UserAgent()
		ct := r.Header.Get("Content-Type")
		isJson := strings.HasPrefix(ct, "application/json")
		isXml := strings.HasPrefix(ct, "text/xml")
		switch {
		case isJson:
			buf, _ := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
			var req interface{}
			_ = json.NewDecoder(rdr1).Decode(&req)
			logFields["RequestBody"] = req
		case isXml:
			buf, _ := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
			var req interface{}
			_ = xml.NewDecoder(rdr1).Decode(&req)
			logFields["RequestBody"] = req
		}
		logFields["URI"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(loggingRW, r)
		var res interface{}
		switch {
		case isJson:
			_ = json.Unmarshal(loggingRW.body, &res)
			logFields["ResponseBody"] = res
		case isXml:
			_ = xml.Unmarshal(loggingRW.body, &res)
			logFields["ResponseBody"] = res
		default:
			logFields["UnsupportedContentType"] = ct
		}
		logFields["Status"] = loggingRW.status
		logFields["ProcessTime"] = time.Since(start).String()
		entry := logrus.WithFields(logFields)
		entry.Logger.SetFormatter(&logrus.JSONFormatter{})
		entry.Println()
	}
	return http.HandlerFunc(fn)
}

type loggingResponseWriter struct {
	status int
	body   []byte
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}
