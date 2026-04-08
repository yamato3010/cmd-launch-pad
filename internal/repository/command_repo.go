package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/cmd-launch-pad/internal/config"
	"github.com/yourname/cmd-launch-pad/internal/models"
	"gopkg.in/yaml.v3"
)

// commandsFile はcommands.yamlの構造体
type commandsFile struct {
	Categories []models.Category `yaml:"categories"`
	Commands   []models.Command  `yaml:"commands"`
}

// defaultCategories はデフォルトカテゴリ一覧
var defaultCategories = []models.Category{
	{ID: "editor", Name: "エディタ", Icon: "✏️", Color: "#7aa2f7"},
	{ID: "git", Name: "Git", Icon: "🌿", Color: "#9ece6a"},
	{ID: "docker", Name: "Docker", Icon: "🐳", Color: "#2ac3de"},
	{ID: "custom", Name: "カスタム", Icon: "⚡", Color: "#e0af68"},
}

// defaultCommands はデフォルトコマンド一覧
var defaultCommands = []models.Command{
	{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Neovim",
		Command:     "nvim",
		Args:        []string{},
		Description: "高機能テキストエディタ",
		CategoryID:  "editor",
		Icon:        "🖊️",
		CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "lazygit",
		Command:     "lazygit",
		Args:        []string{},
		Description: "TUI git クライアント",
		CategoryID:  "git",
		Icon:        "🌿",
		CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "lazydocker",
		Command:     "lazydocker",
		Args:        []string{},
		Description: "TUI docker 管理ツール",
		CategoryID:  "docker",
		Icon:        "🐳",
		CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	},
}

// CommandRepository はコマンドのCRUD操作を提供する
type CommandRepository struct {
	filePath string
}

// NewCommandRepository は新しいCommandRepositoryを生成する
func NewCommandRepository() (*CommandRepository, error) {
	dir, err := config.EnsureConfigDir()
	if err != nil {
		return nil, err
	}
	return &CommandRepository{
		filePath: filepath.Join(dir, config.CommandsFileName),
	}, nil
}

// load はファイルからデータを読み込む。ファイルが存在しない場合はデフォルト値を返す。
func (r *CommandRepository) load() (*commandsFile, error) {
	data, err := os.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return &commandsFile{
			Categories: defaultCategories,
			Commands:   defaultCommands,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("commands.yamlの読み込みに失敗しました: %w", err)
	}

	var cf commandsFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("commands.yamlのパースに失敗しました: %w", err)
	}
	return &cf, nil
}

// save はデータをファイルに書き込む
func (r *CommandRepository) save(cf *commandsFile) error {
	data, err := yaml.Marshal(cf)
	if err != nil {
		return fmt.Errorf("commands.yamlのシリアライズに失敗しました: %w", err)
	}
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("commands.yamlの書き込みに失敗しました: %w", err)
	}
	return nil
}

// ListCommands は全コマンドを返す
func (r *CommandRepository) ListCommands() ([]models.Command, error) {
	cf, err := r.load()
	if err != nil {
		return nil, err
	}
	return cf.Commands, nil
}

// ListCategories は全カテゴリを返す
func (r *CommandRepository) ListCategories() ([]models.Category, error) {
	cf, err := r.load()
	if err != nil {
		return nil, err
	}
	return cf.Categories, nil
}

// GetCommand はIDでコマンドを取得する
func (r *CommandRepository) GetCommand(id string) (*models.Command, error) {
	cf, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, cmd := range cf.Commands {
		if cmd.ID == id {
			c := cmd
			return &c, nil
		}
	}
	return nil, fmt.Errorf("コマンドが見つかりません: %s", id)
}

// AddCommand は新しいコマンドを追加する
func (r *CommandRepository) AddCommand(cmd *models.Command) error {
	cf, err := r.load()
	if err != nil {
		return err
	}
	if cmd.ID == "" {
		cmd.ID = uuid.New().String()
	}
	now := time.Now()
	cmd.CreatedAt = now
	cmd.UpdatedAt = now
	cf.Commands = append(cf.Commands, *cmd)
	return r.save(cf)
}

// UpdateCommand は既存のコマンドを更新する
func (r *CommandRepository) UpdateCommand(cmd *models.Command) error {
	cf, err := r.load()
	if err != nil {
		return err
	}
	for i, c := range cf.Commands {
		if c.ID == cmd.ID {
			cmd.CreatedAt = c.CreatedAt
			cmd.UpdatedAt = time.Now()
			cf.Commands[i] = *cmd
			return r.save(cf)
		}
	}
	return fmt.Errorf("コマンドが見つかりません: %s", cmd.ID)
}

// DeleteCommand はIDでコマンドを削除する
func (r *CommandRepository) DeleteCommand(id string) error {
	cf, err := r.load()
	if err != nil {
		return err
	}
	newCmds := make([]models.Command, 0, len(cf.Commands))
	found := false
	for _, cmd := range cf.Commands {
		if cmd.ID == id {
			found = true
			continue
		}
		newCmds = append(newCmds, cmd)
	}
	if !found {
		return fmt.Errorf("コマンドが見つかりません: %s", id)
	}
	cf.Commands = newCmds
	return r.save(cf)
}

// InitDefaults はデフォルトデータでcommands.yamlを初期化する（ファイルが存在しない場合のみ）
func (r *CommandRepository) InitDefaults() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return r.save(&commandsFile{
			Categories: defaultCategories,
			Commands:   defaultCommands,
		})
	}
	return nil
}
