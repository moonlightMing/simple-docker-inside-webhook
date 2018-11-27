package conf

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	BindHost           string      `json:"bind_host"`
	DockerHost         string      `json:"docker_host"`
	AuthKey            string      `json:"auth_key"`
	DockerApiVersion   string      `json:"docker_api_version"`
	DockerRegistryAuth Auth        `json:"docker_registry_auth"`
	CronEvents         []CronEvent `json:"cron"`
	Email              Email       `json:"email"`
}

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type CronEvent struct {
	Event string `json:"event"`
	Spec  string `json:"spec"`
}

type Email struct {
	Open      bool   `json:"open"`
	SmtpHost  string `json:"smtp_host"`
	SmtpPort  string `json:"smtp_port"`
	UserEmail string `json:"user_email"`
	Password  string `json:"password"`
	SendTo    string `json:"send_to"`
}

func NewConfig() Config {
	once.Do(func() {
		cfgbuf, err := ioutil.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(cfgbuf, &conf)
		if err != nil {
			panic(err)
		}
	})
	return conf
}

func (c *Config) IsAuthentication(key string) bool {
	if c.AuthKey == key {
		return true
	}
	return false
}
