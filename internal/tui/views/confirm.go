package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
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
	command *models.Command
	cursor  int // 0: はい, 1: いいえ
	width   int
	height  int
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

	sb.WriteString(styles.ErrorStyle.Render(i18n.T("confirm.title")))
	sb.WriteString("\n\n")

	sb.WriteString(i18n.T("confirm.question"))
	sb.WriteString(fmt.Sprintf(i18n.T("confirm.name"), styles.AppTitle.Render(m.command.Name)))
	sb.WriteString(fmt.Sprintf(i18n.T("confirm.command"), styles.TabInactive.Render(m.command.Command)))
	if m.command.Description != "" {
		sb.WriteString(fmt.Sprintf(i18n.T("confirm.desc"), m.command.Description))
	}
	sb.WriteString("\n")
	sb.WriteString(styles.ErrorStyle.Render(i18n.T("confirm.irreversible")))
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
		yesBtn = focusedStyle.Render(i18n.T("confirm.yes"))
		noBtn = normalStyle.Render(i18n.T("confirm.no_inactive"))
	} else {
		yesBtn = normalStyle.Render(i18n.T("confirm.yes_inactive"))
		noBtn = focusedStyle.Render(i18n.T("confirm.no"))
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, "    ", noBtn))
	sb.WriteString("\n\n")
	sb.WriteString(styles.TabInactive.Render(i18n.T("confirm.hint")))

	return sb.String()
}
