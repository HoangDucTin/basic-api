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

// This handler set the header to have no cache.
func setNoCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache")                                 // HTTP 1.0
		w.Header().Set("Expires", "0")                                       // Proxies
		next.ServeHTTP(w, r)
	})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger *logrus.Logger
}

func newLogMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	c := &StructuredLogger{
		Logger: logger,
	}
	return c.logMiddleware
}

func (l *StructuredLogger) logMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		entry := &StructuredLoggerEntry{Logger: l.Logger}
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
		switch {
		case strings.HasPrefix(ct, "application/json"):
			buf, _ := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
			var req interface{}
			_ = json.NewDecoder(rdr1).Decode(&req)
			logFields["RequestBody"] = req
		case strings.HasPrefix(ct, "text/xml"):
			buf, _ := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
			var req interface{}
			_ = xml.NewDecoder(rdr1).Decode(&req)
			logFields["RequestBody"] = req
		}
		if strings.HasPrefix(ct, "application/json") {

		}
		logFields["URI"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(loggingRW, r)
		var res interface{}
		_ = json.Unmarshal(loggingRW.body, &res)
		logFields["ResponseBody"] = res
		logFields["Status"] = loggingRW.status
		logFields["ProcessTime"] = fmt.Sprintf("%v", time.Since(start))
		entry.Logger = entry.Logger.WithFields(logFields).Logger
		entry.Logger.Info("")
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
