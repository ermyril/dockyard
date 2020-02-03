package config

import (
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host      string
	Database  string
	User      string
	Password  string
	Directory string
}


func GetConfig() Config {

	configContents := getConfigFileContents("dockyard.yml")
	config := makeConfigObject(configContents)

	return config
}

func getConfigFileContents(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	configContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return configContents
}

func makeConfigObject(configFileContents []byte) Config {
	config := Config{}
	err := yaml.Unmarshal([]byte(configFileContents), &config)
	if err != nil {
		fmt.Println("error: %v", err)
	}

	return config
}
