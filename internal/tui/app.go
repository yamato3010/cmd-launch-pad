package tui

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yamato3010/cmd-launch-pad/internal/config"
	gitpkg "github.com/yamato3010/cmd-launch-pad/internal/git"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/repository"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/components"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/views"
)

// ViewState は現在表示中の画面状態
type ViewState int

const (
	ViewLauncher ViewState = iota
	ViewDetail
	ViewSearch
	ViewGit
	ViewHelp
	ViewCategory
)

// execCommandMsg はコマンド実行後のメッセージ
type execCommandMsg struct {
	err error
}

// gitResultMsg はGit操作完了後のメッセージ
type gitResultMsg struct {
	text string
	err  error
}

// categoryResultMsg はカテゴリ操作完了後のメッセージ
type categoryResultMsg struct {
	text string
	err  error
}

// App はTUIアプリのルートモデル
type App struct {
	state      ViewState
	repo       *repository.CommandRepository
	appCfg     *config.AppConfig
	commands   []models.Command
	categories []models.Category
	gitMgr     *gitpkg.GitManager // nil の場合はGit未初期化

	// サブビュー
	launcher views.LauncherModel
	detail   views.DetailModel
	search   views.SearchModel
	gitView  views.GitViewModel
	help     views.HelpModel
	catView  views.CategoryViewModel

	width  int
	height int
	err    string // エラーメッセージ
}

// NewApp は新しいAppを生成して初期化する
func NewApp() (*App, error) {
	repo, err := repository.NewCommandRepository()
	if err != nil {
		return nil, fmt.Errorf("リポジトリの初期化に失敗しました: %w", err)
	}
	if err := repo.InitDefaults(); err != nil {
		return nil, fmt.Errorf("デフォルトデータの初期化に失敗しました: %w", err)
	}

	appCfg, err := config.LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}
	// デフォルト設定を保存（初回のみ）
	if err := config.SaveAppConfig(appCfg); err != nil {
		return nil, err
	}

	commands, err := repo.ListCommands()
	if err != nil {
		return nil, err
	}
	categories, err := repo.ListCategories()
	if err != nil {
		return nil, err
	}

	cols := appCfg.Columns
	if cols <= 0 {
		cols = 4
	}

	// Git managerを試みる（エラーは無視）
	var gitMgr *gitpkg.GitManager
	cfgDir, _ := config.ConfigDir()
	gitMgr, _ = gitpkg.NewGitManager(cfgDir)

	app := &App{
		state:      ViewLauncher,
		repo:       repo,
		appCfg:     appCfg,
		commands:   commands,
		categories: categories,
		gitMgr:     gitMgr,
		launcher:   views.NewLauncherModel(commands, categories, cols),
		help:       views.NewHelpModel(),
		gitView:    views.NewGitViewModel(),
	}
	return app, nil
}

// Init はAppの初期化コマンドを返す
func (a App) Init() tea.Cmd {
	return a.launcher.Init()
}

// Update はメッセージを処理してモデルを更新する
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// 各サブビューにも伝播
		a.launcher, _ = a.launcher.Update(msg)
		a.detail, _ = a.detail.Update(msg)
		a.search, _ = a.search.Update(msg)
		a.gitView, _ = a.gitView.Update(msg)
		a.help, _ = a.help.Update(msg)
		a.catView, _ = a.catView.Update(msg)
		return a, nil

	case tea.KeyMsg:
		// 終了キー（どの画面でも有効）
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
		if a.state == ViewLauncher && msg.String() == "q" {
			return a, tea.Quit
		}

	case execCommandMsg:
		if msg.err != nil {
			a.err = fmt.Sprintf("コマンド実行エラー: %v", msg.err)
		} else {
			a.err = ""
		}
		return a, nil

	case gitResultMsg:
		text := msg.text
		if msg.err != nil {
			text = fmt.Sprintf("❌ エラー: %v", msg.err)
		}
		a.gitView.SetStatus(text)
		return a, nil

	case categoryResultMsg:
		// ここでリロード（Updateのバリューレシーバ内なのでa自体を更新できる）
		if msg.err == nil {
			a.reloadCommands()
		}
		// 最新データでcatViewを再構築
		cmdCounts := make(map[string]int)
		for _, cmd := range a.commands {
			cmdCounts[cmd.CategoryID]++
		}
		a.catView = views.NewCategoryViewModel(a.categories, cmdCounts)
		if a.width > 0 {
			a.catView, _ = a.catView.Update(tea.WindowSizeMsg{Width: a.width, Height: a.height})
		}
		// エラーがあればcatView内にインライン表示
		if msg.err != nil {
			a.catView.SetError(fmt.Sprintf("カテゴリ操作エラー: %v", msg.err))
		}
		return a, nil
	}

	// 現在のビューにメッセージをルーティング
	switch a.state {
	case ViewLauncher:
		return a.updateLauncher(msg)
	case ViewDetail:
		return a.updateDetail(msg)
	case ViewSearch:
		return a.updateSearch(msg)
	case ViewGit:
		return a.updateGit(msg)
	case ViewHelp:
		return a.updateHelp(msg)
	case ViewCategory:
		return a.updateCategory(msg)
	}
	return a, nil
}

func (a App) updateLauncher(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.launcher, cmd = a.launcher.Update(msg)

	// LauncherMsgを処理
	if lMsg, ok := msg.(views.LauncherMsg); ok {
		switch lMsg.Action {
		case views.LauncherActionExec:
			if lMsg.Command != nil {
				return a, a.execCommand(lMsg.Command)
			}
		case views.LauncherActionNew:
			a.detail = views.NewDetailModel(a.categories)
			a.state = ViewDetail
			return a, a.detail.Init()
		case views.LauncherActionEdit:
			if lMsg.Command != nil {
				a.detail = views.NewEditModel(lMsg.Command, a.categories)
				a.state = ViewDetail
				return a, a.detail.Init()
			}
		case views.LauncherActionDelete:
			if lMsg.Command != nil {
				if err := a.repo.DeleteCommand(lMsg.Command.ID); err != nil {
					a.err = fmt.Sprintf("削除エラー: %v", err)
				} else {
					a.err = ""
					a.reloadCommands()
				}
			}
		case views.LauncherActionSearch:
			a.search = views.NewSearchModel(a.commands, a.categories)
			a.state = ViewSearch
			return a, a.search.Init()
		case views.LauncherActionGit:
			a.gitView = views.NewGitViewModel()
			a.state = ViewGit
			return a, a.gitView.Init()
		case views.LauncherActionHelp:
			a.state = ViewHelp
			return a, a.help.Init()
		case views.LauncherActionCategory:
			return a, a.openCategoryView()
		}
	}
	return a, cmd
}

func (a App) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.detail, cmd = a.detail.Update(msg)

	if dMsg, ok := msg.(views.DetailDoneMsg); ok {
		if dMsg.Saved {
			cmdData := dMsg.Command
			var err error
			if cmdData.ID == "" {
				err = a.repo.AddCommand(&cmdData)
			} else {
				err = a.repo.UpdateCommand(&cmdData)
			}
			if err != nil {
				a.err = fmt.Sprintf("保存エラー: %v", err)
			} else {
				a.err = ""
				a.reloadCommands()
			}
		}
		a.state = ViewLauncher
		return a, nil
	}
	return a, cmd
}

func (a App) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.search, cmd = a.search.Update(msg)

	if sMsg, ok := msg.(views.SearchDoneMsg); ok {
		if sMsg.Selected != nil {
			a.state = ViewLauncher
			return a, a.execCommand(sMsg.Selected)
		}
		a.state = ViewLauncher
		return a, nil
	}
	return a, cmd
}

func (a App) updateGit(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.gitView, cmd = a.gitView.Update(msg)

	if _, ok := msg.(views.BackMsg); ok {
		a.state = ViewLauncher
		return a, nil
	}
	if gMsg, ok := msg.(views.GitDoMsg); ok {
		return a, a.handleGitAction(gMsg)
	}
	return a, cmd
}

func (a App) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.help, cmd = a.help.Update(msg)

	if _, ok := msg.(views.BackMsg); ok {
		a.state = ViewLauncher
		return a, nil
	}
	return a, cmd
}

// View は現在のビューを描画する
// ランチャー以外の画面はランチャーを背景にモーダルとして重ねて表示する
func (a App) View() string {
	launcherBg := a.launcher.View()

	switch a.state {
	case ViewLauncher:
		return launcherBg
	default:
		modal := a.modalContent()
		if modal == "" {
			return launcherBg
		}
		boxed := components.ModalBox.Render(modal)
		return components.PlaceOverlay(launcherBg, boxed, a.width, a.height-1)
	}
}

// modalContent はモーダルとして表示するコンテンツを返す（ボーダーなし）
func (a App) modalContent() string {
	switch a.state {
	case ViewDetail:
		return a.detail.ModalView()
	case ViewSearch:
		return a.search.ModalView()
	case ViewGit:
		return a.gitView.ModalView()
	case ViewHelp:
		return a.help.ModalView()
	case ViewCategory:
		return a.catView.ModalView()
	}
	return ""
}

// openCategoryView はカテゴリ管理画面を開く
func (a *App) openCategoryView() tea.Cmd {
	// コマンド数マップを作成
	cmdCounts := make(map[string]int)
	for _, cmd := range a.commands {
		cmdCounts[cmd.CategoryID]++
	}
	a.catView = views.NewCategoryViewModel(a.categories, cmdCounts)
	a.state = ViewCategory
	return a.catView.Init()
}

// updateCategory はカテゴリ管理画面のメッセージを処理する
func (a App) updateCategory(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.catView, cmd = a.catView.Update(msg)

	if _, ok := msg.(views.BackMsg); ok {
		a.state = ViewLauncher
		return a, nil
	}

	if cMsg, ok := msg.(views.CategoryDoneMsg); ok {
		return a, a.handleCategoryAction(cMsg)
	}
	return a, cmd
}

// handleCategoryAction はカテゴリ操作を実行する
// NOTE: goroutine内でa.reloadCommands()を呼ばないこと。
// リロードはcategoryResultMsgを受け取ったUpdate内で行う。
func (a *App) handleCategoryAction(msg views.CategoryDoneMsg) tea.Cmd {
	switch msg.Action {
	case views.CategoryActionAdd:
		cat := msg.Category
		repo := a.repo
		return func() tea.Msg {
			if err := repo.AddCategory(&cat); err != nil {
				return categoryResultMsg{err: err}
			}
			return categoryResultMsg{text: "✅ カテゴリを追加しました"}
		}
	case views.CategoryActionEdit:
		cat := msg.Category
		repo := a.repo
		return func() tea.Msg {
			if err := repo.UpdateCategory(&cat); err != nil {
				return categoryResultMsg{err: err}
			}
			return categoryResultMsg{text: "✅ カテゴリを更新しました"}
		}
	case views.CategoryActionDelete:
		cat := msg.Category
		withCmds := msg.WithCommands
		repo := a.repo
		return func() tea.Msg {
			if err := repo.DeleteCategory(cat.ID, withCmds); err != nil {
				return categoryResultMsg{err: err}
			}
			return categoryResultMsg{text: "✅ カテゴリを削除しました"}
		}
	}
	return nil
}

// reloadCommands はリポジトリからコマンド・カテゴリを再読み込みする
func (a *App) reloadCommands() {
	cmds, err := a.repo.ListCommands()
	if err != nil {
		a.err = fmt.Sprintf("再読み込みエラー: %v", err)
		return
	}
	cats, err := a.repo.ListCategories()
	if err != nil {
		a.err = fmt.Sprintf("再読み込みエラー: %v", err)
		return
	}
	a.commands = cmds
	a.categories = cats
	a.launcher.SetCommands(cmds)
	a.launcher.SetCategories(cats)
}

// execCommand はコマンドをExecProcessで実行する (TUIを一時停止)
func (a *App) execCommand(cmd *models.Command) tea.Cmd {
	args := append([]string{}, cmd.Args...)
	c := exec.Command(cmd.Command, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return execCommandMsg{err: err}
	})
}

// handleGitAction はGit操作を実行する
func (a *App) handleGitAction(msg views.GitDoMsg) tea.Cmd {
	cfgDir, err := config.ConfigDir()
	if err != nil {
		return func() tea.Msg { return gitResultMsg{err: err} }
	}

	switch msg.Action {
	case views.GitActionInit:
		return func() tea.Msg {
			mgr, err := gitpkg.Init(cfgDir)
			if err != nil {
				return gitResultMsg{err: err}
			}
			a.gitMgr = mgr
			return gitResultMsg{text: "✅ Gitリポジトリを初期化しました"}
		}

	case views.GitActionStatus:
		return func() tea.Msg {
			if a.gitMgr == nil {
				return gitResultMsg{err: fmt.Errorf("Gitリポジトリが初期化されていません")}
			}
			status, err := a.gitMgr.Status()
			if err != nil {
				return gitResultMsg{err: err}
			}
			return gitResultMsg{text: views.FormatGitStatus(status)}
		}

	case views.GitActionCommit:
		return func() tea.Msg {
			if a.gitMgr == nil {
				return gitResultMsg{err: fmt.Errorf("Gitリポジトリが初期化されていません")}
			}
			if err := a.gitMgr.AddAll(); err != nil {
				return gitResultMsg{err: err}
			}
			commitMsg := msg.Payload
			if commitMsg == "" {
				commitMsg = "Update commands"
			}
			if err := a.gitMgr.Commit(commitMsg); err != nil {
				return gitResultMsg{err: err}
			}
			return gitResultMsg{text: fmt.Sprintf("✅ コミット完了: %s", commitMsg)}
		}

	case views.GitActionPush:
		return func() tea.Msg {
			if a.gitMgr == nil {
				return gitResultMsg{err: fmt.Errorf("Gitリポジトリが初期化されていません")}
			}
			if err := a.gitMgr.Push("origin", a.appCfg.Git.Branch, nil); err != nil {
				return gitResultMsg{err: err}
			}
			return gitResultMsg{text: "✅ プッシュ完了"}
		}

	case views.GitActionPull:
		return func() tea.Msg {
			if a.gitMgr == nil {
				return gitResultMsg{err: fmt.Errorf("Gitリポジトリが初期化されていません")}
			}
			if err := a.gitMgr.Pull("origin", a.appCfg.Git.Branch, nil); err != nil {
				return gitResultMsg{err: err}
			}
			a.reloadCommands()
			return gitResultMsg{text: "✅ プル完了"}
		}

	case views.GitActionRemote:
		return func() tea.Msg {
			if a.gitMgr == nil {
				return gitResultMsg{err: fmt.Errorf("Gitリポジトリが初期化されていません")}
			}
			if err := a.gitMgr.SetRemote("origin", msg.Payload); err != nil {
				return gitResultMsg{err: err}
			}
			a.appCfg.Git.Remote = msg.Payload
			_ = config.SaveAppConfig(a.appCfg)
			return gitResultMsg{text: fmt.Sprintf("✅ リモートURLを設定しました: %s", msg.Payload)}
		}
	}
	return nil
}

// Run はTUIアプリを起動する
func Run() error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
