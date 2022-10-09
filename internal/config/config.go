package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type CurrencyInfo struct {
	Name    string `yaml:"name"`
	Display string `yaml:"display"`
}

type Config struct {
	Token                    string         `yaml:"token"`
	CurrencyURLCurrent       string         `yaml:"currency_url_current"`
	CurrencyURLPast          string         `yaml:"currency_url_past"`
	CurrencyRateEarliestDate time.Time      `yaml:"currency_rate_earliest_date"`
	CurrencyMain             string         `yaml:"currency_main"`
	CurrencyRateLoadInterval int            `yaml:"currency_rate_load_interval"`
	Currencies               []CurrencyInfo `yaml:"currencies"`
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) GetConfig() Config {
	return s.config
}
