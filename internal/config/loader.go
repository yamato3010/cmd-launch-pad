package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	AppDirName       = "cmd-launch-pad"
	ConfigFileName   = "config.yaml"
	CommandsFileName = "commands.yaml"
)

// ConfigDir はアプリの設定ディレクトリパスを返す (~/.config/cmd-launch-pad/)
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", AppDirName), nil
}

// EnsureConfigDir は設定ディレクトリが存在しない場合は作成する
func EnsureConfigDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	return dir, nil
}

// LoadAppConfig は設定ファイルを読み込む。存在しない場合はデフォルト設定を返す。
func LoadAppConfig() (*AppConfig, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ConfigFileName)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultAppConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &cfg, nil
}

// SaveAppConfig は設定をファイルに書き込む
func SaveAppConfig(cfg *AppConfig) error {
	dir, err := EnsureConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, ConfigFileName)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config file: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
