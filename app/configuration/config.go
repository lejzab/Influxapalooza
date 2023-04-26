package configuration

import (
	"github.com/BurntSushi/toml"
)

// struct for influx config, with fields Host, Port, Username, Password, Database, Timeout
// default value for Host is localhost, Port is 8086, Username is admin, Password is admin, Database is influxdb, Timeout is 60 s
type Influx struct {
	Host      string `json:"host" toml:"host" yaml:"host"`
	Port      int    `json:"port" toml:"port" yaml:"port"`
	Username  string `json:"username" toml:"username" yaml:"username"`
	Password  string `json:"password" toml:"password" yaml:"password"`
	Database  string `json:"database" toml:"database" yaml:"database"`
	Timeout   int    `json:"timeout" toml:"timeout" yaml:"timeout"`
	Protocol  string `json:"protocol" toml:"protocol" yaml:"protocol"`
	Precision string `json:"precision" toml:"precision" yaml:"precision"`
}

type Application struct {
	LogLevel string `json:"log_level" toml:"log_level" yaml:"log_level"`
}

// Global app configuration
type Configuration struct {
	Application Application `json:"application" toml:"application" yaml:"application"`
	Influx      Influx      `json:"influx" toml:"influx" yaml:"influx"`
}

/*
Function creates default configuration for application
If there is config file, it will be read and default values will be overwritten.
If config file cannot be read, error will be returned.
Function returns pointer to Configuration struct
*/
func NewConfiguration(config_file string) (*Configuration, error) {
	Influx := Influx{
		Host:      "localhost",
		Port:      8086,
		Username:  "admin",
		Password:  "admin",
		Database:  "testdb",
		Timeout:   60,
		Protocol:  "http",
		Precision: "s",
	}
	Application := Application{
		LogLevel: "INFO",
	}

	Configuration := Configuration{
		Application: Application,
		Influx:      Influx,
	}
	if config_file != "" {
		_, err := toml.DecodeFile(config_file, &Configuration)
		if err != nil {
			return nil, err
		}
	}
	return &Configuration, nil
}
