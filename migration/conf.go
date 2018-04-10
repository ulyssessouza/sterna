package migration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Description string      `yaml:"description,omitempty"`
	Migrations  []Migration `yaml:"migrations,omitempty"`
}

type Migration struct {
	Name            string `yaml:"name,omitempty"`
	MigrationType   string `yaml:"type,omitempty"`
	ServiceName     string `yaml:"service,omitempty"`
	PreCloneScript  string `yaml:"preCloneScript,omitempty"`
	PostCloneScript string `yaml:"postCloneScript,omitempty"`
}

func (c *Config) Load(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %#v", err)
	}
}
