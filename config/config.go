package config

import (
	"freeway/common"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const path = "E:/go/src/freeway/resources/"

var configStore = &Configure{}

// Configure 全局配置
type Configure struct {
	Server struct {
		Websocket WebsocketInfo `yaml:"websocket"`
		HTTP      HTTPInfo      `yaml:"http"`
	}
	Database struct {
		Default DatabaseInfo `yaml:"default"`
	}
	Logger struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	}
	inited bool
}

// WebsocketInfo websocket服务器配置
type WebsocketInfo struct {
	Port              string        `yaml:"port"`
	PingPeriodSeconds time.Duration `yaml:"pingPeriodSeconds"`
	ReadBufferSize    int        `yaml:"readBufferSize"`
	WriteBufferSize   int        `yaml:"writeBufferSize"`
	AllowOrigin       string        `yaml:"allowOrigin"`
}

// HTTPInfo http服务器配置
type HTTPInfo struct {
	Port        string `yaml:"port"`
	AllowOrigin string `yaml:"allowOrigin"`
}

// DatabaseInfo 数据库配置
type DatabaseInfo struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbType   string `yaml:"dbtype"`
	DbName   string `yaml:"dbname"`
}

func (c *Configure) initConfig(profile common.Profile) error {
	filename := path + profile.String() + ".yaml"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(bytes, c)
	return nil
}

// Init 根据profile加载对应的配置文件
func Init(profile common.Profile) error {
	return configStore.initConfig(profile)
}

// Get 获取配置信息
func Get() *Configure {
	return configStore
}
