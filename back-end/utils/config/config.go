package config

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/devanadindra/signlink-mobile/back-end/utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host        string      `envconfig:"host"`
	Port        int         `envconfig:"port" validate:"number,required"`
	Environment Environment `envconfig:"environment" validate:"oneof=DEVELOPMENT TEST STAGING PRODUCTION"`
	AIBaseUrl   string      `envconfig:"ai_base_url"`
	Version     string      `envconfig:"version" default:"development"`
	Database    Database    `envconfig:"database"`
	Logger      Logger      `envconfig:"logger"`
	Auth        Auth        `envconfig:"auth"`
	GoogleAuth  GoogleAuth
	RateLimiter RateLimiter `envconfig:"rate_limiter"`
	RajaOngkir  RajaOngkir  `envconfig:"raja_ongkir"`
	Midtrans    Midtrans    `envconfig:"midtrans"`
}

type Database struct {
	CustomerUsername string `envconfig:"customer_username"`
	CustomerPassword string `envconfig:"customer_password"`
	AdminUsername    string `envconfig:"admin_username"`
	AdminPassword    string `envconfig:"admin_password"`
	RootUsername     string `envconfig:"root_username"`
	RootPassword     string `envconfig:"root_password"`
	Host             string `envconfig:"host"`
	Port             string `envconfig:"port"`
	Name             string `envconfig:"name"`
}

type Logger struct {
	Level string `envconfig:"level" validate:"oneof=TRACE DEBUG INFO WARN ERROR FATAL PANIC"`
}

type Auth struct {
	JWT   JWT   `envconfig:"jwt" validate:"required"`
	Basic Basic `envconfig:"basic" validate:"required"`
}

type GoogleAuth struct {
	Enabled  bool   `json:"enabled"`
	ClientID string `json:"client_id"`
}

type JWT struct {
	Username  string        `envconfig:"username" validate:"required"`
	Password  string        `envconfig:"password" validate:"required"`
	ExpireIn  time.Duration `envconfig:"expire_in" default:"1000m"`
	SecretKey string        `envconfig:"secret_key" validate:"required"`
}

type Basic struct {
	Username string `envconfig:"username" validate:"required"`
	Password string `envconfig:"password" validate:"required"`
}

type RateLimiter struct {
	Rps    int `envconfig:"rps" default:"10"`
	Bursts int `envconfig:"bursts" default:"5"`
}

type RajaOngkir struct {
	BaseUrl string `envconfig:"base_url"`
	ApiKey  string `envconfig:"api_key"`
}

type Midtrans struct {
	BaseUrlMidtrans string `envconfig:"base_url_midtrans"`
	ServerKey       string `envconfig:"server_key_midtrans"`
}

var config *Config

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Trace(context.Background(), "Failed to load .env : %v", err)
	}

	var c Config
	err = envconfig.Process("backend", &c)
	if err != nil {
		panic("Failed to Process env : " + err.Error())
	}

	config = &c

	return config
}

func GetConfig() *Config {
	if config != nil {
		return config
	}
	return NewConfig()
}

func (m *Midtrans) AuthHeader() string {
	return base64.StdEncoding.EncodeToString([]byte(m.ServerKey + ":"))
}
