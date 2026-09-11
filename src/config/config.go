package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	TargetMonth int    `yaml:"target_month"`
	TargetYear  int    `yaml:"target_year"`
	Conbine     int    `yaml:"conbine"`
	AttFolder   string `yaml:"att_folder"`
}

var mConfig *Config

func GetConfig() *Config {
	return mConfig
}

func InitConfig() {
	data, err := os.ReadFile("../conf/conf.yaml")

	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	config.AttFolder = filepath.Join("..", config.AttFolder)
	fmt.Printf("%v", config)
	mConfig = &config
	fmt.Printf("month:%d\n year:%d\n conbine:%d", mConfig.TargetMonth, mConfig.TargetYear, mConfig.Conbine)
}
