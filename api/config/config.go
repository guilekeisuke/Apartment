package config

import (
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type config struct {
	Settings      settings      `yaml:"settings"`
	Gmail         gmail         `yaml:"gmail"`
	MovieUpcoming movieUpcoming `yaml:"movieUpcoming"`
}

type settings struct {
	BaseUrl string `yaml:"baseUrl"`
	ApiKey  string `yaml:"apiKey"`
}

type gmail struct {
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Password string `yaml:"password"`
	Smtp     string `yaml:"smtp"`
	Port     string `yaml:"port"`
}

type movieUpcoming struct {
	FromName string `yaml:"fromName"`
	Subject  string `yaml:"subject"`
}

func LoadConfigForYaml(path string) (*config, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("loadConfigForYaml os.Open err:", err)
		return nil, err
	}
	defer f.Close()

	var cfg config
	err = yaml.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
