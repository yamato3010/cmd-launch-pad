package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourname/cmd-launch-pad/internal/tui/styles"
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

	sb.WriteString(styles.AppTitle.Render("❓  キーバインド一覧"))
	sb.WriteString("\n\n")

	type helpEntry struct {
		key  string
		desc string
	}

	entries := []helpEntry{
		{"↑↓←→ / hjkl", "カーソル移動"},
		{"Enter", "コマンド実行"},
		{"n", "新規コマンド登録"},
		{"e", "選択中のコマンド編集"},
		{"d", "選択中のコマンド削除"},
		{"Tab", "カテゴリタブ切り替え"},
		{"/", "検索モード"},
		{"c", "カテゴリ管理"},
		{"g", "Git操作画面"},
		{"?", "ヘルプ表示/非表示"},
		{"q / Ctrl+C", "アプリ終了"},
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
	sb.WriteString(styles.TabInactive.Render("q / Esc / ? で閉じる"))

	return sb.String()
}

// BackMsg はメイン画面に戻るメッセージ
type BackMsg struct{}
