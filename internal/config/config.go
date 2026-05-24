package config

import (
	"os"
	"strconv"
)

const (
	LISTEN      = "0.0.0.0"
	PORT        = 8080
	TELEMT_HOST = "localhost"
	TELEMT_PORT = 9091
)

type Config struct {
	Listen        string `json:"listen"`
	Port          int64  `json:"port"`
	TelemtHost    string `json:"telemt_host"`
	TelemtPort    int64  `json:"telemt_port"`
	defaultConfig *Config
}

func Default() *Config {
	return &Config{
		Listen:     LISTEN,
		Port:       PORT,
		TelemtHost: TELEMT_HOST,
		TelemtPort: TELEMT_PORT,
	}
}

func (c *Config) Default() *Config {
	if c.defaultConfig == nil {
		c.defaultConfig = Default()
	}

	return c.defaultConfig
}

func (c *Config) Validate() {
	if c.Listen == "" {
		c.Listen = c.Default().Listen
	}

	if c.Port <= 0 {
		c.Port = c.Default().Port
	}

	if c.TelemtHost == "" {
		c.TelemtHost = c.Default().TelemtHost
	}

	if c.TelemtPort <= 0 {
		c.TelemtPort = c.Default().TelemtPort
	}
}

func FromEnv() *Config {
	var (
		port        int64
		telemt_port int64
		err         error
	)

	if port, err = strconv.ParseInt(os.Getenv("STUB_PORT"), 10, 64); err != nil || port <= 0 {
		port = -1
	}

	if telemt_port, err = strconv.ParseInt(os.Getenv("STUB_TELEMT_PORT"), 10, 64); err != nil || telemt_port <= 0 {
		telemt_port = -1
	}

	return &Config{
		Listen:     os.Getenv("STUB_HOST"),
		Port:       port,
		TelemtHost: os.Getenv("STUB_TELEMT_HOST"),
		TelemtPort: telemt_port,
	}
}
