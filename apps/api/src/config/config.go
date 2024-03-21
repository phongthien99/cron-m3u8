package config

import (
	"example/src/config/types"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sigmaott/gest/package/technique/configuration"
)

var config *types.EnvironmentVariable

func init() {
	configPath := flag.String("c", "./src/config/default.yaml", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	var err error
	config, err = configuration.LoadConfigYaml[types.EnvironmentVariable](*configPath)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfiguration() *types.EnvironmentVariable {
	return config
}
