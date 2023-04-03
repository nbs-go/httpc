package httpc

import (
	"errors"
	"net/http"
	"net/url"
)

type SetRequestOptionFn func(o *requestOptions)

type PreRequestFn func(req *http.Request, reqBody []byte)

func AddHeader(args ...string) SetRequestOptionFn {
	return func(o *requestOptions) {
		// Panic args length is 0 or odd
		count := len(args)
		if count == 0 || count%2 == 1 {
			panic(errors.New("httpc: Invalid AddHeader() args count must >= 2 and even"))
		}
		// Apply headers
		for i := 0; i < count; i += 2 {
			o.header[args[i]] = args[i+1]
		}
	}
}

func AddQuery(args ...string) SetRequestOptionFn {
	return func(o *requestOptions) {
		// Panic args length is 0 or odd
		count := len(args)
		if count == 0 || count%2 == 1 {
			panic(errors.New("httpc: Invalid AddQuery() args count must >= 2 and even"))
		}
		// Apply headers
		for i := 0; i < count; i += 2 {
			o.query.Add(args[i], args[i+1])
		}
	}
}

func SetBody(body interface{}) SetRequestOptionFn {
	return func(o *requestOptions) {
		if body == nil {
			return
		}
		// Set options
		o.body = body
	}
}

func SetJsonBody(body interface{}) SetRequestOptionFn {
	return func(o *requestOptions) {
		if body == nil {
			return
		}
		// Set options
		o.header[HeaderContentType] = MimeTypeJson
		o.body = body
	}
}

func Timeout(ms int) SetRequestOptionFn {
	return func(o *requestOptions) {
		o.timeout = ms
	}
}

func PreRequest(fn PreRequestFn) SetRequestOptionFn {
	return func(o *requestOptions) {
		o.preRequest = fn
	}
}

func SetUrlEncodedFormBody(body url.Values) SetRequestOptionFn {
	return func(o *requestOptions) {
		if body == nil {
			return
		}
		// Set options
		o.header[HeaderContentType] = MimeTypeUrlEncodedForm
		o.body = body
	}
}

type requestOptions struct {
	header     map[string]string
	query      url.Values
	body       interface{}
	timeout    int
	preRequest PreRequestFn
}

// evaluateClientOptions evaluates Client options and override default value
func evaluateRequestOptions(args []SetRequestOptionFn) *requestOptions {
	b := requestOptions{
		header:  make(map[string]string),
		query:   make(url.Values),
		body:    nil,
		timeout: 10000, // Set default timeout to 10 second
	}
	for _, fn := range args {
		fn(&b)
	}
	return &b
}
