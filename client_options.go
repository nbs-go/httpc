package httpc

type SetClientOptionsFn func(o *clientOptions)

type clientOptions struct {
	namespace string
}

// Namespace override default Client namespace value
func Namespace(n string) SetClientOptionsFn {
	return func(o *clientOptions) {
		o.namespace = n
	}
}

// evaluateClientOptions evaluates Client options and override default value
func evaluateClientOptions(args []SetClientOptionsFn) *clientOptions {
	o := clientOptions{
		namespace: "httpc",
	}
	for _, fn := range args {
		fn(&o)
	}
	return &o
}
