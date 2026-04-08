package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/components"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// SearchDoneMsg は検索画面の完了メッセージ
type SearchDoneMsg struct {
	Selected *models.Command // nil の場合はキャンセル
	Action   LauncherAction
}

// SearchModel は検索画面のモデル
type SearchModel struct {
	input      textinput.Model
	commands   []models.Command // 全コマンド
	filtered   []models.Command // フィルタ後
	categories []models.Category
	cursor     int
	width      int
	height     int
}

// NewSearchModel は新しいSearchModelを生成する
func NewSearchModel(commands []models.Command, categories []models.Category) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "コマンド名・説明・カテゴリで検索..."
	ti.CharLimit = 100
	ti.Focus()

	m := SearchModel{
		input:      ti,
		commands:   commands,
		categories: categories,
	}
	m.applyFilter()
	return m
}

// applyFilter はクエリに基づいてコマンドをフィルタリングする
func (m *SearchModel) applyFilter() {
	query := strings.ToLower(m.input.Value())
	if query == "" {
		m.filtered = make([]models.Command, len(m.commands))
		copy(m.filtered, m.commands)
		return
	}
	filtered := make([]models.Command, 0)
	for _, cmd := range m.commands {
		if strings.Contains(strings.ToLower(cmd.Name), query) ||
			strings.Contains(strings.ToLower(cmd.Description), query) ||
			strings.Contains(strings.ToLower(cmd.CategoryID), query) ||
			strings.Contains(strings.ToLower(cmd.Command), query) {
			filtered = append(filtered, cmd)
		}
	}
	m.filtered = filtered
}

// Init はSearchModelの初期化コマンドを返す
func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update はキー入力を処理する
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return SearchDoneMsg{Selected: nil} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if len(m.filtered) > 0 {
				cmd := m.filtered[m.cursor]
				return m, func() tea.Msg {
					return SearchDoneMsg{Selected: &cmd, Action: LauncherActionExec}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.applyFilter()
	// カーソル範囲を調整
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	return m, cmd
}

// View は検索画面を描画する
func (m SearchModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m SearchModel) ModalView() string {
	var sb strings.Builder

	sb.WriteString(styles.AppTitle.Render("🔍  検索"))
	sb.WriteString("\n\n")

	sb.WriteString(m.input.View())
	sb.WriteString("\n\n")

	if len(m.filtered) == 0 {
		sb.WriteString(styles.TabInactive.Render("該当するコマンドが見つかりません"))
	} else {
		sb.WriteString(fmt.Sprintf("%s\n\n", styles.TabInactive.Render(
			fmt.Sprintf("%d 件ヒット", len(m.filtered)),
		)))
		grid := components.RenderGrid(m.filtered, m.cursor, 3, false)
		sb.WriteString(grid)
	}

	sb.WriteString("\n")
	sb.WriteString(styles.TabInactive.Render("↑↓/jk: 移動  Enter: 実行  Esc: 閉じる"))

	return sb.String()
}
