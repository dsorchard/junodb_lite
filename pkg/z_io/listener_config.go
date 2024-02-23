package z_io

type (
	ServiceEndpoint struct {
		Addr       string
		Network    string
		SSLEnabled bool
	}

	ListenerConfig struct {
		ServiceEndpoint
		Name string
	}
)
