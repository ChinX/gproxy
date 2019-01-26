package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"net"
)

var (
	defaultSocks = "127.0.0.1:9953"
)

type config struct {
	Server   *Server   `yaml:"server"`
	Listener *Listener `yaml:"listener"`
}

type Server struct {
	Addr     string `yaml:"address"`
	Password string `yaml:"password"`
	Method   string `yaml:"method"`
}

type Listener struct {
	Socks    string `yaml:"socks"`
	HTTP     string `yaml:"http"`
}

func LoadConfig(configFile string) (*config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file \"%s\" faild, error: %s", configFile, err)
	}
	conf := &config{}
	if err = yaml.Unmarshal(data, conf); err != nil {
		return nil, fmt.Errorf("unmarshal yaml file \"%s\" faild, error: %s", configFile, err)
	}

	if conf.Server == nil {
		return nil, fmt.Errorf("server configuration is empty")
	} else {
		conf.Server.Addr = formatAddress(conf.Server.Addr)
		if conf.Server.Password == "" {
			return nil, fmt.Errorf("must specify server password")
		}

		if conf.Server.Method == "" {
			conf.Server.Method = "aes-256-cfb"
		}
	}

	if conf.Listener == nil {
		conf.Listener = &Listener{Socks: defaultSocks}
	}

	if  conf.Listener.Socks == "" {
		conf.Listener.Socks = defaultSocks
	}
	return conf, nil

}

func formatAddress(hostAddr string) string {
	host, _, err := net.SplitHostPort(hostAddr)
	if err != nil {
		host = net.JoinHostPort(hostAddr, "80")
	} else {
		host = hostAddr
	}
	return host
}
