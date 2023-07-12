package httpc_test

import (
	"context"
	"github.com/nbs-go/httpc"
	"os"
	"testing"
)

func TestEnableOtelInstrumentation(t *testing.T) {
	// Enable OpenTelemetry tracing
	err := os.Setenv("HTTPC_INSTRUMENTATION", httpc.InstrumentationOpenTelemetry)
	// Init client
	oc := httpc.NewClient("https://httpbin.nbs.dev", httpc.Namespace("httpc_otel"))
	_, _, err = oc.DoRequest(context.Background(), "POST", "/anything", httpc.SetJsonBody(nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
}
