package apollo

// Option apollo option
type Option func(opts *Options)

// Options apollo options
type Options struct {
	// server ip addr
	Addr string
	// application id
	AppId string
	// cluster environment
	Cluster string
	// namespace name
	NamespaceNames []string
	// access key secret
	Secret string
}

// WithAddr Set apollo server ip adder
func WithAddr(addr string) Option {
	return func(opts *Options) { opts.Addr = addr }
}

// WithAppId Set apollo appId
func WithAppId(appId string) Option {
	return func(opts *Options) { opts.AppId = appId }
}

// WithClusterSet apollo cluster
func WithCluster(cluster string) Option {
	return func(opts *Options) { opts.Cluster = cluster }
}

// WithNamespaceNames set apollo namespaceNames
func WithNamespaceNames(namespaceNames []string) Option {
	return func(opts *Options) { opts.NamespaceNames = namespaceNames }
}

// WithSecret set apollo access key secret
func WithSecret(secret string) Option {
	return func(opts *Options) { opts.Secret = secret }
}

func (s *Options) GetNamespaceNames() []string {
	return s.NamespaceNames
}
