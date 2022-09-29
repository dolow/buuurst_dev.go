package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/dolow/buuurst_dev.go/collector"
)

func Middleware(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		writer := NewResponseWriter(w)
		h.ServeHTTP(writer, r)
		c := collector.New()
		c.Collect(writer, w, r, bodyBytes)
	})
}

func ServeHttp(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	writer := NewResponseWriter(w)
	c := collector.New()
	c.Collect(writer, w, r, bodyBytes)
}
