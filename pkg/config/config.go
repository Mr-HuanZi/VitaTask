package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var Instances *Config

type Config struct {
	Jwt    JwtConfig    `yaml:"auth"`
	Mysql  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	App    AppConfig    `yaml:"app"`
	Member MemberConfig `yaml:"member"`
}

type JwtConfig struct {
	Key           string `yaml:"key"`
	ExpireSeconds int    `yaml:"expire"`
	Issuer        string `yaml:"issuer"`
}

type MySQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbname"`
	Prefix   string `yaml:"prefix"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Auth string `yaml:"auth"`
}

type AppConfig struct {
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
	Debug bool   `yaml:"debug"`
}

type MemberConfig struct {
	DefaultPass string `yaml:"defaultPass"`
}

func NewConfig() *Config {
	return &Config{
		Jwt: JwtConfig{
			Key:           "",
			ExpireSeconds: 7200,
			Issuer:        "",
		},
		Mysql: MySQLConfig{
			User:     "root",
			Password: "",
			Host:     "localhost",
			Port:     3306,
			DbName:   "",
			Prefix:   "",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			Auth: "",
		},
		App: AppConfig{
			Port: 8081,
			Host: "localhost",
		},
	}
}

func Load(filePath string) error {
	// 初始化配置
	Instances = NewConfig()
	// 加载文件
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		// 读文件失败直接宕机
		return errors.New(fmt.Sprintf("Read file err %v\n", err))
	}

	// 解析Yaml
	err = yaml.Unmarshal(yamlFile, Instances)
	if err != nil {
		// 解析失败直接宕机
		return errors.New(fmt.Sprintf("Yaml parsing error %v\n", err))
	}
	return nil
}

// Get 获取配置实例
func Get() *Config {
	return Instances
}
