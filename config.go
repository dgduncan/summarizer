package summarizer

import (
	"io"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

// Config represents the structure of the configuration file
type Config struct {
	YouTube  []YouTubeConfig `yaml:"youtube"`
	Podcasts []PodcastConfig `yaml:"podcast"`
	Logging  Logging         `yaml:"logging"`
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

type Logging struct {
	Level string `yaml:"level"`
}

func (lc Logging) ToLogger(w io.Writer) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(lc.Level)); err != nil {
		level = slog.LevelInfo
	}

	leveler := new(slog.LevelVar)
	leveler.Set(level)
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: leveler,
	}))
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
