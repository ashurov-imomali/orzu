package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configs struct {
	Srv      Server     `json:"srv" yaml:"srv"`
	Orzu     OrzuParams `json:"orzu" yaml:"orzu"`
	Otp      OtpParams  `json:"otp_params" yaml:"otp_params"`
	Redis    RDb        `json:"redis" yaml:"redis"`
	Postgres Postgres   `json:"postgres" yaml:"postgres"`
}

type Postgres struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DbName   string `json:"db_name" yaml:"db_name"`
}

func New() (*Configs, error) {
	file, err := os.ReadFile("./config/configs.yaml")
	if err != nil {
		return nil, err
	}
	c := &Configs{}
	return c, yaml.Unmarshal(file, &c)
}

type Server struct {
	Host         string `json:"host" yaml:"host"`
	Port         string `json:"port" yaml:"port"`
	ReadTimeout  int    `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" yaml:"write_timeout"`
	Token        string `json:"token" yaml:"token"`
}

type OrzuParams struct {
	Url       string `json:"url" yaml:"url"`
	ServiceId string `json:"service_id" yaml:"service_id"`
	Token     string `json:"token" yaml:"token"`
	TToken    string `json:"t_token" yaml:"t_token"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}

type RDb struct {
	Url      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database int    `json:"db" yaml:"db"`
}

type OtpParams struct {
	Url          string `json:"url" yaml:"url"`
	LifeTime     int64  `json:"life_time" yaml:"life_time"`
	ConfirmLimit int64  `json:"confirm_limit" yaml:"confirm_limit"`
}
