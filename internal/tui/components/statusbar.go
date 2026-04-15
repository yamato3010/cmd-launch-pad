package components

import (
	"strings"

	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// KeyBinding はキーとその説明のペア
type KeyBinding struct {
	Key  string
	Desc string
}

// RenderStatusBar はキーバインド一覧をステータスバーとして描画する
func RenderStatusBar(bindings []KeyBinding, width int) string {
	parts := make([]string, 0, len(bindings))
	for _, b := range bindings {
		key := styles.StatusBarKey.Render(b.Key)
		desc := styles.StatusBar.Render(":" + b.Desc)
		parts = append(parts, key+desc)
	}
	bar := strings.Join(parts, styles.StatusBar.Render("  "))
	return styles.StatusBar.Copy().Width(width).Render(bar)
}

// RenderDescPanel はフォーカス中のコマンドの説明パネルを描画する
// name: コマンド名, description: 説明文, width: パネルの幅
func RenderDescPanel(name, description string, width int) string {
	if name == "" && description == "" {
		// 追加カードにフォーカスしている場合などは空パネルを返す
		inner := styles.DescPanelText.Render(i18n.T("desc.add_hint"))
		innerWidth := width - 4 // ボーダー + パディング分
		if innerWidth < 1 {
			innerWidth = 1
		}
		return styles.DescPanel.Copy().Width(innerWidth).Render(inner)
	}

	title := styles.DescPanelTitle.Render("▶ " + name)
	desc := description
	if desc == "" {
		desc = "—"
	}
	body := styles.DescPanelText.Render(desc)
	inner := title + "  " + body

	innerWidth := width - 4 // ボーダー(2) + パディング(2)
	if innerWidth < 1 {
		innerWidth = 1
	}
	return styles.DescPanel.Copy().Width(innerWidth).Render(inner)
}
