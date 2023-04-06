package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/nlogger/v2"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewClient(baseUrl string, args ...SetClientOptionsFn) *Client {
	// Evaluate options
	o := evaluateClientOptions(args)
	// Init logger
	cl := nlogger.Get().NewChild(logOption.WithNamespace(o.namespace))
	// Set default timeout
	c := &http.Client{}
	// Init client
	return &Client{
		baseUrl:    baseUrl,
		httpClient: c,
		log:        cl,
	}
}

type Client struct {
	baseUrl    string
	httpClient *http.Client
	log        nlogger.Logger
}

func (c *Client) DoRequest(ctx context.Context, method Method, endpointPath string, args ...SetRequestOptionFn) (*http.Response, []byte, error) {
	o := evaluateRequestOptions(args)
	return c.doRequest(ctx, method, endpointPath, o)
}

func (c *Client) composeRequestBody(method Method, o *requestOptions) ([]byte, error) {
	if method == MethodGet || o.body == nil {
		return nil, nil
	}
	// Check body is already composed to []byte
	body, ok := o.body.([]byte)
	if ok {
		return body, nil
	}
	// Compose body by encoding-type
	ct := o.header[HeaderContentType]
	switch ct {
	case MimeTypeJson:
		j, err := json.Marshal(o.body)
		if err != nil {
			return nil, fmt.Errorf("httpc: Failed to compose request body. ContentType = %s, Error = %w", ct, err)
		}
		return j, nil
	case MimeTypeUrlEncodedForm:
		// Check if type is url.Values
		form, fOk := o.body.(url.Values)
		if !fOk {
			return nil, fmt.Errorf("httpc: Unable to compose URL-Encoded Form, body is not url.Values type. Type = %T", o.body)
		}
		return []byte(form.Encode()), nil
	}
	c.log.Warnf("Unsupported Content-Type in %s in request body", ct)
	return nil, nil
}

func (c *Client) doRequest(ctx context.Context, method Method, endpointPath string, o *requestOptions) (*http.Response, []byte, error) {
	// Validate context
	if ctx == nil {
		return nil, nil, errors.New("httpc: ctx is required")
	}
	// Compose url
	ub := bytes.NewBufferString(c.baseUrl)
	ub.WriteString(endpointPath)
	// If query is set, then add query params
	if len(o.query) > 0 {
		ub.WriteString("?")
		ub.WriteString(o.query.Encode())
	}
	u := ub.String()
	// Compose request body
	reqBody, err := c.composeRequestBody(method, o)
	if err != nil {
		return nil, nil, err
	}
	// Set timeout
	var cancel context.CancelFunc
	hCtx := ctx
	if o.timeout > 0 {
		hCtx, cancel = context.WithTimeout(ctx, time.Duration(o.timeout)*time.Millisecond)
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	// Create request
	req, err := http.NewRequestWithContext(hCtx, method, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}
	// Set header
	for k, v := range o.header {
		req.Header.Set(k, v)
	}
	// Call pre-request hook if set
	if o.preRequest != nil {
		o.preRequest(req, reqBody)
	}
	// Do request
	t := time.Now()
	reqId := c.getRequestId(ctx)
	c.log.Debugf("HTTP Request  (Id=%s) URL=\"%s %s\" Header=%s Body=%s", reqId, method, u, composeHeaderLog(req.Header), reqBody)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("HTTP Request  (Id=%s) Failed to do request", logOption.Format(reqId), logOption.Error(err))
		return nil, nil, err
	}
	c.log.Debugf("HTTP Response (Id=%s) Status=%s, TimeElapsed=%s", reqId, resp.Status, time.Since(t))
	// Read response body
	defer func() {
		wErr := resp.Body.Close()
		if wErr != nil {
			c.log.Warnf("HTTP Response (Id=%s) Failed to close Body reader. ID = %s, Error = %s", reqId, wErr)
		}
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	c.log.Debugf("HTTP Response (Id=%s) Header=%s Body=%s", reqId, composeHeaderLog(resp.Header), respBody)
	return resp, respBody, nil
}

// getRequestId retrieve requestId value from context. If no requestId in context, then requestId wil be generated
func (c *Client) getRequestId(ctx context.Context) string {
	val := ctx.Value(ContextRequestId)
	reqId, ok := val.(string)
	if !ok {
		return uuid.New().String()
	}
	return reqId
}

func composeHeaderLog(header http.Header) string {
	if len(header) == 0 {
		return ""
	}
	s := bytes.NewBufferString("")
	for k := range header {
		s.WriteString(`("`)
		s.WriteString(k)
		s.WriteString(`"="`)
		s.WriteString(header.Get(k))
		s.WriteString(`"),`)
	}
	// Remove last char
	return strings.TrimSuffix(s.String(), ",")
}
