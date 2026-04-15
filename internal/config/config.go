package config

// AppConfig はアプリ全体の設定 (~/.config/cmd-launch-pad/config.yaml)
type AppConfig struct {
	Version  string    `yaml:"version"`
	Theme    string    `yaml:"theme"`    // dark / light
	Columns  int       `yaml:"columns"`  // グリッドの列数
	Language string    `yaml:"language"` // 言語設定: "en" / "ja" (空の場合は自動検出)
	Git      GitConfig `yaml:"git"`
}

// GitConfig はGit連携に関する設定
type GitConfig struct {
	Remote   string `yaml:"remote"`    // リモートリポジトリURL
	AutoPush bool   `yaml:"auto_push"` // コマンド変更時に自動プッシュ
	Branch   string `yaml:"branch"`
}

// DefaultAppConfig はデフォルト設定を返す
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Version: "1",
		Theme:   "dark",
		Columns: 4,
		Git: GitConfig{
			Remote:   "",
			AutoPush: false,
			Branch:   "main",
		},
	}
}
