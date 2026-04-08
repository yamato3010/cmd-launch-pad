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

// DetailMode は詳細/編集画面のモード
type DetailMode int

const (
	DetailModeNew  DetailMode = iota // 新規作成
	DetailModeEdit                   // 編集
)

// DetailDoneMsg は詳細画面の完了メッセージ
type DetailDoneMsg struct {
	Saved   bool
	Command models.Command
}

// DetailModel はコマンド詳細・編集画面のモデル
type DetailModel struct {
	mode       DetailMode
	original   *models.Command // 編集時の元データ
	categories []models.Category

	// 入力フィールド
	inputs    []textinput.Model
	focusIdx  int
	inputKeys []string // フィールド名

	width  int
	height int
}

const (
	fieldName     = 0
	fieldCommand  = 1
	fieldArgs     = 2
	fieldDesc     = 3
	fieldCategory = 4
	fieldIcon     = 5
	fieldCount    = 6
)

// NewDetailModel は新規作成用のDetailModelを生成する
func NewDetailModel(categories []models.Category) DetailModel {
	return newDetailModel(DetailModeNew, nil, categories)
}

// NewEditModel は編集用のDetailModelを生成する
func NewEditModel(cmd *models.Command, categories []models.Category) DetailModel {
	return newDetailModel(DetailModeEdit, cmd, categories)
}

func newDetailModel(mode DetailMode, cmd *models.Command, categories []models.Category) DetailModel {
	inputs := make([]textinput.Model, fieldCount)
	placeholders := []string{"コマンド名", "コマンド (例: nvim)", "引数 (スペース区切り)", "説明", "カテゴリID", "アイコン (例: 🖊️)"}
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 100
		inputs[i] = ti
	}

	if cmd != nil {
		inputs[fieldName].SetValue(cmd.Name)
		inputs[fieldCommand].SetValue(cmd.Command)
		inputs[fieldArgs].SetValue(strings.Join(cmd.Args, " "))
		inputs[fieldDesc].SetValue(cmd.Description)
		inputs[fieldCategory].SetValue(cmd.CategoryID)
		inputs[fieldIcon].SetValue(cmd.Icon)
	}

	inputs[fieldName].Focus()

	return DetailModel{
		mode:       mode,
		original:   cmd,
		categories: categories,
		inputs:     inputs,
		focusIdx:   0,
		inputKeys:  []string{"名前", "コマンド", "引数", "説明", "カテゴリID", "アイコン"},
	}
}

// Init はDetailModelの初期化コマンドを返す
func (m DetailModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update はキー入力を処理する
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return DetailDoneMsg{Saved: false} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.focusIdx == fieldCount-1 {
				// 最後のフィールドでEnterなら保存
				return m, m.save()
			}
			m.focusIdx++
			return m, m.updateFocus()
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "down"))):
			m.focusIdx = (m.focusIdx + 1) % (fieldCount + 1) // +1 for save button
			if m.focusIdx == fieldCount {
				// セーブボタンにフォーカス
				for i := range m.inputs {
					m.inputs[i].Blur()
				}
				return m, nil
			}
			return m, m.updateFocus()
		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))):
			m.focusIdx--
			if m.focusIdx < 0 {
				m.focusIdx = fieldCount
			}
			if m.focusIdx == fieldCount {
				for i := range m.inputs {
					m.inputs[i].Blur()
				}
				return m, nil
			}
			return m, m.updateFocus()
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			return m, m.save()
		}
		// セーブボタンフォーカス時にEnterで保存
		if m.focusIdx == fieldCount {
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
				return m, m.save()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// アクティブなinputにメッセージを渡す（inputsが未初期化の場合はスキップ）
	if m.focusIdx < fieldCount && len(m.inputs) > m.focusIdx {
		var cmd tea.Cmd
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
		return m, cmd
	}
	return m, nil
}

// updateFocus はフォーカスを更新する
func (m *DetailModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, fieldCount)
	for i := range m.inputs {
		if i == m.focusIdx {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

// save は入力内容を保存してDetailDoneMsgを返す
func (m *DetailModel) save() tea.Cmd {
	cmd := models.Command{
		Name:        m.inputs[fieldName].Value(),
		Command:     m.inputs[fieldCommand].Value(),
		Description: m.inputs[fieldDesc].Value(),
		CategoryID:  m.inputs[fieldCategory].Value(),
		Icon:        m.inputs[fieldIcon].Value(),
	}
	argsStr := m.inputs[fieldArgs].Value()
	if argsStr != "" {
		cmd.Args = strings.Fields(argsStr)
	} else {
		cmd.Args = []string{}
	}
	if m.original != nil {
		cmd.ID = m.original.ID
	}
	return func() tea.Msg {
		return DetailDoneMsg{Saved: true, Command: cmd}
	}
}

// View は詳細・編集画面を描画する（フルスクリーン用 - 後方互換）
func (m DetailModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m DetailModel) ModalView() string {
	var sb strings.Builder

	title := "新規コマンド登録"
	if m.mode == DetailModeEdit {
		title = "コマンド編集"
	}
	sb.WriteString(styles.AppTitle.Render("✚  " + title))
	sb.WriteString("\n\n")

	// カテゴリ一覧のヒントを表示
	if len(m.categories) > 0 {
		catIDs := make([]string, 0, len(m.categories))
		for _, cat := range m.categories {
			catIDs = append(catIDs, fmt.Sprintf("%s%s", cat.Icon, cat.ID))
		}
		sb.WriteString(styles.TabInactive.Render("カテゴリ: " + strings.Join(catIDs, "  ")))
		sb.WriteString("\n\n")
	}

	// 入力フィールド
	for i, label := range m.inputKeys {
		labelStr := styles.InputLabel.Render(label + ":")
		inputStr := m.inputs[i].View()
		sb.WriteString(fmt.Sprintf("%s  %s\n", labelStr, inputStr))
	}

	sb.WriteString("\n")

	// 保存ボタン
	if m.focusIdx == fieldCount {
		sb.WriteString(styles.CardFocused.Copy().Width(22).Height(1).Render("[ 保存 (Ctrl+S) ]"))
	} else {
		sb.WriteString(styles.CardNormal.Copy().Width(22).Height(1).Render("[ 保存 (Ctrl+S) ]"))
	}
	sb.WriteString("\n\n")
	sb.WriteString(styles.TabInactive.Render("Tab: 次へ  Shift+Tab: 前へ  Ctrl+S: 保存  Esc: 閉じる"))

	return sb.String()
}
