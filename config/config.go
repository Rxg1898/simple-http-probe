package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpListenAddr         string `yaml:"http_listen_addr"`
	HttpProbeTimeoutSecond int    `yaml:"http_probe_timeout_second"`
}

// 定义一个全局超时时间
var GlobalTwsec int

// 加载配置文件
func Load(in []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(in, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// 读取配置文件，解析
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(content)
	if err != nil {
		log.Printf("加载配置文件错误: %v", err)
		return nil, err
	}

	// 超时配置是否配置
	if cfg.HttpProbeTimeoutSecond == 0 {
		GlobalTwsec = 5
	} else {
		GlobalTwsec = cfg.HttpProbeTimeoutSecond
	}
	return cfg, nil
}
