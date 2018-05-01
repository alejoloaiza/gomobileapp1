package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/mobile/asset"
)

type Configuration struct {
	InstaUser      string
	InstaPass      string
	Sentences      []string
	BlacklistUsers []string
	FemaleNames    []string
}

var Localconfig *Configuration

func GetConfig(configpath string) *Configuration {
	file, err := asset.Open(configpath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	Localconfig = &configuration
	return &configuration
}
