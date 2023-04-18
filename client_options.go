package httpc

type SetClientOptionsFn func(o *clientOptions)

type clientOptions struct {
	namespace    string
	logDump      bool
	disableHTTP2 bool
}

// Namespace override default Client namespace value
func Namespace(n string) SetClientOptionsFn {
	return func(o *clientOptions) {
		o.namespace = n
	}
}

// LogDump enable log HTTP request and response dump
func LogDump(enable bool) SetClientOptionsFn {
	return func(o *clientOptions) {
		o.logDump = enable
	}
}

// DisableHTTP2 disable HTTP/2 alternate protocol
func DisableHTTP2() SetClientOptionsFn {
	return func(o *clientOptions) {
		o.disableHTTP2 = true
	}
}

// evaluateClientOptions evaluates Client options and override default value
func evaluateClientOptions(args []SetClientOptionsFn) *clientOptions {
	o := clientOptions{
		namespace:    "httpc",
		logDump:      false,
		disableHTTP2: false,
	}
	for _, fn := range args {
		fn(&o)
	}
	return &o
}
