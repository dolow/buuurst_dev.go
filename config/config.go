package config

var Current *Config = nil

type Config struct {
	Eenable      string   `json:"enable"`
	ProjectID    string   `json:"project_id"`
	ServiceKey   string   `json:"service_key"`
	CustomHeader []string `json:"custom_header"`
	IgnorePaths  []string `json:"appignore_paths"`
}

func init() {
	Current = &Config{
		// TODO: parse config
	}
}
