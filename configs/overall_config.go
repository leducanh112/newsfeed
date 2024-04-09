package configs

import (
	"fmt"
	"io"
	"os"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
)

// allInOneConfig combines all service configs
type allInOneConfig struct {
	MySQL                     mysql.Config              `yaml:"my_sql"`
	Redis                     redis.Options             `yaml:"redis"`
	AuthenticateAndPostConfig AuthenticateAndPostConfig `yaml:"authenticate_and_post_config"`
	WebConfig                 WebConfig                 `yaml:"web_config"`
}

// AuthenticateAndPostConfig is config for AAP service
type AuthenticateAndPostConfig struct {
	Port  int           `yaml:"port"`
	MySQL mysql.Config  `yaml:"my_sql"`
	Redis redis.Options `yaml:"redis"`
}

// WebConfig is config for webapp service
type WebConfig struct {
	Port                int `yaml:"port"`
	AuthenticateAndPost struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"authenticate_and_post"`
	Newsfeed struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"newsfeed"`
}

// getAllInOneConfig parse configs from yaml file
func getAllInOneConfig(path string) (*allInOneConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open all in one config (path=%s) error: %s", path, err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read all in one config (path=%s) error: %s", path, err)
	}

	allInOneConf := &allInOneConfig{}
	if err := yaml.Unmarshal(bs, allInOneConf); err != nil {
		return nil, fmt.Errorf("unmarshal all in one config (path=%s) error: %s", path, err)
	}
	return allInOneConf, nil
}

func GetAuthenticateAndPostConfig(path string) (AuthenticateAndPostConfig, error) {
	allInOneConf, err := getAllInOneConfig(path)
	if err != nil {
		return AuthenticateAndPostConfig{}, err
	}
	return allInOneConf.AuthenticateAndPostConfig, nil
}

func GetWebConfig(path string) (WebConfig, error) {
	allInOneConf, err := getAllInOneConfig(path)
	if err != nil {
		return WebConfig{}, err
	}
	return allInOneConf.WebConfig, nil
}
