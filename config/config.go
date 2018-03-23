package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"bytes"
)

type Config struct {
	SDKSrv struct {
		Username string `toml:"username,omitempty"`
		ApiKey   string `toml:"apikey,omitempty"`
	} `toml:"sdksrv,omitempty"`
	SMTP struct {
		Host       string `toml:"host,omitempty"`
		ServerAddr string `toml:"serveraddr,omitempty"`
		User       string `toml:"user,omitempty"`
		Password   string `toml:"password,omitempty"`
		Salt       string `toml:"salt,omitempty"`
	} `toml:"smtp,omitempty"`
}

var (
	Cfg Config
)

const (
	cfgFileName = "./config.toml"
)

const defaultConfig = `
[sdksrv]
username = "superuser"
apikey = "api key"
[smtp]
host = "smtp.163.com"
serveraddr = "smtp.163.com:25"
user = "xxx@xx.com"
password = "xxxxxx"
salt = "salt"
`

func printCfg(flag string, cfg *Config) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(cfg); err != nil {
		panic(err)
	}

	fmt.Printf("==========%s================\n%s\n", flag, buf.String())
}

func LoadConfig() {
	defer printCfg("final", &Cfg)

	if _, err := toml.Decode(defaultConfig, &Cfg); err != nil {
		panic(err)
	}

	if _, err := os.Stat(cfgFileName); err != nil {
		return
	}

	var newCfg Config
	if _, err := toml.DecodeFile(cfgFileName, &newCfg); err != nil {
		panic(err)
	}

	printCfg(cfgFileName, &newCfg)
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
	if len(newCfg.SMTP.Salt) > 0 {
		cfg.SMTP.Salt = newCfg.SMTP.Salt
	}
}
