# CHANGELOG

## v0.7.0

- fix: Set Go runtime requirement to v1.17
- BREAKING CHANGE: Remove built-in OpenTelemetry dependencies. Use SetGlobalTransporterOverrider to wrap httpc.Client Transporter instead
- BREAKING CHANGE: Add SetGlobalTransporterOverrider to override initiated httpc.Client Transporter

## v0.6.2

- fix(log): Add context to include metadata from context

## v0.6.1

- fix(instrumentation): Add mutex to lock global instrumentation variable
- feat(instrumentation): Add env OTEL_TRACE_HTTPC to enable OpenTelemetry

## v0.6.0

- feat(instrumentation): Add OpenTelemetry instrumentation
- BREAKING CHANGE: Upgrade Go minimum version to 1.19

## v0.5.1

- fix: Skip parsing json if response body is empty

## v0.5.0

- feat: Add NonCanonicalHeader request options to preserve case
- feat: Add option to disable automatic switch to HTTP/2

## v0.4.1

- fix: Fix missing request body on log dump

## v0.4.0

- feat: Add DumpLog to log http request and response dump

## v0.3.1

- fix: Fix request and response log format and add Header in log

## v0.3.0

- feat: Add PreRequest hook

## v0.2.0

- feat: Support raw encoded []byte request body
- feat: Add request body encoding support for url-encoded form
- fix(httpc): Change header data type to map[string]string and pass header on do request

## v0.1.0

- feat(rest): Add REST API Request Builder
- feat: Add HTTP Client wrapper

