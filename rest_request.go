package httpc

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

// NewRESTRequest creates a builder for REST API style request that use json as request and response body
func NewRESTRequest(c *Client, method Method, endpointPath string, args ...SetRequestOptionFn) *RESTRequest {
	if len(args) == 0 {
		args = make([]SetRequestOptionFn, 0)
	}
	return &RESTRequest{
		Id:           uuid.New().String(),
		client:       c,
		method:       method,
		endpointPath: endpointPath,
		args:         args,
	}
}

type RESTRequest struct {
	Id           string
	client       *Client
	method       Method
	endpointPath string
	args         []SetRequestOptionFn
}

func (rr *RESTRequest) AddOption(fn ...SetRequestOptionFn) *RESTRequest {
	rr.args = append(rr.args, fn...)
	return rr
}

func (rr *RESTRequest) AddHeader(args ...string) *RESTRequest {
	rr.args = append(rr.args, AddHeader(args...))
	return rr
}

func (rr *RESTRequest) AddQuery(args ...string) *RESTRequest {
	rr.args = append(rr.args, AddQuery(args...))
	return rr
}

func (rr *RESTRequest) Body(b interface{}) *RESTRequest {
	rr.args = append(rr.args, SetJsonBody(b))
	return rr
}

func (rr *RESTRequest) PreRequest(fn PreRequestFn) *RESTRequest {
	rr.args = append(rr.args, PreRequest(fn))
	return rr
}

// Do prepare REST request, do and parse response body to JSON dst
func (rr *RESTRequest) Do(ctx context.Context, dst interface{}) (*http.Response, error) {
	// Set "accept" header to Json mime type
	rr.AddHeader("Accept", MimeTypeJson)
	// Set request id in context
	ctx = context.WithValue(ctx, ContextRequestId, rr.Id)
	// Do request
	resp, respBody, err := rr.client.DoRequest(ctx, rr.method, rr.endpointPath, rr.args...)
	if err != nil {
		return nil, err
	}
	// Skip parsing json body if destination is nil or body is nil
	if dst == nil || len(respBody) == 0 {
		return resp, nil
	}
	// Parse response body as json
	err = json.Unmarshal(respBody, dst)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
