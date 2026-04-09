package models

import "time"

// Command はコマンドランチャーに登録されるコマンドエンティティ
type Command struct {
	ID            string    `yaml:"id"`
	Name          string    `yaml:"name"`
	Command       string    `yaml:"command"`
	Args          []string  `yaml:"args"`
	Description   string    `yaml:"description"`
	CategoryID    string    `yaml:"category_id"`
	Icon          string    `yaml:"icon"`
	CaptureOutput bool      `yaml:"capture_output"` // trueの場合、実行結果をポップアップ表示
	CreatedAt     time.Time `yaml:"created_at"`
	UpdatedAt     time.Time `yaml:"updated_at"`
}
