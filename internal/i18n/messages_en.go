package i18n

// messagesEn は英語のメッセージマップ
var messagesEn = map[string]string{
	// ── cmd/root.go ──────────────────────────────────────────────────────────
	"root.short": "cmd-launch-pad - TUI command launcher",
	"root.long": `cmd-launch-pad (clp) is a TUI command launcher for terminal users.
Visually manage and launch commands like nvim, lazygit, lazydocker,
just like a GUI "Launchpad".`,

	// ── cmd/add.go ───────────────────────────────────────────────────────────
	"add.short":              "Add a command from CLI",
	"add.long":               "Add a command to the launcher from CLI.",
	"add.flag.name":          "Command name (required)",
	"add.flag.command":       "Executable command (required)",
	"add.flag.category":      "Category ID",
	"add.flag.desc":          "Command description",
	"add.flag.icon":          "Icon (emoji)",
	"add.flag.args":          "Arguments (space-separated)",
	"add.err.required":       "--name and --command are required",
	"add.err.failed":         "failed to add command: %w",
	"add.success":            "✅ Command added: %s (%s)",

	// ── cmd/export.go ────────────────────────────────────────────────────────
	"export.short":           "Export command definitions to YAML",
	"export.long":            "Output registered commands to stdout in YAML format.",
	"export.flag.output":     "Output file path (default: stdout)",
	"export.err.serialize":   "YAML serialization failed: %w",
	"export.err.write":       "File write failed: %w",
	"export.success":         "✅ Export completed: %s",

	// ── cmd/list.go ──────────────────────────────────────────────────────────
	"list.short":             "Interactively select a command and write it to the terminal",
	"list.long": `List registered commands and write the selected one to the terminal.
Just press Enter to execute.

Shell integration setup:
  clp list --shell-init >> ~/.zshrc   # zsh
  clp list --shell-init-bash >> ~/.bashrc  # bash

After setup, call with the _clp_pick command.`,
	"list.flag.shell_init":      "Print zsh shell integration code",
	"list.flag.shell_init_bash": "Print bash shell integration code",
	"list.empty":                "No commands registered. Run `clp` to add some.",
	"list.err.tty":              "could not open terminal: %w",
	"list.header":               "  cmd-launch-pad — Select command",
	"list.hint":                 "  ↑↓ / jk: move   Enter: select   q: cancel",
	"list.footer":               "  Selected command will be written to the terminal",

	// ── cmd/sync.go ──────────────────────────────────────────────────────────
	"sync.short":              "Sync settings via Git",
	"sync.long":               "Sync configuration files using a Git repository.",
	"sync.push.short":         "Push settings to remote",
	"sync.pull.short":         "Pull settings from remote",
	"sync.status.short":       "Show Git status",
	"sync.init.short":         "Initialize config directory as a Git repository",
	"sync.err.not_init":       "Git repository not initialized. Run `clp sync init` first: %w",
	"sync.err.not_init_bare":  "Git repository not initialized: %w",
	"sync.push.success":       "✅ Push completed",
	"sync.pull.success":       "✅ Pull completed",
	"sync.status.clean":       "✅ No changes (clean)",
	"sync.init.success":       "✅ Git repository initialized: %s",

	// ── internal/config/loader.go ────────────────────────────────────────────
	"config.err.homedir":      "failed to get home directory: %w",
	"config.err.mkdir":        "failed to create config directory: %w",
	"config.err.read":         "failed to read config file: %w",
	"config.err.parse":        "failed to parse config file: %w",
	"config.err.serialize":    "failed to serialize config file: %w",
	"config.err.write":        "failed to write config file: %w",

	// ── internal/tui/app.go ──────────────────────────────────────────────────
	"app.err.repo_init":       "failed to initialize repository: %w",
	"app.err.defaults_init":   "failed to initialize default data: %w",
	"app.err.config_load":     "failed to load config file: %w",
	"app.err.exec":            "command execution error: %v",
	"app.err.save":            "save error: %v",
	"app.err.delete":          "delete error: %v",
	"app.err.reload":          "reload error: %v",
	"app.err.category":        "category operation error: %v",
	"app.git.err.prefix":      "❌ Error: %v",
	"app.git.not_init":        "Git repository not initialized",
	"app.git.commit_done":     "✅ Commit completed: %s",
	"app.git.push_done":       "✅ Push completed",
	"app.git.pull_done":       "✅ Pull completed",
	"app.git.remote_done":     "✅ Remote URL set: %s",
	"app.git.init_done":       "✅ Git repository initialized",
	"app.cat.add_done":        "✅ Category added",
	"app.cat.edit_done":       "✅ Category updated",
	"app.cat.delete_done":     "✅ Category deleted",

	// ── views/launcher.go ────────────────────────────────────────────────────
	"launcher.key.move":       "move",
	"launcher.key.exec":       "exec",
	"launcher.key.new":        "new",
	"launcher.key.edit":       "edit",
	"launcher.key.delete":     "delete",
	"launcher.key.tab":        "tab",
	"launcher.key.search":     "search",
	"launcher.key.category":   "category",
	"launcher.key.git":        "Git",
	"launcher.key.help":       "help",
	"launcher.key.quit":       "quit",
	"launcher.tab.all":        "📁 All",

	// ── views/help.go ────────────────────────────────────────────────────────
	"help.title":              "❓  Key Bindings",
	"help.move":               "Move cursor",
	"help.exec":               "Execute command",
	"help.new":                "Register new command",
	"help.edit":               "Edit selected command",
	"help.delete":             "Delete selected command",
	"help.tab":                "Switch category tab",
	"help.search":             "Search mode",
	"help.category":           "Category management",
	"help.git":                "Git operations",
	"help.help":               "Show/hide help",
	"help.quit":               "Quit app",
	"help.close":              "q / Esc / ? to close",

	// ── views/detail.go ──────────────────────────────────────────────────────
	"detail.title.new":        "New Command",
	"detail.title.edit":       "Edit Command",
	"detail.placeholder.name":    "Command name",
	"detail.placeholder.command": "Command (e.g. nvim)",
	"detail.placeholder.args":    "Arguments (space-separated)",
	"detail.placeholder.desc":    "Description",
	"detail.placeholder.cat":     "Category ID",
	"detail.placeholder.icon":    "Icon (e.g. 🖊️)",
	"detail.field.name":       "Name",
	"detail.field.command":    "Command",
	"detail.field.args":       "Args",
	"detail.field.desc":       "Desc",
	"detail.field.cat":        "Category",
	"detail.field.icon":       "Icon",
	"detail.label.categories":"Categories",
	"detail.toggle.label":     "Capture output:",
	"detail.toggle.off":       "[ ] Off (normal exec)",
	"detail.toggle.on":        "[x] On (show result popup)",
	"detail.toggle.hint":      "         Space/Enter to toggle",
	"detail.save_btn":         "[ Save (Ctrl+S) ]",
	"detail.hint":             "Tab: next  Shift+Tab: prev  Ctrl+S: save  Esc: close",

	// ── views/confirm.go ─────────────────────────────────────────────────────
	"confirm.title":           "🗑  Delete Command",
	"confirm.question":        "Are you sure you want to delete the following command?\n\n",
	"confirm.name":            "  Name: %s\n",
	"confirm.command":         "  Command: %s\n",
	"confirm.desc":            "  Desc: %s\n",
	"confirm.irreversible":    "This action cannot be undone.",
	"confirm.yes":             "▶ Yes",
	"confirm.no":              "▶ No",
	"confirm.yes_inactive":    "  Yes",
	"confirm.no_inactive":     "  No",
	"confirm.hint":            "←→/Tab: toggle  Enter: confirm  y: delete  Esc/n: cancel",

	// ── views/search.go ──────────────────────────────────────────────────────
	"search.title":            "🔍  Search",
	"search.placeholder":      "Search by name, description, category...",
	"search.no_results":       "No commands found",
	"search.results":          "%d result(s)",
	"search.hint":             "↑↓/jk: move  Enter: exec  Esc: close",

	// ── views/git.go ─────────────────────────────────────────────────────────
	"git.title":               "🌿  Git Operations",
	"git.menu.init":           "📂 Initialize Git repository",
	"git.menu.status":         "📋 Check status",
	"git.menu.commit":         "💾 Commit",
	"git.menu.push":           "⬆️  Push",
	"git.menu.pull":           "⬇️  Pull",
	"git.menu.remote":         "🔗 Set remote URL",
	"git.input.label":         "Input:",
	"git.commit.placeholder":  "Enter commit message...",
	"git.remote.placeholder":  "Remote URL (e.g. https://github.com/user/repo.git)",
	"git.hint":                "↑↓/jk: move  Enter: exec  q/Esc: close",
	"git.input.hint":          "Enter: exec  Esc: cancel",
	"git.status.clean":        "✅ No changes (clean)",
	"git.status.changed":      "📋 Changes:\n",

	// ── views/output.go ──────────────────────────────────────────────────────
	"output.title":            "🖥  Result",
	"output.title_with_name":  "🖥  Result: ",
	"output.status_ok":        "Exit status: OK",
	"output.status_err":       "Exit status: Error",
	"output.empty":            "(no output)",
	"output.error_prefix":     "Error: ",
	"output.remaining_lines":  "\n... (%d more lines)",
	"output.scroll_hint":      "↑↓/PgUp/PgDn: scroll  %d%%",
	"output.close_hint":       "q / Esc / Enter to close",

	// ── views/category.go ────────────────────────────────────────────────────
	"cat.title.list":          "🗂  Category Management",
	"cat.empty":               "No categories. Press n to add.",
	"cat.header.id":           "ID",
	"cat.header.name":         "Name",
	"cat.header.icon":         "Icon",
	"cat.header.color":        "Color",
	"cat.header.cmd":          "Cmd",
	"cat.list.hint":           "↑↓/jk: move  n: add  e: edit  d: delete  q/Esc: close",
	"cat.form.title.add":      "Add Category",
	"cat.form.title.edit":     "Edit Category",
	"cat.form.placeholder.id":    "Category ID (e.g. editor)",
	"cat.form.placeholder.name":  "Category name (e.g. Editors)",
	"cat.form.placeholder.icon":  "Icon (e.g. ✏️)",
	"cat.form.placeholder.color": "Color (e.g. #7aa2f7)",
	"cat.form.label.id":       "ID:",
	"cat.form.label.name":     "Name:",
	"cat.form.label.icon":     "Icon:",
	"cat.form.label.color":    "Color:",
	"cat.form.save_btn":       "[ Save (Ctrl+S) ]",
	"cat.form.hint":           "Tab/Enter: next  Ctrl+S: save  Esc: cancel",
	"cat.confirm.title":       "Delete Category",
	"cat.confirm.msg":         "Delete category %s \"%s %s\".\n\n",
	"cat.confirm.has_cmds":    "This category has %s command(s).\n",
	"cat.confirm.choose":      "Choose how to delete:\n\n",
	"cat.confirm.with_cmds":   "Delete with commands (%d)",
	"cat.confirm.without_cmds":"Delete category only (keep commands)",
	"cat.confirm.no_cmds":     "No commands in this category.\n\n",
	"cat.confirm.do_delete":   "▶ Delete",
	"cat.confirm.hint":        "←→/Tab: toggle  Enter: confirm  Esc: cancel",

	// ── components/statusbar.go ──────────────────────────────────────────────
	"desc.add_hint":           "Press Enter to add a new command",
}
