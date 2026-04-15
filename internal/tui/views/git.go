package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// GitAction はGit操作の種別
type GitAction int

const (
	GitActionInit   GitAction = iota // Gitリポジトリ初期化
	GitActionStatus                  // ステータス確認
	GitActionCommit                  // コミット
	GitActionPush                    // プッシュ
	GitActionPull                    // プル
	GitActionRemote                  // リモート設定
)

// GitDoMsg はGit操作実行メッセージ
type GitDoMsg struct {
	Action  GitAction
	Payload string // コミットメッセージ or リモートURL など
}

// GitViewModel はGit操作画面のモデル
type GitViewModel struct {
	cursor      int
	statusText  string // 操作結果メッセージ
	isInputMode bool
	input       textinput.Model
	inputAction GitAction
	width       int
	height      int
}

// gitMenuItem はメニュー項目
type gitMenuItem struct {
	labelKey string
	action   GitAction
}

var gitMenuItemDefs = []gitMenuItem{
	{"git.menu.init", GitActionInit},
	{"git.menu.status", GitActionStatus},
	{"git.menu.commit", GitActionCommit},
	{"git.menu.push", GitActionPush},
	{"git.menu.pull", GitActionPull},
	{"git.menu.remote", GitActionRemote},
}

// NewGitViewModel は新しいGitViewModelを生成する
func NewGitViewModel() GitViewModel {
	ti := textinput.New()
	ti.CharLimit = 200
	return GitViewModel{
		input: ti,
	}
}

// SetStatus はステータステキストを設定する
func (m *GitViewModel) SetStatus(text string) {
	m.statusText = text
}

// Init はGitViewModelの初期化コマンドを返す
func (m GitViewModel) Init() tea.Cmd {
	return nil
}

// Update はキー入力を処理する
func (m GitViewModel) Update(msg tea.Msg) (GitViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isInputMode {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.isInputMode = false
				m.input.SetValue("")
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				val := m.input.Value()
				m.isInputMode = false
				m.input.SetValue("")
				return m, func() tea.Msg {
					return GitDoMsg{Action: m.inputAction, Payload: val}
				}
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))):
			return m, func() tea.Msg { return BackMsg{} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(gitMenuItemDefs)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			return m, m.handleSelect()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// handleSelect は選択されたメニュー項目を処理する
func (m *GitViewModel) handleSelect() tea.Cmd {
	action := gitMenuItemDefs[m.cursor].action
	switch action {
	case GitActionCommit:
		m.isInputMode = true
		m.inputAction = GitActionCommit
		m.input.Placeholder = i18n.T("git.commit.placeholder")
		m.input.Focus()
		return textinput.Blink
	case GitActionRemote:
		m.isInputMode = true
		m.inputAction = GitActionRemote
		m.input.Placeholder = i18n.T("git.remote.placeholder")
		m.input.Focus()
		return textinput.Blink
	default:
		return func() tea.Msg {
			return GitDoMsg{Action: action}
		}
	}
}

// View はGit操作画面を描画する
func (m GitViewModel) View() string {
	return m.ModalView()
}

// ModalView はモーダル表示用コンテンツを返す
func (m GitViewModel) ModalView() string {
	var sb strings.Builder

	sb.WriteString(styles.AppTitle.Render(i18n.T("git.title")))
	sb.WriteString("\n\n")

	// メニュー
	for i, item := range gitMenuItemDefs {
		label := i18n.T(item.labelKey)
		if i == m.cursor {
			sb.WriteString(styles.CardFocused.Copy().
				Width(32).Height(1).
				Render(" " + label + " "))
		} else {
			sb.WriteString(styles.CardNormal.Copy().
				Width(32).Height(1).
				Render(" " + label + " "))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// 入力モード
	if m.isInputMode {
		sb.WriteString(styles.InputLabel.Render(i18n.T("git.input.label")))
		sb.WriteString("  ")
		sb.WriteString(m.input.View())
		sb.WriteString("\n")
		sb.WriteString(styles.TabInactive.Render(i18n.T("git.input.hint")))
		sb.WriteString("\n\n")
	}

	// ステータステキスト
	if m.statusText != "" {
		sb.WriteString(styles.SuccessStyle.Render(m.statusText))
		sb.WriteString("\n\n")
	}

	sb.WriteString(styles.TabInactive.Render(i18n.T("git.hint")))

	return sb.String()
}

// FormatGitStatus はGitステータスを見やすくフォーマットする
func FormatGitStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return i18n.T("git.status.clean")
	}
	var sb strings.Builder
	sb.WriteString(i18n.T("git.status.changed"))
	for _, line := range strings.Split(strings.TrimSpace(status), "\n") {
		if line != "" {
			sb.WriteString("  " + line + "\n")
		}
	}
	return sb.String()
}
