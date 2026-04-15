package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
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
	isEdit     bool   // フォームが編集モードか
	errMsg     string // インライン表示エラー

	// フォーム入力
	inputs   []textinput.Model
	focusIdx int

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
		mode:       CategoryModeList,
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
		i18n.T("cat.form.placeholder.id"),
		i18n.T("cat.form.placeholder.name"),
		i18n.T("cat.form.placeholder.icon"),
		i18n.T("cat.form.placeholder.color"),
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
	sb.WriteString(styles.AppTitle.Render(i18n.T("cat.title.list")))
	sb.WriteString("\n\n")

	// インラインエラー表示
	if m.errMsg != "" {
		sb.WriteString(styles.ErrorStyle.Render("⚠  " + m.errMsg))
		sb.WriteString("\n\n")
	}

	if len(m.categories) == 0 {
		sb.WriteString(styles.TabInactive.Render(i18n.T("cat.empty")))
		sb.WriteString("\n")
	} else {
		// ヘッダー行
		sb.WriteString(styles.InputLabel.Copy().Width(10).Render(i18n.T("cat.header.id")))
		sb.WriteString(styles.InputLabel.Copy().Width(14).Render(i18n.T("cat.header.name")))
		sb.WriteString(styles.InputLabel.Copy().Width(8).Render(i18n.T("cat.header.icon")))
		sb.WriteString(styles.InputLabel.Copy().Width(12).Render(i18n.T("cat.header.color")))
		sb.WriteString(styles.InputLabel.Render(i18n.T("cat.header.cmd")))
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
	sb.WriteString(styles.TabInactive.Render(i18n.T("cat.list.hint")))

	return sb.String()
}

func (m CategoryViewModel) viewForm() string {
	var sb strings.Builder

	title := i18n.T("cat.form.title.add")
	if m.isEdit {
		title = i18n.T("cat.form.title.edit")
	}
	sb.WriteString(styles.AppTitle.Render("⌨  cmd-launch-pad - " + title))
	sb.WriteString("\n\n")

	labels := []string{
		i18n.T("cat.form.label.id"),
		i18n.T("cat.form.label.name"),
		i18n.T("cat.form.label.icon"),
		i18n.T("cat.form.label.color"),
	}
	for i, label := range labels {
		if i < len(m.inputs) {
			labelStr := styles.InputLabel.Render(label)
			sb.WriteString(fmt.Sprintf("%s  %s\n", labelStr, m.inputs[i].View()))
		}
	}

	sb.WriteString("\n")
	// 保存ボタン
	if m.focusIdx == catFieldCount {
		sb.WriteString(styles.CardFocused.Copy().Width(20).Height(1).Render(i18n.T("cat.form.save_btn")))
	} else {
		sb.WriteString(styles.CardNormal.Copy().Width(20).Height(1).Render(i18n.T("cat.form.save_btn")))
	}
	sb.WriteString("\n\n")
	sb.WriteString(styles.TabInactive.Render(i18n.T("cat.form.hint")))

	return sb.String()
}

func (m CategoryViewModel) viewConfirm() string {
	if m.deleteTarget == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(styles.AppTitle.Render("⌨  cmd-launch-pad - " + i18n.T("cat.confirm.title")))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf(i18n.T("cat.confirm.msg"),
		styles.ErrorStyle.Render(m.deleteTarget.ID),
		m.deleteTarget.Icon,
		m.deleteTarget.Name,
	))

	if m.cmdCountOfDel > 0 {
		sb.WriteString(fmt.Sprintf(i18n.T("cat.confirm.has_cmds"),
			styles.ErrorStyle.Render(fmt.Sprintf("%d", m.cmdCountOfDel)),
		))
		sb.WriteString(i18n.T("cat.confirm.choose"))

		// 選択肢
		opt0 := fmt.Sprintf(i18n.T("cat.confirm.with_cmds"), m.cmdCountOfDel)
		opt1 := i18n.T("cat.confirm.without_cmds")
		var btn0, btn1 string
		if m.deleteConfirm == 0 {
			btn0 = styles.CardFocused.Copy().Width(30).Render("▶ " + opt0)
			btn1 = styles.CardNormal.Copy().Width(30).Render("  " + opt1)
		} else {
			btn0 = styles.CardNormal.Copy().Width(30).Render("  " + opt0)
			btn1 = styles.CardFocused.Copy().Width(30).Render("▶ " + opt1)
		}
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, btn0, "  ", btn1))
		sb.WriteString("\n\n")
	} else {
		sb.WriteString(i18n.T("cat.confirm.no_cmds"))
		sb.WriteString(styles.CardFocused.Copy().Width(20).Height(1).Render(i18n.T("cat.confirm.do_delete")))
		sb.WriteString("\n\n")
	}

	sb.WriteString(styles.TabInactive.Render(i18n.T("cat.confirm.hint")))

	return styles.DialogBox.Render(sb.String())
}
