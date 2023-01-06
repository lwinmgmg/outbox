package config

import (
	"os"

	"github.com/lwinmgmg/outbox/helper"

	"gopkg.in/yaml.v2"
)

var (
	configData *Config = nil
)

type Topic struct {
	Name              string   `yaml:"name"`
	DbHost            string   `yaml:"db_host"`
	DbPort            int      `yaml:"db_port"`
	DbUser            string   `yaml:"db_user"`
	DbPassword        string   `yaml:"db_password"`
	DbName            string   `yaml:"db_name"`
	Driver            string   `yaml:"driver"`
	Table             string   `yaml:"table"`
	Order             string   `yaml:"order"`
	Brokers           []string `yaml:"brokers"`
	SecurityProtocol  string   `yaml:"security_protocol"`
	SaslMechanism     string   `yaml:"sasl_mechanism"`
	KafkaUser         string   `yaml:"kafka_user"`
	KafkaPassword     string   `yaml:"kafka_password"`
	ProduceCount      int      `yaml:"produce_count"`
	ProduceIntervalMs int      `yaml:"produce_interval_ms"`
}

func (tpc *Topic) Validate() {
	helper.PanicEmptyString(tpc.Name, "topic name")
	helper.PanicEmptyString(tpc.DbHost, "db_host")
	helper.PanicEmptyString(tpc.DbPort, "db_port")
	helper.PanicEmptyString(tpc.DbName, "db_name")
	helper.PanicEmptyString(tpc.DbUser, "db_user")
	helper.PanicEmptyString(tpc.DbPassword, "db_password")
	tpc.Driver = helper.DefaultString(tpc.Driver, "postgresql")
	tpc.Table = helper.DefaultString(tpc.Table, "outbox")
	tpc.Order = helper.DefaultString(tpc.Order, "created_at, id")
	tpc.SaslMechanism = helper.DefaultString(tpc.SaslMechanism, "PLAINTEXT")
	tpc.SecurityProtocol = helper.DefaultString(tpc.SecurityProtocol, "SASL_PLAINTEXT")
	helper.PanicEmptyString(tpc.KafkaUser, "kafka_user")
	helper.PanicEmptyString(tpc.KafkaPassword, "kafka_password")
	tpc.ProduceCount = helper.DefaultInt(tpc.ProduceCount, 10)             // produce count 10 by default
	tpc.ProduceIntervalMs = helper.DefaultInt(tpc.ProduceIntervalMs, 2000) // each topic produces every 2 seconds
}

type ConfigServer struct {
	Name          string  `yaml:"name"`
	LogDir        string  `yaml:"log_dir"`
	RetryCount    int     `yaml:"retry_count"`
	RespawnTimeMs int     `yaml:"respawn_time_ms"`
	Topics        []Topic `yaml:"topics"`
}

func (confSr *ConfigServer) Validate() {
	confSr.RetryCount = helper.DefaultInt(confSr.RetryCount, 3)
	confSr.RespawnTimeMs = helper.DefaultInt(confSr.RespawnTimeMs, 5000) //5 seconds by default
	for i := 0; i < len(confSr.Topics); i++ {
		confSr.Topics[i].Validate()
	}
}

type Config struct {
	Version string       `yaml:"version"`
	Server  ConfigServer `yaml:"server"`
}

func (conf *Config) Validate() {
	conf.Server.Validate()
}

func GetConfig(filename string) (*Config, error) {
	if configData == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		configData = &Config{}
		yaml.Unmarshal(data, configData)
	}
	configData.Validate()
	return configData, nil
}
