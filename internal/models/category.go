package models

// Category はコマンドを分類するカテゴリエンティティ
type Category struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon"`
	Color string `yaml:"color"`
}
