package configinternal

// Common - common settings for microservice (server options, socket port and etc)
type Common struct {
	HttpServer   `yaml:"http_server"`
	SocketServer `yaml:"socket_server"`
	WebSocket    `yaml:"web_socket"`
}

type WebSocket struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	Url     string `yaml:"url"`
}

type HttpServer struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type SocketServer struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}
