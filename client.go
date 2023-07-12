package httpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/nlogger/v2"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

var instrumentation Instrumentation

func init() {
	LoadEnv()
}

func LoadEnv() {
	// Set instrumentation using OpenTelemetry env
	ev := os.Getenv("OTEL_TRACE_HTTPC")
	if ev == "true" {
		instrumentation = InstrumentationOpenTelemetry
		return
	}
	// Set instrumentation using httpc
	ev = os.Getenv("HTTPC_INSTRUMENTATION")
	if ev == "opentelemetry" {
		instrumentation = InstrumentationOpenTelemetry
		return
	}
}

func NewClient(baseUrl string, args ...SetClientOptionsFn) *Client {
	// Evaluate options
	o := evaluateClientOptions(args)
	// Init logger
	cl := nlogger.Get().NewChild(logOption.WithNamespace(o.namespace))
	// Init Client
	c := &http.Client{}
	// Set transport
	if o.disableHTTP2 {
		c.Transport = &http.Transport{
			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		}
		cl.Debugf("HTTP/2 automatic switch is disabled")
	}
	// Wrap OpenTelemetry instrumentation
	if instrumentation == InstrumentationOpenTelemetry {
		c.Transport = otelhttp.NewTransport(c.Transport)
	}
	// Init client
	return &Client{
		baseUrl:    baseUrl,
		httpClient: c,
		log:        cl,
		logDump:    o.logDump,
	}
}

type Client struct {
	baseUrl    string
	httpClient *http.Client
	log        nlogger.Logger
	logDump    bool
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
	if o.canonicalHeader {
		for k, v := range o.header {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range o.header {
			req.Header[k] = []string{v}
		}
	}

	// Call pre-request hook if set
	if o.preRequest != nil {
		o.preRequest(req, reqBody)
	}
	// Do request
	t := time.Now()
	reqId := c.getRequestId(ctx)
	c.logDumpRequest(req, reqId)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error("HTTP Request  (Id=%s) Failed to do request", logOption.Format(reqId), logOption.Error(err))
		return nil, nil, err
	}
	c.logDumpResponse(resp, reqId)
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
	c.log.Debugf("HTTP Request  (Id=%s) URL=\"%s %s\" ResponseStatus=\"%s\" TimeElapsed=\"%s\"", reqId, req.Method, req.URL.String(), resp.Status, time.Since(t))
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

func (c *Client) logDumpRequest(req *http.Request, reqId string) {
	if !c.logDump {
		return
	}
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		c.log.Warnf("Unable to dump request. Error = %s", err)
		return
	}
	c.log.Debugf("\n---------- HTTP Request Dump -----------\n(RequestId=%s)\n%s\n----------------------------------------", reqId, dump)
}

func (c *Client) logDumpResponse(resp *http.Response, reqId string) {
	if !c.logDump {
		return
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		c.log.Warnf("Unable to dump response. Error = %s", err)
		return
	}
	c.log.Debugf("\n---------- HTTP Response Dump ----------\n(RequestId=%s)\n%s\n----------------------------------------", reqId, dump)
}
