package summarizer

import (
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

// Config represents the structure of the configuration file
type Config struct {
	YouTube []YouTubeConfig `yaml:"youtube"`
	Podcast []PodcastConfig `yaml:"podcast"`
}

// YouTubeConfig represents the configuration for YouTube
type YouTubeConfig struct {
	Name   string `yaml:"name"`
	RssURL string `yaml:"rss_url"`
	Prompt string `yaml:"prompt"`
}

// PodcastConfig represents the configuration for Podcast
type PodcastConfig struct {
	Name   string `yaml:"name"`
	RssURL string `yaml:"rss_url"`
	Prompt string `yaml:"prompt"`
}

// LoadConfig reads the YAML configuration file and unmarshals it into a Config struct
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
