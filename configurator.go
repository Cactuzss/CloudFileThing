package main

import (
	"encoding/json"
	"os"
	"os/user"
	"strings"
)

type Configurator struct {
	HostAddress string
	Username    string
}

func NewConfigurator() Configurator {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return Configurator{
		"localhost:1109",
		user.Username,
	}
}

func SaveAsJson(path string, config Configurator) {
	data, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			delimPos := strings.LastIndex(path, "/")
			woName := path[:delimPos+1]

			if _, err := os.Stat(woName); os.IsNotExist(err) {
				os.Mkdir(woName, 0777)
			}

			os.Create(path)
		}
	}
}

func LoadFromJson(path string) Configurator {
	var conf = NewConfigurator()

	data, err := os.ReadFile(path)
	if err != nil {
		os.Create(path)
		CloseConfigurator(path, conf)
		return conf
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		os.Create(path)
		CloseConfigurator(path, conf)
		return conf
	}

	return conf
}

func CloseConfigurator(file string, cfg Configurator) {
	SaveAsJson(workingDir+cfgFile, cfg)
}
