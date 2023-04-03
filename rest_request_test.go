package httpc_test

import (
	"context"
	"github.com/nbs-go/httpc"
	"net/http"
	"net/url"
	"testing"
)

type HttpBinResult struct {
	Url     string            `json:"url"`
	Args    map[string]string `json:"args"`
	Json    map[string]string `json:"json"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	Data    string            `json:"data"`
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
	expected := "https://httpbin.org/get"
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
	if actual != expected {
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
	expected := "https://httpbin.org/post"
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
	if actual != expected {
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
	if actual != expected {
		t.Errorf("unexpected actual value. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestRestPreRequest(t *testing.T) {
	req := httpc.NewRESTRequest(c, "GET", "/anything").
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
	if v, _ := respBody.Headers["Signature"]; v != "some-random-string" {
		t.Errorf("unexpected condition: signature does not passed to headers. Actual = %s", v)
		return
	}
	// Assert query
	if v, _ := respBody.Args["message"]; v != "hello" {
		t.Errorf("unexpected condition: message mismatch")
		return
	}
}
