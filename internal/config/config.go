package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/subosito/gotenv"
)

type Config struct {
	Port uint16 `env:"PORT" envDefault:"8080"`
	// NetworkAddress is empty by default, since that implies listening on
	// all interfaces. For development usecases, on windows for example, this
	// is very annoying, as windows will nag you with firewall prompts.
	NetworkAddress string `env:"NETWORK_ADDRESS"`
	CPUProfilePath string `env:"CPU_PROFILE_PATH"`
}

// Load loads the configuration from the environment. If a .env file is
// available, it will be loaded as well. Values found in the environment
// will overwrite whatever is load from the .env file.
func Load() (*Config, error) {
	localEnvVars := os.Environ()
	envVars := make(map[string]string, len(localEnvVars))

	// Add local environment variables to EnvVars map
	for _, keyValuePair := range localEnvVars {
		pair := strings.SplitN(keyValuePair, "=", 2)
		envVars[pair[0]] = pair[1]
	}

	dotEnvPath := ".env"
	if _, err := os.Stat(dotEnvPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking for existence of .env file: %w", err)
		}
	} else {
		envFileContent, err := gotenv.Read(dotEnvPath)
		if err != nil {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
		for key, value := range envFileContent {
			if _, keyExistsInEnvVars := envVars[key]; !keyExistsInEnvVars {
				envVars[key] = value
			}
		}
	}

	var config Config
	if err := env.ParseWithOptions(&config, env.Options{
		Environment: envVars,
		OnSet: func(key string, value any, isDefault bool) {
			if !reflect.ValueOf(value).IsZero() {
				log.Printf("Setting '%s' to '%v' (isDefault: %v)\n", key, value, isDefault)
			}
		},
	}); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}
	return &config, nil
}