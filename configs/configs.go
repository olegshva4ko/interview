package configs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	pathToConfig = "./internal/configs/"
	settings     = "api_key = \"YourApiKeyToken\"\naddr = \":8080\"\n"
)

var (
	key string
)

func init() {
	flag.StringVar(&key, "File with API key",
		"./internal/configs/config.toml", "path to config file")
}

//Config struct with settings
type Config struct {
	APIKey string `toml:"api_key"`
	Addr   string `toml:"addr"`
}

// MakeConfig reads toml file with settings.
// If file does not exist it creates directory with basic config file.
// Reads parameter passed (should be api key string)
func MakeConfig() *Config {
	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) { // create folder with config if not exists
		if err := os.Mkdir("internal", os.ModePerm); err != nil {
			panic(err)
		}
		if err := os.Mkdir(pathToConfig, os.ModePerm); err != nil {
			panic(err)
		}
	}	

	if _, err := os.Stat(pathToConfig + "config.toml"); os.IsNotExist(err) { // create config file if not exists
		if err := ioutil.WriteFile(pathToConfig+"config.toml", []byte(settings), 0777); err != nil { // rwxrwxrwx
			panic(err)
		}
	}
	flag.Parse()

	config := new(Config)
	_, err := toml.DecodeFile(key, config)
	if err != nil {
		panic(err)
	}

	// check for external APIkey
	if len(os.Args) > 1 && os.Args[1] != "" {
		config.APIKey = os.Args[1]
	} else {
		fmt.Printf("API key you are using: %s\nPass your API key as the second parameter\n", config.APIKey)
	}

	return config
}
