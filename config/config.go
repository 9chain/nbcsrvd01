package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	SDKSrv struct {
		Username string `yaml:"username"`
		ApiKey   string `yaml:"apikey"`
	} `yaml:"sdksrv"`
	SMTP struct {
		Host       string `yaml:"host"`
		ServerAddr string `yaml:"serveraddr"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
	} `yaml:"smtp"`
}

var (
	Cfg Config
)

const (
	cfgFileName = "./config.yaml"
)

const defaultConfig = `
sdksrv:
  username: superuser
  apikey: api key
smtp:
  host: smtp.163.com
  serveraddr: smtp.163.com:25
  user: xxx@xx.com
  password: xxxxxx
`

func LoadConfig() {
	defer func() {
		bs, _ := yaml.Marshal(Cfg)
		fmt.Printf("config:\n%s\n", string(bs))
	}()

	if err := yaml.Unmarshal([]byte(defaultConfig), &Cfg); err != nil {
		panic(err)
	}

	fmt.Printf("default config:\n\t%+v\n", Cfg)

	if _, err := os.Stat(cfgFileName); err != nil {
		return
	}

	source, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		panic(err)
	}

	var newCfg Config
	if err := yaml.Unmarshal(source, &newCfg); err != nil {
		panic(err)
	}

	fmt.Printf("config from %s:\n\t%+v\n", cfgFileName, newCfg)
	compareReset(&Cfg, &newCfg)
}

// TODO change a better way
func compareReset(cfg *Config, newCfg *Config) {
	if len(newCfg.SDKSrv.Username) > 0 {
		cfg.SDKSrv.Username = newCfg.SDKSrv.Username
	}
	if len(newCfg.SDKSrv.ApiKey) > 0 {
		cfg.SDKSrv.ApiKey = newCfg.SDKSrv.ApiKey
	}

	if len(newCfg.SMTP.Host) > 0 {
		cfg.SMTP.Host = newCfg.SMTP.Host
	}
	if len(newCfg.SMTP.ServerAddr) > 0 {
		cfg.SMTP.ServerAddr = newCfg.SMTP.ServerAddr
	}
	if len(newCfg.SMTP.User) > 0 {
		cfg.SMTP.User = newCfg.SMTP.User
	}
	if len(newCfg.SMTP.Password) > 0 {
		cfg.SMTP.Password = newCfg.SMTP.Password
	}
}
