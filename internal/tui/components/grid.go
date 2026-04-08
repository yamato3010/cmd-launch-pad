package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourname/cmd-launch-pad/internal/models"
)

// RenderGrid はコマンド一覧をグリッド形式で描画する
// commands: 表示するコマンド一覧
// cursor: フォーカス中のインデックス (-1 は追加ボタン)
// cols: 列数
// showAdd: 追加カードを表示するか
func RenderGrid(commands []models.Command, cursor int, cols int, showAdd bool) string {
	items := make([]string, 0, len(commands)+1)
	for i, cmd := range commands {
		items = append(items, RenderCard(cmd, cursor == i))
	}
	if showAdd {
		items = append(items, RenderAddCard(cursor == len(commands)))
	}

	if len(items) == 0 {
		return "コマンドが登録されていません。 n で新規追加できます。"
	}

	rows := []string{}
	for i := 0; i < len(items); i += cols {
		end := i + cols
		if end > len(items) {
			end = len(items)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, items[i:end]...)
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}
