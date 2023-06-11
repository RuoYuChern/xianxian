package common

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type gzhConf struct {
	Token string `yaml:"token"`
	Api   string `yaml:"api"`
}

type xxcConf struct {
	Token string `yaml:"token"`
	Api   string `yaml:"api"`
}

type infraConf struct {
	DbDns string `yaml:"db_dns"`
}

type loggerConf struct {
	Level      string `yaml:"level"`
	File       string `yaml:"log_file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Env        string `yaml:"env"`
}

type TaoConf struct {
	Gzh   gzhConf    `yaml:"gzh"`
	Xxc   xxcConf    `yaml:"xxc"`
	Infra infraConf  `yaml:"infra"`
	Log   loggerConf `yaml:"logger"`
}

func (c *TaoConf) LoadTaoConf(path string) {
	ymlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ReadFile failed:%s", err.Error())
		panic(err)
	}
	err = yaml.Unmarshal(ymlFile, c)
	if err != nil {
		log.Printf("Unmarshal failed:%s", err.Error())
		panic(err)
	}
}