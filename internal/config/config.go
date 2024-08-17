package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port              int    `envconfig:"PORT" default:"8080"`
	R2AccountID       string `envconfig:"R2_ACCOUNT_ID" required:"true"`
	R2AccessKeyID     string `envconfig:"R2_ACCESS_KEY_ID" required:"true"`
	R2AccessKeySecret string `envconfig:"R2_ACCESS_KEY_SECRET" required:"true"`
	R2Bucket          string `envconfig:"R2_BUCKET" required:"true"`
}

func Load() (*Config, error) {
	// Debug: Print all environment variables
	fmt.Println("Environment variables:")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	// Debug: Print loaded configuration
	fmt.Printf("Loaded configuration: %+v\n", cfg)

	return &cfg, nil
}
