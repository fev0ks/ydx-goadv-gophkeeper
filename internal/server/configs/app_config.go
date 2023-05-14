package configs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	defaultPort           = ":3200"
	defaultSecretKey      = ""
	defaultDBConfig       = ""
	defaultPrivateKeyPath = "cmd/server/privkey.pem"
)

type AppConfig struct {
	ServerPort       string `env:"SERVER_PORT" json:"server_port"`
	TokenKey         string `env:"TOKEN_KEY" json:"token_key"`
	PrivateKey       *rsa.PrivateKey
	PrivateKeyPath   string `env:"CRYPTO_KEY_PATH" json:"crypto_key_path"`
	DBConnection     string `env:"DV_CONNECTION" json:"db_connection"`
	DBMaxConnections int    `env:"DB_MAX_CONNECTIONS" json:"db_max_connections"`
}

func InitAppConfig(configPath string) (*AppConfig, error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}
	setupConfigByFlags(config)
	err = setupRSAKey(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitTestAppConfig(configPath string) (*AppConfig, error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}
	setupConfigByFlags(config)
	err = setupRSAKey(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func setupConfigByFlags(cfg *AppConfig) {
	var serverPortF string
	pflag.StringVarP(&serverPortF, "a", "a", defaultPort, "Port of the proto server")

	var privateKeyPathF string
	pflag.StringVarP(&privateKeyPathF, "f", "f", defaultPrivateKeyPath, "Path of Backup store file")

	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", defaultDBConfig, "Postgres DB DSN")

	var cryptoKeyF string
	pflag.StringVarP(&cryptoKeyF, "crypto-key", "c", "", "Path to private key")

	var dbMaxConnF string
	pflag.StringVarP(&dbMaxConnF, "t", "t", "", "DB Max connections")

	var tokenKeyF string
	pflag.StringVarP(&tokenKeyF, "tk", "k", "", "Token key")

	pflag.Parse()

	if cfg.ServerPort != "" && serverPortF != "" {
		cfg.ServerPort = serverPortF
	}
	if cfg.ServerPort != "" && privateKeyPathF != "" {
		cfg.PrivateKeyPath = privateKeyPathF
	}
	if dbDsnF != "" {
		cfg.DBConnection = dbDsnF
	}
	if dbMaxConnF != "" {
		cfg.DBMaxConnections, _ = strconv.Atoi(dbMaxConnF)
	}
	if tokenKeyF != "" {
		cfg.TokenKey = tokenKeyF
	}
}

func setupRSAKey(config *AppConfig) error {
	if config.PrivateKeyPath != "" {
		key, err := readRsaPrivateKey(config.PrivateKeyPath)
		if err != nil {
			return err
		}
		config.PrivateKey = key
	}
	return nil
}

func readConfig(configFilePath string) (*AppConfig, error) {
	if configFilePath == "" {
		return nil, errors.New("failed to init configuration: file path is not specified")
	}
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configFile by '%s': %v", configFilePath, err)
	}
	var config AppConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config json '%s': %v", string(configBytes), err)
	}
	return &config, nil
}

func readRsaPrivateKey(cryptoKeyPath string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(cryptoKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read publicKey by '%s': %v", cryptoKeyPath, err)
	}
	block, _ := pem.Decode(pemBytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse publicKey: %v", err)
	}
	return key, nil
}
