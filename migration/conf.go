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
	MigrationType   string `yaml:"type,omitempty"`
	Name            string `yaml:"name,omitempty"`
	ClonedName      string `yaml:"clonedName,omitempty"`
	PreCloneScript  string `yaml:"preCloneScript,omitempty"`
	PostCloneScript string `yaml:"postCloneScript,omitempty"`
	Selector        string `yaml:"selector,omitempty"`
	Inplace         bool   `yaml:"inplace,omitempty"`
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
