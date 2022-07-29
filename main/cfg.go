package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

const serviceName = "projecttemplate"

var BuildVersion = "default"

//Config configure intent
var Config struct {
	Port string
	Log  struct {
		Level string
	}
	MySql map[string]struct {
		IP       string
		Port     string
		User     string
		PassWord string
		DBname   string
	}
	MongoDB map[string]struct {
		IP             string
		Port           string
		User           string
		Password       string
		DBName         string
		MongoPoolLimit int
	}
}

func loadcfg() {
	showVersion := flag.Bool("v", false, "show version")
	configFilePath := flag.String("c", "", "configure file path")
	configFileTest := flag.Bool("t", false, "configure file test")
	flag.Parse()
	if *showVersion {
		fmt.Println(BuildVersion)
		os.Exit(0)
	}
	Info("cfg path :", *configFilePath)
	if *configFilePath == "" {
		//*configFilePath = "./output/projecttemplate.toml"
		*configFilePath = "./output/test.toml"
	}

	var err error
	if _, err = toml.DecodeFile(*configFilePath, &Config); err != nil {
		Fatal("toml fail to parse file :", err)
		os.Exit(-1)
	}

	Infof("%+v", Config)

	err = LoadVersionV2(serviceName, *configFilePath, BuildVersion)
	if err != nil {
		Error(err)
	}

	if *configFileTest {
		fmt.Printf("configuration file %s test is successful\n", *configFilePath)
		os.Exit(0)
	}
}

func functionInit() {
	mysqlInit()
}
