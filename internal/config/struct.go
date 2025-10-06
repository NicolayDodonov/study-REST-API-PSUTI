package config

import "time"

type Config struct {
	Log      *Log      `yaml:"log"`
	HTTP     *Http     `yaml:"http"`
	Postgres *Postgres `yaml:"postgres"`
}
type Log struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type Http struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	Timeout      time.Duration `yaml:"timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"dbname"`
}
