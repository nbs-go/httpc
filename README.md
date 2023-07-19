# httpc

[![Go Report Card](https://goreportcard.com/badge/github.com/nbs-go/httpc)](https://goreportcard.com/report/github.com/nbs-go/httpc)
[![GitHub license](https://img.shields.io/github/license/nbs-go/httpc)](https://github.com/nbs-go/httpc/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/nbs-go/httpc/branch/master/graph/badge.svg?token=NXJHYTA06I)](https://codecov.io/gh/nbs-go/httpc)

A tiny HTTP Client wrapper based on Go http package.

## Installing

> **WARNING**
>
> API is not yet stable and we might introduce breaking changes until we reached version v1.0.0. See [Breaking Changes](#breaking-changes) section for deprecation notes.

```shell
go get -u github.com/nbs-go/httpc
```

## Breaking Changes

### v0.7.0

- Revert Go minimum version to 1.17
- Remove built-in OpenTelemetry instrumentation, replaced with `SetGlobalTransporterOverrider()` function

### v0.6.0

- Upgrade Go minimum version to 1.19

## Usage

> TODO

### Enable OpenTelemetry Instrumentation

```
package main

import (
	"github.com/nbs-go/httpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// On init() package, SetGlobalTransporterOverrider to enable OpenTelemetry instrumentation
// across all initiated httpc.Client
func init() {
	httpc.SetGlobalTransporterOverrider(func(existingTransporter http.RoundTripper) http.RoundTripper {
		// Wrap existingTransporter with otelhttp.Transporter to enable instrumentation
		return otelhttp.NewTransport(existingTransporter)
	})
}
```

## Contributors

<a href="https://github.com/nbs-go/nsql/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=nbs-go/nsql" alt="contributors" />
</a>