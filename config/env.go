package config

var Env *buuurstDevEnv = nil

type buuurstDevEnv struct {
	PathInfo       string `env:"PATH_INFO"`
	RequestMethod  string `env:"REQUEST_METHOD"`
	QueryString    string `env:"QUERY_STRING"`
	HttpXRequestID string `env:"HTTP_X_REQUEST_ID"`
	Body           []byte
}

func init() {
	Env = &buuurstDevEnv{
		// TODO: parse env
	}
}
