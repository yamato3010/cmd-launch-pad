package i18n

// messagesJa は日本語のメッセージマップ
var messagesJa = map[string]string{
	// ── cmd/root.go ──────────────────────────────────────────────────────────
	"root.short": "cmd-launch-pad - TUIコマンドランチャー",
	"root.long": `cmd-launch-pad (clp) はターミナルユーザー向けのTUIコマンドランチャーです。
nvim、lazygit、lazydockerなどのコマンドをGUIの「Launchpad」のように
視覚的に管理・起動できます。`,

	// ── cmd/add.go ───────────────────────────────────────────────────────────
	"add.short":              "コマンドをCLIから追加する",
	"add.long":               "CLIからコマンドランチャーにコマンドを追加します。",
	"add.flag.name":          "コマンド名 (必須)",
	"add.flag.command":       "実行コマンド (必須)",
	"add.flag.category":      "カテゴリID",
	"add.flag.desc":          "コマンドの説明",
	"add.flag.icon":          "アイコン (絵文字)",
	"add.flag.args":          "引数 (スペース区切り)",
	"add.err.required":       "--name と --command は必須です",
	"add.err.failed":         "コマンドの追加に失敗しました: %w",
	"add.success":            "✅ コマンドを追加しました: %s (%s)",

	// ── cmd/export.go ────────────────────────────────────────────────────────
	"export.short":           "コマンド定義をYAML形式でエクスポート",
	"export.long":            "登録済みコマンドをYAML形式で標準出力に出力します。",
	"export.flag.output":     "出力ファイルパス (省略時は標準出力)",
	"export.err.serialize":   "YAMLシリアライズに失敗しました: %w",
	"export.err.write":       "ファイル書き込みに失敗しました: %w",
	"export.success":         "✅ エクスポート完了: %s",

	// ── cmd/list.go ──────────────────────────────────────────────────────────
	"list.short":             "コマンドをインタラクティブに選択してターミナルに書き込む",
	"list.long": `登録済みコマンドを一覧表示し、選択したコマンドをターミナルに書き込みます。
Enter を押すだけで実行できます。

シェル統合のセットアップ:
  clp list --shell-init >> ~/.zshrc   # zsh
  clp list --shell-init-bash >> ~/.bashrc  # bash

セットアップ後は _clp_pick コマンドで呼び出せます。`,
	"list.flag.shell_init":      "zsh用シェル統合コードを出力する",
	"list.flag.shell_init_bash": "bash用シェル統合コードを出力する",
	"list.empty":                "コマンドが登録されていません。`clp` を起動して追加してください。",
	"list.err.tty":              "ターミナルを開けませんでした: %w",
	"list.header":               "  cmd-launch-pad — コマンド選択",
	"list.hint":                 "  ↑↓ / jk: 移動   Enter: 選択   q: キャンセル",
	"list.footer":               "  選択するとコマンドがターミナルに書き込まれます",

	// ── cmd/sync.go ──────────────────────────────────────────────────────────
	"sync.short":              "Gitによる設定同期",
	"sync.long":               "Gitリポジトリを使って設定ファイルを同期します。",
	"sync.push.short":         "設定をリモートにプッシュ",
	"sync.pull.short":         "リモートから設定をプル",
	"sync.status.short":       "Gitステータスを表示",
	"sync.init.short":         "設定ディレクトリをGitリポジトリとして初期化",
	"sync.err.not_init":       "Gitリポジトリが初期化されていません。先に `clp sync init` を実行してください: %w",
	"sync.err.not_init_bare":  "Gitリポジトリが初期化されていません: %w",
	"sync.push.success":       "✅ プッシュ完了",
	"sync.pull.success":       "✅ プル完了",
	"sync.status.clean":       "✅ 変更なし (クリーン)",
	"sync.init.success":       "✅ Gitリポジトリを初期化しました: %s",

	// ── internal/config/loader.go ────────────────────────────────────────────
	"config.err.homedir":      "ホームディレクトリの取得に失敗しました: %w",
	"config.err.mkdir":        "設定ディレクトリの作成に失敗しました: %w",
	"config.err.read":         "設定ファイルの読み込みに失敗しました: %w",
	"config.err.parse":        "設定ファイルのパースに失敗しました: %w",
	"config.err.serialize":    "設定ファイルのシリアライズに失敗しました: %w",
	"config.err.write":        "設定ファイルの書き込みに失敗しました: %w",

	// ── internal/tui/app.go ──────────────────────────────────────────────────
	"app.err.repo_init":       "リポジトリの初期化に失敗しました: %w",
	"app.err.defaults_init":   "デフォルトデータの初期化に失敗しました: %w",
	"app.err.config_load":     "設定ファイルの読み込みに失敗しました: %w",
	"app.err.exec":            "コマンド実行エラー: %v",
	"app.err.save":            "保存エラー: %v",
	"app.err.delete":          "削除エラー: %v",
	"app.err.reload":          "再読み込みエラー: %v",
	"app.err.category":        "カテゴリ操作エラー: %v",
	"app.git.err.prefix":      "❌ エラー: %v",
	"app.git.not_init":        "Gitリポジトリが初期化されていません",
	"app.git.commit_done":     "✅ コミット完了: %s",
	"app.git.push_done":       "✅ プッシュ完了",
	"app.git.pull_done":       "✅ プル完了",
	"app.git.remote_done":     "✅ リモートURLを設定しました: %s",
	"app.git.init_done":       "✅ Gitリポジトリを初期化しました",
	"app.cat.add_done":        "✅ カテゴリを追加しました",
	"app.cat.edit_done":       "✅ カテゴリを更新しました",
	"app.cat.delete_done":     "✅ カテゴリを削除しました",

	// ── views/launcher.go ────────────────────────────────────────────────────
	"launcher.key.move":       "移動",
	"launcher.key.exec":       "実行",
	"launcher.key.new":        "新規",
	"launcher.key.edit":       "編集",
	"launcher.key.delete":     "削除",
	"launcher.key.tab":        "タブ切替",
	"launcher.key.search":     "検索",
	"launcher.key.category":   "カテゴリ",
	"launcher.key.git":        "Git",
	"launcher.key.help":       "ヘルプ",
	"launcher.key.quit":       "終了",
	"launcher.tab.all":        "📁 全て",

	// ── views/help.go ────────────────────────────────────────────────────────
	"help.title":              "❓  キーバインド一覧",
	"help.move":               "カーソル移動",
	"help.exec":               "コマンド実行",
	"help.new":                "新規コマンド登録",
	"help.edit":               "選択中のコマンド編集",
	"help.delete":             "選択中のコマンド削除",
	"help.tab":                "カテゴリタブ切り替え",
	"help.search":             "検索モード",
	"help.category":           "カテゴリ管理",
	"help.git":                "Git操作画面",
	"help.help":               "ヘルプ表示/非表示",
	"help.quit":               "アプリ終了",
	"help.close":              "q / Esc / ? で閉じる",

	// ── views/detail.go ──────────────────────────────────────────────────────
	"detail.title.new":        "新規コマンド登録",
	"detail.title.edit":       "コマンド編集",
	"detail.placeholder.name":    "コマンド名",
	"detail.placeholder.command": "コマンド (例: nvim)",
	"detail.placeholder.args":    "引数 (スペース区切り)",
	"detail.placeholder.desc":    "説明",
	"detail.placeholder.cat":     "カテゴリID",
	"detail.placeholder.icon":    "アイコン (例: 🖊️)",
	"detail.field.name":       "名前",
	"detail.field.command":    "コマンド",
	"detail.field.args":       "引数",
	"detail.field.desc":       "説明",
	"detail.field.cat":        "カテゴリID",
	"detail.field.icon":       "アイコン",
	"detail.label.categories":"カテゴリ",
	"detail.toggle.label":     "出力キャプチャ:",
	"detail.toggle.off":       "[ ] オフ (通常実行)",
	"detail.toggle.on":        "[x] オン (結果ポップアップ表示)",
	"detail.toggle.hint":      "         Space/Enter でオン/オフ切り替え",
	"detail.save_btn":         "[ 保存 (Ctrl+S) ]",
	"detail.hint":             "Tab: 次へ  Shift+Tab: 前へ  Ctrl+S: 保存  Esc: 閉じる",

	// ── views/confirm.go ─────────────────────────────────────────────────────
	"confirm.title":           "🗑  コマンドの削除",
	"confirm.question":        "以下のコマンドを削除してもよいですか？\n\n",
	"confirm.name":            "  名前: %s\n",
	"confirm.command":         "  コマンド: %s\n",
	"confirm.desc":            "  説明: %s\n",
	"confirm.irreversible":    "この操作は元に戻せません。",
	"confirm.yes":             "▶ はい",
	"confirm.no":              "▶ いいえ",
	"confirm.yes_inactive":    "  はい",
	"confirm.no_inactive":     "  いいえ",
	"confirm.hint":            "←→/Tab: 選択切替  Enter: 実行  y: 削除  Esc/n: キャンセル",

	// ── views/search.go ──────────────────────────────────────────────────────
	"search.title":            "🔍  検索",
	"search.placeholder":      "コマンド名・説明・カテゴリで検索...",
	"search.no_results":       "該当するコマンドが見つかりません",
	"search.results":          "%d 件ヒット",
	"search.hint":             "↑↓/jk: 移動  Enter: 実行  Esc: 閉じる",

	// ── views/git.go ─────────────────────────────────────────────────────────
	"git.title":               "🌿  Git操作",
	"git.menu.init":           "📂 Gitリポジトリ初期化",
	"git.menu.status":         "📋 ステータス確認",
	"git.menu.commit":         "💾 コミット",
	"git.menu.push":           "⬆️  プッシュ",
	"git.menu.pull":           "⬇️  プル",
	"git.menu.remote":         "🔗 リモートURL設定",
	"git.input.label":         "入力:",
	"git.commit.placeholder":  "コミットメッセージを入力...",
	"git.remote.placeholder":  "リモートURL (例: https://github.com/user/repo.git)",
	"git.hint":                "↑↓/jk: 移動  Enter: 実行  q/Esc: 閉じる",
	"git.input.hint":          "Enter: 実行  Esc: キャンセル",
	"git.status.clean":        "✅ 変更なし (クリーン)",
	"git.status.changed":      "📋 変更あり:\n",

	// ── views/output.go ──────────────────────────────────────────────────────
	"output.title":            "🖥  実行結果",
	"output.title_with_name":  "🖥  実行結果: ",
	"output.status_ok":        "終了ステータス: 正常終了",
	"output.status_err":       "終了ステータス: エラー",
	"output.empty":            "(出力なし)",
	"output.error_prefix":     "エラー: ",
	"output.remaining_lines":  "\n... (残り %d 行)",
	"output.scroll_hint":      "↑↓/PgUp/PgDn: スクロール  %d%%",
	"output.close_hint":       "q / Esc / Enter で閉じる",

	// ── views/category.go ────────────────────────────────────────────────────
	"cat.title.list":          "🗂  カテゴリ管理",
	"cat.empty":               "カテゴリがありません。n で新規追加できます。",
	"cat.header.id":           "ID",
	"cat.header.name":         "名前",
	"cat.header.icon":         "アイコン",
	"cat.header.color":        "カラー",
	"cat.header.cmd":          "Cmd",
	"cat.list.hint":           "↑↓/jk: 移動  n: 追加  e: 編集  d: 削除  q/Esc: 閉じる",
	"cat.form.title.add":      "カテゴリ新規追加",
	"cat.form.title.edit":     "カテゴリ編集",
	"cat.form.placeholder.id":    "カテゴリID (例: editor)",
	"cat.form.placeholder.name":  "カテゴリ名 (例: エディタ)",
	"cat.form.placeholder.icon":  "アイコン (例: ✏️)",
	"cat.form.placeholder.color": "カラー (例: #7aa2f7)",
	"cat.form.label.id":       "ID:",
	"cat.form.label.name":     "名前:",
	"cat.form.label.icon":     "アイコン:",
	"cat.form.label.color":    "カラー:",
	"cat.form.save_btn":       "[ 保存 (Ctrl+S) ]",
	"cat.form.hint":           "Tab/Enter: 次へ  Ctrl+S: 保存  Esc: キャンセル",
	"cat.confirm.title":       "カテゴリ削除確認",
	"cat.confirm.msg":         "カテゴリ %s「%s %s」を削除します。\n\n",
	"cat.confirm.has_cmds":    "このカテゴリには %s 件のコマンドが属しています。\n",
	"cat.confirm.choose":      "削除方法を選択してください:\n\n",
	"cat.confirm.with_cmds":   "コマンドも一緒に削除 (%d件)",
	"cat.confirm.without_cmds":"カテゴリのみ削除（コマンドは残す）",
	"cat.confirm.no_cmds":     "このカテゴリにはコマンドが属していません。\n\n",
	"cat.confirm.do_delete":   "▶ 削除する",
	"cat.confirm.hint":        "←→/Tab: 選択切替  Enter: 実行  Esc: キャンセル",

	// ── components/statusbar.go ──────────────────────────────────────────────
	"desc.add_hint":           "新しいコマンドを追加するには Enter を押してください",
}
