package config
import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Site struct {
		Title  string `yaml:"title"`
		Author string `yaml:"author"`
		Bio    string `yaml:"bio"`
	} `yaml:"site"`

	Paths struct {
		Content   string `yaml:"content"`
		Templates string `yaml:"templates"`
		Output    string `yaml:"output"`
	} `yaml:"paths"`
}

func LoadConf() (*Config, error) {
	data, err := os.ReadFile("koyo.config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}
