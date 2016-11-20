package config

import (
	"errors"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds data needed to run a server instance
type Config struct {
	SlackToken string `yaml:"slack_token"`
}

// Load reads configuration from config.yml and parses it
func Load(src io.Reader) (*Config, error) {
	c := &Config{}
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return c, err
	}

	if err = yaml.Unmarshal(data, c); err != nil {
		return c, err
	}
	err = c.validate()
	return c, err
}

func (c *Config) validate() error {
	if c.SlackToken == "" {
		return errors.New("Invalid token.")
	}
	return nil
}
