package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourname/cmd-launch-pad/internal/models"
	"github.com/yourname/cmd-launch-pad/internal/tui/styles"
)

// CategoryAction はカテゴリ管理画面のアクション種別
type CategoryAction int

const (
	CategoryActionNone   CategoryAction = iota
	CategoryActionAdd                   // 新規追加
	CategoryActionEdit                  // 編集
	CategoryActionDelete                // 削除
)

// CategoryDoneMsg はカテゴリ管理画面の完了メッセージ
type CategoryDoneMsg struct {
	Action       CategoryAction
	Category     models.Category
	WithCommands bool // 削除時: コマンドも一緒に削除するか
}

// CategoryViewMode はカテゴリ管理画面のモード
type CategoryViewMode int

const (
	CategoryModeList    CategoryViewMode = iota // 一覧表示
	CategoryModeForm                            // 追加/編集フォーム
	CategoryModeConfirm                         // 削除確認ダイアログ
)

// catFormField はフォームフィールドのインデックス
const (
	catFieldID    = 0
	catFieldName  = 1
	catFieldIcon  = 2
	catFieldColor = 3
	catFieldCount = 4
)

// CategoryViewModel はカテゴリ管理画面のモデル
type CategoryViewModel struct {
	mode       CategoryViewMode
	categories []models.Category
	cmdCounts  map[string]int // カテゴリIDごとのコマンド数
	cursor     int
	isEdit     bool // フォームが編集モードか
	errMsg     string // インライン表示エラー

	// フォーム入力
	inputs    []textinput.Model
	focusIdx  int

	// 削除確認
	deleteTarget  *models.Category
	deleteConfirm int // 0: コマンドも削除, 1: カテゴリのみ
	cmdCountOfDel int // 削除対象カテゴリのコマンド数

	width  int
	height int
}

// NewCategoryViewModel は新しいCategoryViewModelを生成する
func NewCategoryViewModel(categories []models.Category, cmdCounts map[string]int) CategoryViewModel {
	return CategoryViewModel{
		mode:      CategoryModeList,
		categories: categories,
		cmdCounts:  cmdCounts,
	}
}

// SetError はエラーメッセージをセットする
func (m *CategoryViewModel) SetError(msg string) {
	m.errMsg = msg
}

// Init は初期化コマンドを返す
func (m CategoryViewModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理する
func (m CategoryViewModel) Update(msg tea.Msg) (CategoryViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch m.mode {
		case CategoryModeList:
			return m.updateList(msg)
		case CategoryModeForm:
			return m.updateForm(msg)
		case CategoryModeConfirm:
			return m.updateConfirm(msg)
		}
	}

	// フォームモードの場合はtextinputにもメッセージを渡す
	if m.mode == CategoryModeForm && len(m.inputs) > m.focusIdx {
		var cmd tea.Cmd
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m CategoryViewModel) updateList(msg tea.KeyMsg) (CategoryViewModel, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
		return m, func() tea.Msg { return BackMsg{} }
	case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
		if m.cursor < len(m.categories)-1 {
			m.cursor++
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
		// 新規追加フォームを開く
		m.mode = CategoryModeForm
		m.isEdit = false
		m.inputs = makeCatInputs(nil)
		m.focusIdx = 0
		return m, textinput.Blink
	case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
		// 編集フォームを開く
		if len(m.categories) > 0 {
			cat := m.categories[m.cursor]
			m.mode = CategoryModeForm
			m.isEdit = true
			m.inputs = makeCatInputs(&cat)
			m.focusIdx = 0
			return m, textinput.Blink
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
		// 削除確認ダイアログを開く
		if len(m.categories) > 0 {
			cat := m.categories[m.cursor]
			m.deleteTarget = &cat
			m.deleteConfirm = 0
			m.cmdCountOfDel = m.cmdCounts[cat.ID]
			m.mode = CategoryModeConfirm
		}
	}
	return m, nil
}

func (m CategoryViewModel) updateForm(msg tea.KeyMsg) (CategoryViewModel, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		m.mode = CategoryModeList
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
		return m, m.saveCategory()
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "down"))):
		m.focusIdx = (m.focusIdx + 1) % (catFieldCount + 1)
		if m.focusIdx == catFieldCount {
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, nil
		}
		return m, m.updateCatFocus()
	case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))):
		m.focusIdx--
		if m.focusIdx < 0 {
			m.focusIdx = catFieldCount
		}
		if m.focusIdx == catFieldCount {
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, nil
		}
		return m, m.updateCatFocus()
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if m.focusIdx == catFieldCount {
			return m, m.saveCategory()
		}
		if m.focusIdx < catFieldCount-1 {
			m.focusIdx++
			return m, m.updateCatFocus()
		}
		return m, m.saveCategory()
	}
	// アクティブフィールドにメッセージを渡す
	if m.focusIdx < catFieldCount && len(m.inputs) > m.focusIdx {
		var cmd tea.Cmd
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m CategoryViewModel) updateConfirm(msg tea.KeyMsg) (CategoryViewModel, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))):
		m.mode = CategoryModeList
		m.deleteTarget = nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h"))):
		if m.deleteConfirm > 0 {
			m.deleteConfirm--
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l"))):
		if m.deleteConfirm < 1 {
			m.deleteConfirm++
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
		m.deleteConfirm = (m.deleteConfirm + 1) % 2
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		if m.deleteTarget == nil {
			m.mode = CategoryModeList
			return m, nil
		}
		target := *m.deleteTarget
		withCmds := m.deleteConfirm == 0
		m.mode = CategoryModeList
		m.deleteTarget = nil
		return m, func() tea.Msg {
			return CategoryDoneMsg{
				Action:       CategoryActionDelete,
				Category:     target,
				WithCommands: withCmds,
			}
		}
	}
	return m, nil
}

func (m *CategoryViewModel) updateCatFocus() tea.Cmd {
	cmds := make([]tea.Cmd, catFieldCount)
	for i := range m.inputs {
		if i == m.focusIdx {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *CategoryViewModel) saveCategory() tea.Cmd {
	if len(m.inputs) < catFieldCount {
		return nil
	}
	cat := models.Category{
		ID:    m.inputs[catFieldID].Value(),
		Name:  m.inputs[catFieldName].Value(),
		Icon:  m.inputs[catFieldIcon].Value(),
		Color: m.inputs[catFieldColor].Value(),
	}
	action := CategoryActionAdd
	if m.isEdit {
		action = CategoryActionEdit
	}
	m.mode = CategoryModeList
	return func() tea.Msg {
		return CategoryDoneMsg{Action: action, Category: cat}
	}
}

// makeCatInputs はカテゴリフォームのinputsを生成する
func makeCatInputs(cat *models.Category) []textinput.Model {
	placeholders := []string{
		"カテゴリID (例: editor)",
		"カテゴリ名 (例: エディタ)",
		"アイコン (例: ✏️)",
		"カラー (例: #7aa2f7)",
	}
	inputs := make([]textinput.Model, catFieldCount)
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 50
		inputs[i] = ti
	}
	if cat != nil {
		inputs[catFieldID].SetValue(cat.ID)
		inputs[catFieldName].SetValue(cat.Name)
		inputs[catFieldIcon].SetValue(cat.Icon)
		inputs[catFieldColor].SetValue(cat.Color)
	}
	inputs[0].Focus()
	return inputs
}

// View はカテゴリ管理画面を描画する
func (m CategoryViewModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m CategoryViewModel) ModalView() string {
	switch m.mode {
	case CategoryModeList:
		return m.viewList()
	case CategoryModeForm:
		return m.viewForm()
	case CategoryModeConfirm:
		return m.viewConfirm()
	}
	return ""
}

func (m CategoryViewModel) viewList() string {
	var sb strings.Builder
	sb.WriteString(styles.AppTitle.Render("🗂  カテゴリ管理"))
	sb.WriteString("\n\n")

	// インラインエラー表示
	if m.errMsg != "" {
		sb.WriteString(styles.ErrorStyle.Render("⚠  " + m.errMsg))
		sb.WriteString("\n\n")
	}

	if len(m.categories) == 0 {
		sb.WriteString(styles.TabInactive.Render("カテゴリがありません。n で新規追加できます。"))
		sb.WriteString("\n")
	} else {
		// ヘッダー行
		sb.WriteString(styles.InputLabel.Copy().Width(10).Render("ID"))
		sb.WriteString(styles.InputLabel.Copy().Width(14).Render("名前"))
		sb.WriteString(styles.InputLabel.Copy().Width(8).Render("アイコン"))
		sb.WriteString(styles.InputLabel.Copy().Width(12).Render("カラー"))
		sb.WriteString(styles.InputLabel.Render("Cmd"))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 52))
		sb.WriteString("\n")

		for i, cat := range m.categories {
			count := m.cmdCounts[cat.ID]
			line := fmt.Sprintf("%-10s %-12s %-7s %-11s %d",
				cat.ID, cat.Name, cat.Icon, cat.Color, count)
			if i == m.cursor {
				sb.WriteString(styles.CardFocused.Copy().
					Width(52).Height(1).
					Render(line))
			} else {
				sb.WriteString(styles.CardNormal.Copy().
					Width(52).Height(1).
					Render(line))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(styles.TabInactive.Render("↑↓/jk: 移動  n: 追加  e: 編集  d: 削除  q/Esc: 閉じる"))

	return sb.String()
}

func (m CategoryViewModel) viewForm() string {
	var sb strings.Builder

	title := "カテゴリ新規追加"
	if m.isEdit {
		title = "カテゴリ編集"
	}
	sb.WriteString(styles.AppTitle.Render("⌨  cmd-launch-pad - " + title))
	sb.WriteString("\n\n")

	labels := []string{"ID:", "名前:", "アイコン:", "カラー:"}
	for i, label := range labels {
		if i < len(m.inputs) {
			labelStr := styles.InputLabel.Render(label)
			sb.WriteString(fmt.Sprintf("%s  %s\n", labelStr, m.inputs[i].View()))
		}
	}

	sb.WriteString("\n")
	// 保存ボタン
	if m.focusIdx == catFieldCount {
		sb.WriteString(styles.CardFocused.Copy().Width(20).Height(1).Render("[ 保存 (Ctrl+S) ]"))
	} else {
		sb.WriteString(styles.CardNormal.Copy().Width(20).Height(1).Render("[ 保存 (Ctrl+S) ]"))
	}
	sb.WriteString("\n\n")
	sb.WriteString(styles.TabInactive.Render("Tab/Enter: 次へ  Ctrl+S: 保存  Esc: キャンセル"))

	return sb.String()
}

func (m CategoryViewModel) viewConfirm() string {
	if m.deleteTarget == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(styles.AppTitle.Render("⌨  cmd-launch-pad - カテゴリ削除確認"))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("カテゴリ %s「%s %s」を削除します。\n\n",
		styles.ErrorStyle.Render(m.deleteTarget.ID),
		m.deleteTarget.Icon,
		m.deleteTarget.Name,
	))

	if m.cmdCountOfDel > 0 {
		sb.WriteString(fmt.Sprintf("このカテゴリには %s 件のコマンドが属しています。\n",
			styles.ErrorStyle.Render(fmt.Sprintf("%d", m.cmdCountOfDel)),
		))
		sb.WriteString("削除方法を選択してください:\n\n")

		// 選択肢
		opt0 := fmt.Sprintf("コマンドも一緒に削除 (%d件)", m.cmdCountOfDel)
		opt1 := "カテゴリのみ削除（コマンドは残す）"
		if m.deleteConfirm == 0 {
			sb.WriteString(styles.CardFocused.Copy().Width(36).Height(1).Render("▶ " + opt0))
			sb.WriteString("  ")
			sb.WriteString(styles.CardNormal.Copy().Width(36).Height(1).Render("  " + opt1))
		} else {
			sb.WriteString(styles.CardNormal.Copy().Width(36).Height(1).Render("  " + opt0))
			sb.WriteString("  ")
			sb.WriteString(styles.CardFocused.Copy().Width(36).Height(1).Render("▶ " + opt1))
		}
		sb.WriteString("\n\n")
	} else {
		sb.WriteString("このカテゴリにはコマンドが属していません。\n\n")
		sb.WriteString(styles.CardFocused.Copy().Width(20).Height(1).Render("▶ 削除する"))
		sb.WriteString("\n\n")
	}

	sb.WriteString(styles.TabInactive.Render("←→/Tab: 選択切替  Enter: 実行  Esc: キャンセル"))

	return styles.DialogBox.Render(sb.String())
}
