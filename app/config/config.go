package config

import (
	//common "github.com/p1cn/tantan-backend-common/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Config ...
type Config struct {
	Database *Database
}

// Database config
type Database struct {
	UserName string
	Password string
	DBName   string
	Host     string
	Port     string
}

var globalConfig *Config

// Get - config getter()
func Get() *Config {
	return globalConfig
}

func parseConfig(rootPath string, fileName string, cfg interface{}) error {

	fp := filepath.Join(rootPath, string(filepath.Separator), fileName)
	err := UnmarshalJSONConfig(fp, cfg)

	if err != nil {
		return err
	}

	return nil
}

// Init init all conf needed
func Init(configPath string) *Config {
	globalConfig = &Config{}
	err := parseConfig(configPath, "db.json", &globalConfig.Database)
	if err != nil {
		log.Fatalf("####read db config failed :%v", err)
	}
	return globalConfig
}

// UnmarshalJSONConfig read conf file into struct
func UnmarshalJSONConfig(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}
