package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/repository"
)

// ============================================================
// シェル統合コード
// ============================================================

const shellInitZsh = `# cmd-launch-pad shell integration for zsh
# Add the following to ~/.zshrc

_clp_pick() {
  local selected
  selected=$(clp list)
  if [[ -n "$selected" ]]; then
    print -z -- "$selected"
  fi
}

# Key binding (e.g. Alt+P to pick a command)
# bindkey -s '^[p' '_clp_pick\n'
`

const shellInitBash = `# cmd-launch-pad shell integration for bash
# Add the following to ~/.bashrc

_clp_pick() {
  local selected
  selected=$(clp list)
  if [[ -n "$selected" ]]; then
    READLINE_LINE="$selected"
    READLINE_POINT=${#READLINE_LINE}
  fi
}

# Key binding (e.g. Alt+P to pick a command)
# bind -x '"\ep": _clp_pick'
`

// ============================================================
// スタイル
// ============================================================

var (
	listHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#c0caf5"))

	listSelectedNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7aa2f7"))

	listSelectedCmdStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#9ece6a"))

	listNormalNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a9b1d6"))

	listNormalCmdStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#565f89"))

	listDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#414868"))

	listFooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89"))

	listCursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7aa2f7"))
)

// ============================================================
// Bubbletea モデル
// ============================================================

type listModel struct {
	commands   []models.Command
	categories map[string]models.Category
	cursor     int
	selected   string
}

func newListModel(commands []models.Command, categories []models.Category) listModel {
	catMap := make(map[string]models.Category, len(categories))
	for _, c := range categories {
		catMap[c.ID] = c
	}
	return listModel{
		commands:   commands,
		categories: catMap,
	}
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.commands) > 0 {
				cmd := m.commands[m.cursor]
				parts := []string{cmd.Command}
				parts = append(parts, cmd.Args...)
				m.selected = strings.Join(parts, " ")
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m listModel) View() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(listHeaderStyle.Render(i18n.T("list.header")))
	sb.WriteString("\n")
	sb.WriteString(listDimStyle.Render(i18n.T("list.hint")))
	sb.WriteString("\n\n")

	for i, cmd := range m.commands {
		cmdStr := cmd.Command
		if len(cmd.Args) > 0 {
			cmdStr += " " + strings.Join(cmd.Args, " ")
		}

		catName := ""
		if cat, ok := m.categories[cmd.CategoryID]; ok {
			catName = cat.Icon + " " + cat.Name
		}

		if i == m.cursor {
			cursor := listCursorStyle.Render("▶")
			name := listSelectedNameStyle.Render(fmt.Sprintf("%-18s", cmd.Name))
			command := listSelectedCmdStyle.Render(fmt.Sprintf("%-30s", cmdStr))
			cat := listDimStyle.Render(catName)
			sb.WriteString(fmt.Sprintf("  %s %s  %s  %s\n", cursor, name, command, cat))
		} else {
			name := listNormalNameStyle.Render(fmt.Sprintf("  %-18s", cmd.Name))
			command := listNormalCmdStyle.Render(fmt.Sprintf("%-30s", cmdStr))
			cat := listDimStyle.Render(catName)
			sb.WriteString(fmt.Sprintf("  %s  %s  %s\n", name, command, cat))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(listFooterStyle.Render(i18n.T("list.footer")))
	sb.WriteString("\n\n")

	return sb.String()
}

// ============================================================
// Cobra コマンド
// ============================================================

var listCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("list.short"),
	Long:  i18n.T("list.long"),
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("shell-init", false, i18n.T("list.flag.shell_init"))
	listCmd.Flags().Bool("shell-init-bash", false, i18n.T("list.flag.shell_init_bash"))
}

func runList(cobraCmd *cobra.Command, args []string) error {
	if ok, _ := cobraCmd.Flags().GetBool("shell-init"); ok {
		fmt.Print(shellInitZsh)
		return nil
	}
	if ok, _ := cobraCmd.Flags().GetBool("shell-init-bash"); ok {
		fmt.Print(shellInitBash)
		return nil
	}

	repo, err := repository.NewCommandRepository()
	if err != nil {
		return err
	}

	commands, err := repo.ListCommands()
	if err != nil {
		return err
	}

	if len(commands) == 0 {
		fmt.Fprintln(os.Stderr, i18n.T("list.empty"))
		return nil
	}

	categories, err := repo.ListCategories()
	if err != nil {
		return err
	}

	// /dev/tty を直接開くことで、$(clp list) でも UI が正常に表示される
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf(i18n.T("list.err.tty"), err)
	}
	defer tty.Close()

	m := newListModel(commands, categories)
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))

	result, err := p.Run()
	if err != nil {
		return err
	}

	if final, ok := result.(listModel); ok && final.selected != "" {
		fmt.Println(final.selected)
	}
	return nil
}
