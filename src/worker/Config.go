package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 程序配置
type Config struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

// Singleton
var (
	G_config *Config
)

// Load config
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	// 1. read config file
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2. deserialize JSON
	if json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf

	fmt.Println(conf)
	return
}
