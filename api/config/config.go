package config

import (
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type settings struct {
	BaseUrl string `yaml:"baseUrl"`
	ApiKey  string `yaml:"apiKey"`
}

type config struct {
	Settings settings `yaml:"settings"`
}

func LoadConfigForYaml() (*config, error) {
	f, err := os.Open("/go/api/config/movieConfig.yml")

	if err != nil {
		log.Fatal("loadConfigForYaml os.Open err:", err)
		return nil, err
	}
	defer f.Close()

	var cfg config
	err = yaml.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
