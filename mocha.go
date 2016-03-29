package mocha

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

// request handling func type to replace gin.HandlerFunc
type RequestFunc func(*gin.Context)

// response handling func type
type ResponseFunc func(*httptest.ResponseRecorder)

type RequestConfig struct {
	Method      string
	Path        string
	Body        string
	Headers     map[string]string
	Middlewares []gin.HandlerFunc
	Handler     RequestFunc
	Callback    ResponseFunc
	Debug       bool
}

func New() *RequestConfig {

	return &RequestConfig{}
}

func (rc *RequestConfig) SetDebug(enable bool) *RequestConfig {
	rc.Debug = enable

	return rc
}

func (rc *RequestConfig) GET(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "GET"

	return rc
}

func (rc *RequestConfig) POST(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "POST"

	return rc
}

func (rc *RequestConfig) PUT(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "PUT"

	return rc
}

func (rc *RequestConfig) DELETE(path string) *RequestConfig {
	rc.Path = path
	rc.Method = "DELETE"

	return rc
}

func (rc *RequestConfig) SetHeader(headers map[string]string) *RequestConfig {
	if len(headers) > 0 {
		rc.Headers = headers
	}

	return rc
}

func (rc *RequestConfig) SetBody(body string) *RequestConfig {
	if len(body) > 0 {
		rc.Body = body
	}

	return rc
}

func (rc *RequestConfig) RunGinEngine(r *gin.Engine, response ResponseFunc) {
	qs := ""
	if strings.Contains(rc.Path, "?") {
		ss := strings.Split(rc.Path, "?")
		rc.Path = ss[0]
		qs = ss[1]
	}

	body := bytes.NewBufferString(rc.Body)

	req, _ := http.NewRequest(rc.Method, rc.Path, body)

	if len(qs) > 0 {
		req.URL.RawQuery = qs
	}

	if len(rc.Headers) > 0 {
		for k, v := range rc.Headers {
			req.Header.Set(k, v)
		}
	} else if rc.Method == "POST" || rc.Method == "PUT" {
		if strings.HasPrefix(rc.Body, "{") {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if rc.Debug {
		log.Printf("Request Method: %s", rc.Method)
		log.Printf("Request Path: %s", rc.Path)
		log.Printf("Request Body: %s", rc.Body)
		log.Printf("Request Headers: %s", rc.Headers)
		log.Printf("Request Header: %s", req.Header)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response(w)
}
