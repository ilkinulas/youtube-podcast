package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	ListenAddr string
	DbFile     string
	S3         S3
}

type S3 struct {
	Endpoint             string
	Regioin              string
	Bucket               string
	Key                  string
	Secret               string
	PresignedUrlDuration duration
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
