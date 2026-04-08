package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourname/cmd-launch-pad/internal/tui/styles"
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

var gitMenuItems = []struct {
	label  string
	action GitAction
}{
	{"📂 Gitリポジトリ初期化", GitActionInit},
	{"📋 ステータス確認", GitActionStatus},
	{"💾 コミット", GitActionCommit},
	{"⬆️  プッシュ", GitActionPush},
	{"⬇️  プル", GitActionPull},
	{"🔗 リモートURL設定", GitActionRemote},
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
			if m.cursor < len(gitMenuItems)-1 {
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
	action := gitMenuItems[m.cursor].action
	switch action {
	case GitActionCommit:
		m.isInputMode = true
		m.inputAction = GitActionCommit
		m.input.Placeholder = "コミットメッセージを入力..."
		m.input.Focus()
		return textinput.Blink
	case GitActionRemote:
		m.isInputMode = true
		m.inputAction = GitActionRemote
		m.input.Placeholder = "リモートURL (例: https://github.com/user/repo.git)"
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

	sb.WriteString(styles.AppTitle.Render("🌿  Git操作"))
	sb.WriteString("\n\n")

	// メニュー
	for i, item := range gitMenuItems {
		if i == m.cursor {
			sb.WriteString(styles.CardFocused.Copy().
				Width(32).Height(1).
				Render(fmt.Sprintf(" %s ", item.label)))
		} else {
			sb.WriteString(styles.CardNormal.Copy().
				Width(32).Height(1).
				Render(fmt.Sprintf(" %s ", item.label)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// 入力モード
	if m.isInputMode {
		sb.WriteString(styles.InputLabel.Render("入力:"))
		sb.WriteString("  ")
		sb.WriteString(m.input.View())
		sb.WriteString("\n")
		sb.WriteString(styles.TabInactive.Render("Enter: 実行  Esc: キャンセル"))
		sb.WriteString("\n\n")
	}

	// ステータステキスト
	if m.statusText != "" {
		sb.WriteString(styles.SuccessStyle.Render(m.statusText))
		sb.WriteString("\n\n")
	}

	sb.WriteString(styles.TabInactive.Render("↑↓/jk: 移動  Enter: 実行  q/Esc: 閉じる"))

	return sb.String()
}

// FormatGitStatus はGitステータスを見やすくフォーマットする
func FormatGitStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "✅ 変更なし (クリーン)"
	}
	var sb strings.Builder
	sb.WriteString("📋 変更あり:\n")
	for _, line := range strings.Split(strings.TrimSpace(status), "\n") {
		if line != "" {
			sb.WriteString("  " + line + "\n")
		}
	}
	return sb.String()
}
