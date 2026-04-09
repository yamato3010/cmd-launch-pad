package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// ConfirmDeleteDoneMsg は削除確認モーダルの完了メッセージ
type ConfirmDeleteDoneMsg struct {
	Confirmed bool
	Command   *models.Command
}

// ConfirmDeleteModel は削除確認モーダルのモデル
type ConfirmDeleteModel struct {
	command  *models.Command
	cursor   int // 0: はい, 1: いいえ
	width    int
	height   int
}

// NewConfirmDeleteModel は新しい削除確認モーダルを生成する
func NewConfirmDeleteModel(cmd *models.Command) ConfirmDeleteModel {
	return ConfirmDeleteModel{
		command: cmd,
		cursor:  1, // デフォルトは「いいえ」(誤削除防止)
	}
}

// Init は初期化コマンドを返す
func (m ConfirmDeleteModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理する
func (m ConfirmDeleteModel) Update(msg tea.Msg) (ConfirmDeleteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "n"))):
			// キャンセル
			cmd := m.command
			return m, func() tea.Msg {
				return ConfirmDeleteDoneMsg{Confirmed: false, Command: cmd}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h", "tab"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l"))):
			if m.cursor < 1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
			// yキーで即座に確認
			cmd := m.command
			return m, func() tea.Msg {
				return ConfirmDeleteDoneMsg{Confirmed: true, Command: cmd}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			confirmed := m.cursor == 0
			cmd := m.command
			return m, func() tea.Msg {
				return ConfirmDeleteDoneMsg{Confirmed: confirmed, Command: cmd}
			}
		}
	}
	return m, nil
}

// View はモーダル表示用コンテンツを返す
func (m ConfirmDeleteModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル内コンテンツを返す（ModalBoxのボーダーなし）
func (m ConfirmDeleteModel) ModalView() string {
	if m.command == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(styles.ErrorStyle.Render("🗑  コマンドの削除"))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("以下のコマンドを削除してもよいですか？\n\n"))
	sb.WriteString(fmt.Sprintf("  名前: %s\n", styles.AppTitle.Render(m.command.Name)))
	sb.WriteString(fmt.Sprintf("  コマンド: %s\n", styles.TabInactive.Render(m.command.Command)))
	if m.command.Description != "" {
		sb.WriteString(fmt.Sprintf("  説明: %s\n", m.command.Description))
	}
	sb.WriteString("\n")
	sb.WriteString(styles.ErrorStyle.Render("この操作は元に戻せません。"))
	sb.WriteString("\n\n")

	// ボタン（lipgloss.JoinHorizontal で横並び）
	focusedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorRed).
		Padding(0, 2)
	normalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(0, 2)

	var yesBtn, noBtn string
	if m.cursor == 0 {
		yesBtn = focusedStyle.Render("▶ はい")
		noBtn  = normalStyle.Render("  いいえ")
	} else {
		yesBtn = normalStyle.Render("  はい")
		noBtn  = focusedStyle.Render("▶ いいえ")
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, "    ", noBtn))
	sb.WriteString("\n\n")
	sb.WriteString(styles.TabInactive.Render("←→/Tab: 選択切替  Enter: 実行  y: 削除  Esc/n: キャンセル"))

	return sb.String()
}
