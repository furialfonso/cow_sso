package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/olebedev/config"
	configYML "github.com/olebedev/config"
)

var cfg *configYML.Config

func init() {
	Start()
}

func Start() {
	scope := os.Getenv("SCOPE")
	if strings.EqualFold(scope, "") {
		scope = "local"
	}
	fmt.Println("SCOPE:", scope)
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		panic("CONFIG_DIR environment variable is not set")
	}
	fileConfig := fmt.Sprintf("%v/properties.yml", configDir)
	b, err := os.ReadFile(fileConfig)
	if err != nil {
		panic(err)
	}

	cfg, err = configYML.ParseYamlBytes(b)
	if err != nil {
		return
	}
}

func Get() *config.Config {
	return cfg
}
