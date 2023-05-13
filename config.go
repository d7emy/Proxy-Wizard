package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type config struct {
	ThreadsCount int `yaml:"threadsCount"`
	Proxy        struct {
		Useproxy bool   `yaml:"useproxy"`
		ForceUse bool   `yaml:"forceuse"`
		URL      string `yaml:"url"`
	} `yaml:"proxy"`
	Urls []URLDATA `yaml:"urls"`
}

type URLDATA struct {
	URL        string `yaml:"url"`
	AutoRegex  bool   `yaml:"autoRegex,omitempty"`
	ToArray    string `yaml:"to_array,omitempty"`
	IPJSON     string `yaml:"ip_json,omitempty"`
	PortJSON   string `yaml:"port_json,omitempty"`
	IPPortJSON string `yaml:"ip_port_json,omitempty"`
	//IPRg       string `yaml:"ip_rg,omitempty"`
	//PortRg     string `yaml:"port_rg,omitempty"`
	IPPortRg string `yaml:"ip_port_rg,omitempty"`
}

var Config = config{}

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(data), &Config)
	if err != nil {
		panic(err)
	}
}
