package config

type Config struct {
	PluginDirectory string `json:"plugin_directory"`
	ConfigDirectory string `json:"config_directory"`

	DefaultStore string `json:"default_store"`

	GracefulShutdownTimeout int `json:"graceful_shutdown_timeout"`
}
