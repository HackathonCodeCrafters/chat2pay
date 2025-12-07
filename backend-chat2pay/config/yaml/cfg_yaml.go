package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App        App        `yaml:"app,omitempty" json:"app"`
	DB         DB         `yaml:"db" json:"db"`
	Websocket  Websocket  `yaml:"web_socket"  json:"web_socket"`
	JWT        JWT        `yaml:"jwt" json:"jwt"`
	Logger     Logger     `yaml:"logger" json:"logger"`
	Kolosal    Kolosal    `yaml:"kolosal"json:"kolosal"`
	Gemini     Gemini     `yaml:"gemini" json:"gemini"`
	OpenAI     OpenAI     `yaml:"open_ai" json:"open_ai"`
	Mistral    Mistral    `yaml:"mistral" json:"mistral"`
	RajaOngkir RajaOngkir `yaml:"rajaongkir" json:"rajaongkir"`
}

type App struct {
	Name string `yaml:"name,omitempty" json:"name"`
	Port string `yaml:"port,omitempty" json:"port"`
	//ReadTimeOut  int    `yaml:"read_time_out" json:"read_time_out"`
	//WriteTimeOut int    `yaml:"write_time_out" json:"write_time_out"`
}

type Websocket struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type DB struct {
	Dialect  string `yaml:"dialect" json:"dialect"`
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	DbName   string `yaml:"db_name" json:"db_name"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	SSLMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	//MaxOpen       int    `yaml:"max_open" json:"max_open"`
	//MaxIdle       int    `yaml:"max_idle" json:"max_idle"`
	//TimeOutSecond int    `yaml:"time_out_second" json:"time_out_second"`
	//LifeTimeMs    int    `yaml:"life_time_ms" json:"life_time_ms"`
	//Charset       string `yaml:"charset" json:"charset"`
}

type Kolosal struct {
	URL       string `yaml:"url" json:"url"`
	APIKey    string `yaml:"api_key" json:"api_key"`
	ModelName string `yaml:"model_name" json:"model_name"`
}

type Gemini struct {
	APIKey string `yaml:"api_key" json:"api_key"`
}

type Mistral struct {
	APIKey string `yaml:"api_key" json:"api_key"`
}

type RajaOngkir struct {
	APIKey string `yaml:"api_key" json:"api_key"`
}

type OpenAI struct {
	APIKey string `yaml:"api_key" json:"api_key"`
}

type JWT struct {
	Key           string `yaml:"key" json:"key"`
	ExpiredMinute int    `yaml:"expired_minute" json:"expired_minute"`
}

type Logger struct {
	Enable bool `yaml:"enable" json:"enable"`
}

func NewConfig() (*Config, error) {
	var config *Config

	yfile, err := os.ReadFile("./config/yaml/app.yaml")

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
