package config

import "github.com/BurntSushi/toml"

type Config struct {
	ListenAddr string
	DbFile     string
	S3         S3
}

type S3 struct {
	Endpoint string
	Regioin  string
	Bucket   string
	Key      string
	Secret   string
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
