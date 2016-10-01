package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	isSmock           bool     `json:"-"`
	AppEnv            string   `json:"app_env"`
	SMSUid            string   `json:"sms_uid"`
	SMSKey            string   `json:"sms_key"`
	SMTPHost          string   `json:"smtp_host"`
	SMTPUser          string   `json:"smtp_user"`
	SMTPPassword      string   `json:"smtp_password"`
	EmailAddressAdmin []string `json:"email_address_admin"`
	EmailAddressDev   []string `json:"email_address_dev"`
	MongoHost         string   `json:"mongo_host"`
	MongoDatabase     string   `json:"mongo_database"`
	MongoUser         string   `json:"mongo_user"`
	MongoPassword     string   `json:"mongo_password"`
}

var conf *Config

func (c *Config) IsSmockServer() bool {
	if c.isSmock {
		return true
	}
	return c.AppEnv != "ONLINE"
}

func Init(path string) *Config {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("load config thxl.conf failed: ", err)
	}
	conf = &Config{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		log.Fatalln("decode config file failed:", string(buf), err)
	}
	return conf
}

func InitWithParams(path string, isSmock bool) *Config {
	conf := Init(path)
	if conf != nil {
		conf.isSmock = isSmock
	}
	return conf
}

func Instance() *Config {
	if conf == nil {
		Init("/root/thxlfzzx_go/server/thxl.conf")
	}
	return conf
}
