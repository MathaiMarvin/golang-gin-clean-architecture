package config

import (
	"fmt" //this package is used for formatted I/O
	"github.com/ilyakaznacheev/cleanenv" 
	// The above package is used for reading configuration files and environment variables.
)

type (
	// Config -. This is the main configuration struct that will encompass all other structs.
	// The struct below represents the complete configuration
	Config struct {

		//THis are embedded structs within the config
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
	}

	// App -.THis struct represents the application section of the configuration
	App struct {

		//The name field below represents the name of the application. It is tagged with the "env-required true" - this indicates that it is a required environment variable.
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		// The field below represents the version of the application.
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.Represents the HTTP section of the configuration
	HTTP struct {

		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.This represents the logging configuration section of the overall configuration
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax        int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		VaTimezone     string `env-required:"true"   yaml:"pg_timezone"    env:"VA_TIMEZONE"`
		DatabaseDriver string `env-required:"true"   yaml:"database_driver"    env:"DATABASE_DRIVER"`
		PostgresUrl    string `env-required:"true"   yaml:"PG_URL"    env:"PG_URL"`
	}
)

// NewConfig returns app config -.Tasked with returning a pointer to a config structure and an error
func NewConfig() (*Config, error) {

	//Initialize a new variable cfg which is of the type config which is a pointer to the config struct
	//the Config{} creates a new instance of the config struct. The curely braces indicate a new struct being created.
	//The ampersand takes the address of the newly created Config struct instance. This creates a pointer to the struct
	// The cfg gets assigned the memory address of the newly created struct allowing for manipulation of the configuration data.
	cfg := &Config{}

	//The cleanenv.ReadConfig function is called to read the configuration from the YAML file. and populate cfg
	err := cleanenv.ReadConfig("./config/config.yml", cfg)

	// if an error occurs during the configuration reading, the function returns nil and error message wrapped with fmt.Errorf()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	//the read env function is called to read the environment variables and override the corresponding fields in the cfg variable.
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type T struct {
	Storage struct {
		File struct {
			Path string `json:"path"`
		} `json:"file"`
	} `json:"storage"`
	Listener struct {
		Tcp struct {
			Address    string `json:"address"`
			TlsDisable bool   `json:"tls_disable"`
		} `json:"tcp"`
	} `json:"listener"`
	Ui              bool   `json:"ui"`
	MaxLeaseTtl     string `json:"max_lease_ttl"`
	DefaultLeaseTtl string `json:"default_lease_ttl"`
	DisableMlock    bool   `json:"disable_mlock"`
}
