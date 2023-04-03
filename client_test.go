package httpc_test

import (
	"context"
	"encoding/json"
	"github.com/nbs-go/httpc"
	"net/http"
	"testing"
)

var c *httpc.Client

func init() {
	c = httpc.NewClient("https://httpbin.org", httpc.Namespace("httpc_test"))
}

func TestTimeout(t *testing.T) {
	ct := httpc.NewClient("https://httpbin.org")
	_, _, err := ct.DoRequest(context.Background(), "GET", "/delay/1", httpc.Timeout(100))
	if err.Error() != `Get "https://httpbin.org/delay/1": context deadline exceeded` {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestInvalidJsonBody(t *testing.T) {
	_, _, err := c.DoRequest(context.Background(), "POST", "/post", httpc.SetJsonBody(json.RawMessage("{")))
	if err.Error() != `httpc: Failed to compose request body. ContentType = application/json, Error = json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input` {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestNilContext(t *testing.T) {
	_, _, err := c.DoRequest(nil, "POST", "/post")
	if err.Error() != `httpc: ctx is required` {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestUnimplementedBody(t *testing.T) {
	resp, _, _ := c.DoRequest(context.Background(), "POST", "/post",
		httpc.AddHeader(httpc.HeaderContentType, "application/xml"),
		httpc.SetBody(map[string]string{"message": "hello"}))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected response status code. StatusCode = %d", resp.StatusCode)
		return
	}
}

func TestInvalidAddHeaderArgs(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("unexpected condition: Code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("unexpected error: Unknown recovered err value. Error = %v", err)
			return
		}
		expected := `httpc: Invalid AddHeader() args count must >= 2 and even`
		if err.Error() != expected {
			t.Errorf("unexpected error: %s", err)
		}
	}()
	_, _, _ = c.DoRequest(context.Background(), "POST", "/post", httpc.AddHeader("key1", "value1", "key1"))
}

func TestInvalidAddQueryArgs(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("unexpected condition: Code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("unexpected error: Unknown recovered err value. Error = %v", err)
			return
		}
		expected := `httpc: Invalid AddQuery() args count must >= 2 and even`
		if err.Error() != expected {
			t.Errorf("unexpected error: %s", err)
		}
	}()
	_, _, _ = c.DoRequest(context.Background(), "POST", "/post", httpc.AddQuery("key1", "value1", "key1"))
}

func TestEmptyBody(t *testing.T) {
	_, respBody, err := c.DoRequest(context.Background(), "POST", "/anything", httpc.SetJsonBody(nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Convert response body
	var result HttpBinResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if result.Data != "" {
		t.Errorf("unexpected condition: No data shall pass to request body")
		return
	}
}