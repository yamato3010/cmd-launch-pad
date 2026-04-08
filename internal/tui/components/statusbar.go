package components

import (
	"strings"

	"github.com/yourname/cmd-launch-pad/internal/tui/styles"
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
