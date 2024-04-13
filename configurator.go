package main

import (
	"encoding/json"
	"os"
	"os/user"
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
		"localhost:1488", // TODO: DO NOT FORGOR IT
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
		panic(err)
	}

	//fmt.Println("Writed " + string(data))
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
		panic(err)
	}

	return conf
}

func CloseConfigurator(file string, cfg Configurator) {
	SaveAsJson(workingDir+cfgFile, cfg)
}
