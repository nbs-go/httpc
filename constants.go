package httpc

const (
	HeaderContentType = "Content-Type"
)

const (
	MimeTypeJson = "application/json"
)

type ContextKey int8

const (
	ContextRequestId = iota + 1
)
