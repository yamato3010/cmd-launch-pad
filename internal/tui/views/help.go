package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// HelpModel はヘルプ画面のモデル
type HelpModel struct {
	width  int
	height int
}

// NewHelpModel は新しいHelpModelを生成する
func NewHelpModel() HelpModel {
	return HelpModel{}
}

// Init はHelpModelの初期化コマンドを返す
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理する
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc", "?"))) {
			return m, func() tea.Msg { return BackMsg{} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View はヘルプ画面を描画する
func (m HelpModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m HelpModel) ModalView() string {
	var sb strings.Builder

	sb.WriteString(styles.AppTitle.Render(i18n.T("help.title")))
	sb.WriteString("\n\n")

	type helpEntry struct {
		key  string
		desc string
	}

	entries := []helpEntry{
		{"↑↓←→ / hjkl", i18n.T("help.move")},
		{"Enter", i18n.T("help.exec")},
		{"n", i18n.T("help.new")},
		{"e", i18n.T("help.edit")},
		{"d", i18n.T("help.delete")},
		{"Tab", i18n.T("help.tab")},
		{"/", i18n.T("help.search")},
		{"c", i18n.T("help.category")},
		{"g", i18n.T("help.git")},
		{"?", i18n.T("help.help")},
		{"q / Ctrl+C", i18n.T("help.quit")},
	}

	for _, e := range entries {
		line := fmt.Sprintf("%s  %s",
			styles.HelpKey.Render(e.key),
			styles.HelpDesc.Render(e.desc),
		)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(styles.TabInactive.Render(i18n.T("help.close")))

	return sb.String()
}

// BackMsg はメイン画面に戻るメッセージ
type BackMsg struct{}
