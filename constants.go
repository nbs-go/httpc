package httpc

const (
	HeaderContentType = "Content-Type"
)

const (
	MimeTypeJson           = "application/json"
	MimeTypeUrlEncodedForm = "application/x-www-form-urlencoded"
)

type ContextKey int8

const (
	ContextRequestId = iota + 1
)
