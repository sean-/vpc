package db

import "time"

type Config struct {
	Scheme             string        `mapstructure:"scheme"`
	User               string        `mapstructure:"user"`
	Password           string        `mapstructure:"password"`
	Host               string        `mapstructure:"host"`
	Port               uint16        `mapstructure:"port"`
	Database           string        `mapstructure:"database"`
	CAPath             string        `mapstructure:"ca_path"`
	CertPath           string        `mapstructure:"cert_path"`
	KeyPath            string        `mapstructure:"key_path"`
	ConnTimeout        time.Duration `mapstructure:"conn_timeout"`
	InsecureSkipVerify bool          `mapstructure:"insecure_skip_verify"`
}
