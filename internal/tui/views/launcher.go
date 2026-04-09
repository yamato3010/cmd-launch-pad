package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/components"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// LauncherMsg はランチャーからのメッセージ
type LauncherMsg struct {
	Action  LauncherAction
	Command *models.Command
}

// LauncherAction はランチャーのアクション種別
type LauncherAction int

const (
	LauncherActionExec     LauncherAction = iota // コマンド実行
	LauncherActionEdit                           // 編集
	LauncherActionDelete                         // 削除
	LauncherActionNew                            // 新規作成
	LauncherActionSearch                         // 検索
	LauncherActionGit                            // Git操作
	LauncherActionHelp                           // ヘルプ
	LauncherActionCategory                       // カテゴリ管理
)

// KeyMap はランチャー画面のキーバインド
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	New      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Tab      key.Binding
	Search   key.Binding
	Git      key.Binding
	Help     key.Binding
	Quit     key.Binding
	Category key.Binding
}

var DefaultKeyMap = KeyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k")),
	Down:     key.NewBinding(key.WithKeys("down", "j")),
	Left:     key.NewBinding(key.WithKeys("left", "h")),
	Right:    key.NewBinding(key.WithKeys("right", "l")),
	Enter:    key.NewBinding(key.WithKeys("enter")),
	New:      key.NewBinding(key.WithKeys("n")),
	Edit:     key.NewBinding(key.WithKeys("e")),
	Delete:   key.NewBinding(key.WithKeys("d")),
	Tab:      key.NewBinding(key.WithKeys("tab")),
	Search:   key.NewBinding(key.WithKeys("/")),
	Git:      key.NewBinding(key.WithKeys("g")),
	Help:     key.NewBinding(key.WithKeys("?")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c")),
	Category: key.NewBinding(key.WithKeys("c")),
}

// LauncherModel はメインランチャー画面のモデル
type LauncherModel struct {
	commands    []models.Command
	categories  []models.Category
	filtered    []models.Command // フィルタ後のコマンド
	cursor      int
	cols        int
	activeTabID string // アクティブなカテゴリID ("" = 全て)
	width       int
	height      int
	keyMap      KeyMap
}

// NewLauncherModel は新しいLauncherModelを生成する
func NewLauncherModel(commands []models.Command, categories []models.Category, cols int) LauncherModel {
	m := LauncherModel{
		commands:    commands,
		categories:  categories,
		cols:        cols,
		activeTabID: "",
		keyMap:      DefaultKeyMap,
	}
	m.applyFilter()
	return m
}

// SetCommands はコマンド一覧を更新する
func (m *LauncherModel) SetCommands(commands []models.Command) {
	m.commands = commands
	m.applyFilter()
	// カーソルが範囲外になった場合は調整
	total := len(m.filtered) + 1 // +1 for add card
	if m.cursor >= total {
		m.cursor = total - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// SetCategories はカテゴリ一覧を更新する
func (m *LauncherModel) SetCategories(categories []models.Category) {
	m.categories = categories
	// アクティブタブのカテゴリが削除された場合は「全て」に戻す
	found := m.activeTabID == ""
	for _, cat := range categories {
		if cat.ID == m.activeTabID {
			found = true
			break
		}
	}
	if !found {
		m.activeTabID = ""
	}
	m.applyFilter()
}

// applyFilter はアクティブタブに基づいてコマンドをフィルタリングする
func (m *LauncherModel) applyFilter() {
	if m.activeTabID == "" {
		m.filtered = make([]models.Command, len(m.commands))
		copy(m.filtered, m.commands)
		return
	}
	filtered := make([]models.Command, 0)
	for _, cmd := range m.commands {
		if cmd.CategoryID == m.activeTabID {
			filtered = append(filtered, cmd)
		}
	}
	m.filtered = filtered
}

// Init はLauncherModelの初期化コマンドを返す
func (m LauncherModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理してモデルを更新する
func (m LauncherModel) Update(msg tea.Msg) (LauncherModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		total := len(m.filtered) + 1 // +1 for add card
		switch {
		case key.Matches(msg, m.keyMap.Up):
			if m.cursor >= m.cols {
				m.cursor -= m.cols
			}
		case key.Matches(msg, m.keyMap.Down):
			if m.cursor+m.cols < total {
				m.cursor += m.cols
			}
		case key.Matches(msg, m.keyMap.Left):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keyMap.Right):
			if m.cursor < total-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keyMap.Tab):
			m.nextTab()
		case key.Matches(msg, m.keyMap.Enter):
			if m.cursor == len(m.filtered) {
				// 追加カード
				return m, func() tea.Msg {
					return LauncherMsg{Action: LauncherActionNew}
				}
			}
			cmd := m.filtered[m.cursor]
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionExec, Command: &cmd}
			}
		case key.Matches(msg, m.keyMap.New):
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionNew}
			}
		case key.Matches(msg, m.keyMap.Edit):
			if m.cursor < len(m.filtered) {
				cmd := m.filtered[m.cursor]
				return m, func() tea.Msg {
					return LauncherMsg{Action: LauncherActionEdit, Command: &cmd}
				}
			}
		case key.Matches(msg, m.keyMap.Delete):
			if m.cursor < len(m.filtered) {
				cmd := m.filtered[m.cursor]
				return m, func() tea.Msg {
					return LauncherMsg{Action: LauncherActionDelete, Command: &cmd}
				}
			}
		case key.Matches(msg, m.keyMap.Search):
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionSearch}
			}
		case key.Matches(msg, m.keyMap.Git):
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionGit}
			}
		case key.Matches(msg, m.keyMap.Help):
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionHelp}
			}
		case key.Matches(msg, m.keyMap.Category):
			return m, func() tea.Msg {
				return LauncherMsg{Action: LauncherActionCategory}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// 画面幅に応じて列数を動的に計算する
		// カード1枚あたりの幅: Width(14) + Padding(0,1)*2 + Border*2 = 18
		const cardWidth = 18
		if m.width > 0 {
			cols := m.width / cardWidth
			if cols < 1 {
				cols = 1
			}
			m.cols = cols
		}
	}
	return m, nil
}

// nextTab は次のカテゴリタブに切り替える
func (m *LauncherModel) nextTab() {
	if m.activeTabID == "" {
		if len(m.categories) > 0 {
			m.activeTabID = m.categories[0].ID
		}
	} else {
		for i, cat := range m.categories {
			if cat.ID == m.activeTabID {
				if i+1 < len(m.categories) {
					m.activeTabID = m.categories[i+1].ID
				} else {
					m.activeTabID = "" // 全てに戻る
				}
				break
			}
		}
	}
	m.applyFilter()
	m.cursor = 0
}

// View はランチャー画面を描画する
func (m LauncherModel) View() string {
	// ステータスバー（最下部固定）
	bindings := []components.KeyBinding{
		{Key: "↑↓←→/hjkl", Desc: "移動"},
		{Key: "Enter", Desc: "実行"},
		{Key: "n", Desc: "新規"},
		{Key: "e", Desc: "編集"},
		{Key: "d", Desc: "削除"},
		{Key: "Tab", Desc: "タブ切替"},
		{Key: "/", Desc: "検索"},
		{Key: "c", Desc: "カテゴリ"},
		{Key: "g", Desc: "Git"},
		{Key: "?", Desc: "ヘルプ"},
		{Key: "q", Desc: "終了"},
	}
	statusBar := components.RenderStatusBar(bindings, m.width)

	// ヘッダー
	header := styles.AppTitle.Render("⌨  cmd-launch-pad")

	// カテゴリタブ
	tabs := m.renderTabs()

	// グリッド
	grid := components.RenderGrid(m.filtered, m.cursor, m.cols, true)

	// 説明パネル（フォーカス中コマンドの名前・説明を表示）
	var descName, descBody string
	if m.cursor < len(m.filtered) {
		descName = m.filtered[m.cursor].Name
		descBody = m.filtered[m.cursor].Description
	}
	descPanel := components.RenderDescPanel(descName, descBody, m.width)

	// メインコンテンツ（ヘッダー + タブ + グリッド）
	content := header + "\n" + tabs + "\n\n" + grid

	// ステータスバーと説明パネルを画面最下部に配置するため、
	// コンテンツとの間にパディングを挿入する
	if m.height > 0 {
		contentLines := strings.Count(content, "\n") + 1
		// ステータスバーの実際の行数を考慮する（改行されると2行以上になる）
		statusBarLines := strings.Count(statusBar, "\n") + 1
		bottomLines := strings.Count(descPanel, "\n") + 1 + 1 + statusBarLines
		padding := m.height - contentLines - bottomLines
		if padding > 0 {
			content += strings.Repeat("\n", padding)
		}
	}

	return content + "\n" + descPanel + "\n" + statusBar
}

// renderTabs はカテゴリタブを描画する
func (m LauncherModel) renderTabs() string {
	parts := []string{}

	// 「全て」タブ
	allLabel := fmt.Sprintf("📁 全て")
	if m.activeTabID == "" {
		parts = append(parts, styles.TabActive.Render(allLabel))
	} else {
		parts = append(parts, styles.TabInactive.Render(allLabel))
	}

	// カテゴリタブ
	for _, cat := range m.categories {
		label := fmt.Sprintf("%s %s", cat.Icon, cat.Name)
		if m.activeTabID == cat.ID {
			parts = append(parts, styles.TabActive.Render(label))
		} else {
			parts = append(parts, styles.TabInactive.Render(label))
		}
	}

	return strings.Join(parts, " ")
}
