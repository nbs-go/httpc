package httpc_test

import (
	"context"
	"encoding/json"
	"github.com/nbs-go/httpc"
	"net/http"
	"net/url"
	"testing"
)

type HttpBinResult struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Origin  string            `json:"origin"`
	Headers url.Values        `json:"headers"`
	Args    url.Values        `json:"args"`
	Data    string            `json:"data"`
	Files   json.RawMessage   `json:"files"`
	Form    url.Values        `json:"form"`
	Json    map[string]string `json:"json"`
}

type HttpEchoApiResult struct {
	Header map[string]string `json:"header"`
}

func TestRestGet(t *testing.T) {
	req := httpc.NewRESTRequest(c, "GET", "/get")
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "https://httpbin.nbs.dev/get"
	if respBody.Url != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, respBody.Url)
	}
}

func TestRestGetQuery(t *testing.T) {
	req := httpc.NewRESTRequest(c, "GET", "/anything").
		AddQuery("message", "hello")
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "hello"
	actual, _ := respBody.Args["message"]
	if len(actual) == 1 && actual[0] != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestRestPost(t *testing.T) {
	req := httpc.NewRESTRequest(c, "POST", "/post")
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "https://httpbin.nbs.dev/post"
	if respBody.Url != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, respBody.Url)
	}
}

func TestRestPostBody(t *testing.T) {
	req := httpc.NewRESTRequest(c, "POST", "/anything").
		Body(map[string]string{"message": "hello"})
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "hello"
	actual, _ := respBody.Json["message"]
	if actual != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestRestSkipResponseBodyParsing(t *testing.T) {
	req := httpc.NewRESTRequest(c, "POST", "/anything").
		Body(map[string]string{"message": "hello"})
	_, err := req.Do(context.Background(), nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestRestPostEmptyBody(t *testing.T) {
	req := httpc.NewRESTRequest(c, "POST", "/anything").
		Body(nil)
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := 0
	actual := len(respBody.Json)
	if actual != expected {
		t.Errorf("unexpected actual value. Expected = %d, Actual = %d", expected, actual)
	}
}

func TestRestAddOption(t *testing.T) {
	req := httpc.NewRESTRequest(c, "GET", "/anything").
		AddOption(httpc.AddQuery("message", "hello"))
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "hello"
	actual, _ := respBody.Args["message"]
	if len(actual) == 1 && actual[0] != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestRestXMLResponse(t *testing.T) {
	req := httpc.NewRESTRequest(c, "GET", "/xml")
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err.Error() != `invalid character '<' looking for beginning of value` {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestRestUrlEncodedForm(t *testing.T) {
	form := make(url.Values)
	form.Add("message", "hello")
	req := httpc.NewRESTRequest(c, "POST", "/anything").AddOption(httpc.SetUrlEncodedFormBody(form))
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert value
	expected := "hello"
	actual, _ := respBody.Form["message"]
	if len(actual) == 1 && actual[0] != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestRestPreRequest(t *testing.T) {
	dc := httpc.NewClient("https://httpbin.nbs.dev", httpc.Namespace("httpc_dump"), httpc.LogDump(true))
	req := httpc.NewRESTRequest(dc, "GET", "/anything").
		PreRequest(func(r *http.Request, rb []byte) {
			// Add header
			r.Header.Add("signature", "some-random-string")
			// Add query
			r.URL.RawQuery += "&message=hello"
		})
	// Do request
	var respBody HttpBinResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert header
	if v, _ := respBody.Headers["Signature"]; len(v) == 1 && v[0] != "some-random-string" {
		t.Errorf("unexpected condition: signature does not passed to headers. Actual = %s", v)
		return
	}
	// Assert query
	if v, _ := respBody.Args["message"]; len(v) == 1 && v[0] != "hello" {
		t.Errorf("unexpected condition: message mismatch")
		return
	}
}

func TestRestUpperCaseHeader(t *testing.T) {
	dc := httpc.NewClient("https://echo-api.nbs.dev", httpc.Namespace("echo-api"), httpc.LogDump(true))
	req := httpc.NewRESTRequest(dc, "GET", "/", httpc.DisableCanonicalHeader()).
		AddHeader("MESSAGE", "hello")
	// Do request
	var respBody HttpEchoApiResult
	_, err := req.Do(context.Background(), &respBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	// Assert result
	actual, _ := respBody.Header["MESSAGE"]
	if actual != "hello" {
		t.Errorf("unexpected value: %s", actual)
	}
}
