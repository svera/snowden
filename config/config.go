package config

import (
	"errors"
	"io/ioutil"
	"os"

	"fmt"

	"gopkg.in/yaml.v2"
)

// Config holds data needed to run a server instance
type Config struct {
	Debug      bool              `yaml:"debug"`
	SlackToken string            `yaml:"slack_token"`
	Watched    map[string][]Rule `yaml:"watched"`
}

// Rule contains which files and folders to watch and who to notify in case one of them is to be changed
type Rule struct {
	Names       []string `yaml:"names"`
	Subscribers []string `yaml:"subscribers"`
	Exceptions  []string `yaml:"exceptions"`
}

// Load reads configuration from the passed file name
func Load(src string) ([]byte, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load configuration file. Check that %s exists and that it can be read. Exiting...", src)
	}
	return ioutil.ReadAll(f)
}

// Parse unmarshals the data into a YAML and validates it
func Parse(data []byte) (*Config, error) {
	var err error
	c := &Config{}
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

	if len(c.Watched) == 0 {
		return errors.New("There must be at least one watched repository.")
	}

	for repository, rules := range c.Watched {
		if len(rules) == 0 {
			return fmt.Errorf("There are no rules defined.")
		}
		for i, rule := range rules {
			if len(rule.Names) == 0 {
				return fmt.Errorf("Rule %d for repository %s has no watched files.", i, repository)
			}

			if len(rule.Subscribers) == 0 {
				return fmt.Errorf("Rule %d for repository %s has no subscribers.", i, repository)
			}
		}
	}

	return nil
}
