package z_conn_mgr

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
