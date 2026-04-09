package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// OutputResultMsg はコマンド出力キャプチャ完了メッセージ
type OutputResultMsg struct {
	CommandName string
	Output      string
	Err         error
}

// OutputViewModel はコマンド実行結果ポップアップのモデル
type OutputViewModel struct {
	commandName string
	output      string
	errText     string
	viewport    viewport.Model
	ready       bool
	width       int
	height      int
}

// NewOutputViewModel は新しいOutputViewModelを生成する
func NewOutputViewModel(commandName, output string, execErr error) OutputViewModel {
	errText := ""
	if execErr != nil {
		errText = execErr.Error()
	}
	return OutputViewModel{
		commandName: commandName,
		output:      output,
		errText:     errText,
	}
}

// Init は初期化コマンドを返す
func (m OutputViewModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理する
func (m OutputViewModel) Update(msg tea.Msg) (OutputViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc", "enter"))):
			return m, func() tea.Msg { return BackMsg{} }
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// ポップアップ内のviewport サイズを調整
		vpW := msg.Width - 10
		vpH := msg.Height - 14
		if vpW < 20 {
			vpW = 20
		}
		if vpH < 5 {
			vpH = 5
		}
		if !m.ready {
			m.viewport = viewport.New(vpW, vpH)
			m.viewport.SetContent(m.buildContent())
			m.ready = true
		} else {
			m.viewport.Width = vpW
			m.viewport.Height = vpH
		}
	}

	if m.ready {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// buildContent はviewportに表示するコンテンツを組み立てる
func (m OutputViewModel) buildContent() string {
	var sb strings.Builder
	if m.errText != "" {
		sb.WriteString(styles.ErrorStyle.Render("エラー: " + m.errText))
		sb.WriteString("\n\n")
	}
	if strings.TrimSpace(m.output) == "" {
		sb.WriteString(styles.TabInactive.Render("(出力なし)"))
	} else {
		sb.WriteString(m.output)
	}
	return sb.String()
}

// View は出力ポップアップを描画する
func (m OutputViewModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m OutputViewModel) ModalView() string {
	var sb strings.Builder

	// タイトル
	titleText := "🖥  実行結果"
	if m.commandName != "" {
		titleText = "🖥  実行結果: " + m.commandName
	}
	sb.WriteString(styles.AppTitle.Render(titleText))
	sb.WriteString("\n")

	if m.errText != "" {
		sb.WriteString(styles.ErrorStyle.Render("終了ステータス: エラー"))
	} else {
		sb.WriteString(styles.SuccessStyle.Render("終了ステータス: 正常終了"))
	}
	sb.WriteString("\n\n")

	// コンテンツ
	if m.ready {
		sb.WriteString(m.viewport.View())
		sb.WriteString("\n")
		// スクロールインジケーター
		pct := int(m.viewport.ScrollPercent() * 100)
		sb.WriteString(styles.TabInactive.Render(
			strings.Repeat("─", 30) + " " + strings.Repeat("─", 10),
		))
		sb.WriteString("\n")
		sb.WriteString(styles.TabInactive.Render(
			"↑↓/PgUp/PgDn: スクロール  " + fmt.Sprintf("%d%%", pct),
		))
	} else {
		if strings.TrimSpace(m.output) == "" && m.errText == "" {
			sb.WriteString(styles.TabInactive.Render("(出力なし)"))
		} else {
			// viewport未初期化時は直接表示（最大20行）
			lines := strings.Split(m.output, "\n")
			max := 20
			if len(lines) < max {
				max = len(lines)
			}
			sb.WriteString(strings.Join(lines[:max], "\n"))
			if len(lines) > 20 {
				sb.WriteString(styles.TabInactive.Render("\n... (残り " + fmt.Sprintf("%d", len(lines)-20) + " 行)"))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(styles.TabInactive.Render("q / Esc / Enter で閉じる"))

	return sb.String()
}
