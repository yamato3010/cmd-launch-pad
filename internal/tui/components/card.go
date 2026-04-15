package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// RenderCard はコマンドカードを描画して文字列を返す
func RenderCard(cmd models.Command, focused bool) string {
	icon := cmd.Icon
	if icon == "" {
		icon = "⚡"
	}

	title := styles.CardTitle.Render(truncate(cmd.Name, 12))

	content := fmt.Sprintf("%s\n%s", icon, title)

	if focused {
		return styles.CardFocused.Render(content)
	}
	return styles.CardNormal.Render(content)
}

// RenderAddCard は「新規追加」カードを描画して文字列を返す
func RenderAddCard(focused bool) string {
	content := fmt.Sprintf("%s\n%s\n%s",
		"➕",
		styles.CardTitle.Render(i18n.T("launcher.key.new")),
		styles.CardDesc.Render(""),
	)
	if focused {
		return styles.CardFocused.Copy().
			BorderForeground(lipgloss.Color("#9ece6a")).
			Render(content)
	}
	return styles.CardNormal.Copy().
		BorderForeground(lipgloss.Color("#3b4261")).
		Render(content)
}

// truncate は文字列を最大長で切り詰める
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
