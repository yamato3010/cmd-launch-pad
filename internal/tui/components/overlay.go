package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yamato3010/cmd-launch-pad/internal/tui/styles"
)

// PlaceOverlay は背景文字列 bg の中央にモーダルコンテンツ modal をオーバーレイして返す。
// screenW, screenH は端末のサイズ。
func PlaceOverlay(bg, modal string, screenW, screenH int) string {
	if screenW == 0 || screenH == 0 {
		return modal
	}

	bgLines := strings.Split(bg, "\n")
	modalLines := strings.Split(modal, "\n")

	// モーダルの実際の幅・高さを計算
	modalW := 0
	for _, line := range modalLines {
		w := lipgloss.Width(line)
		if w > modalW {
			modalW = w
		}
	}
	modalH := len(modalLines)

	// 中央配置の開始位置
	startX := (screenW - modalW) / 2
	startY := (screenH - modalH) / 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}

	// 背景行数が足りない場合は補完
	for len(bgLines) < screenH {
		bgLines = append(bgLines, strings.Repeat(" ", screenW))
	}

	// 各行を合成
	result := make([]string, screenH)
	for y := 0; y < screenH; y++ {
		bgLine := ""
		if y < len(bgLines) {
			bgLine = bgLines[y]
		}

		// モーダル行の範囲内か
		modalY := y - startY
		if modalY < 0 || modalY >= modalH {
			result[y] = bgLine
			continue
		}

		modalLine := modalLines[modalY]
		result[y] = overlayLine(bgLine, modalLine, startX, screenW)
	}

	return strings.Join(result, "\n")
}

// overlayLine は背景行の startX の位置からモーダル行を上書きした文字列を返す。
// ANSI エスケープシーケンスを考慮するため lipgloss を使って処理する。
func overlayLine(bgLine, modalLine string, startX, screenW int) string {
	// 背景行を screenW にパディング
	bgW := lipgloss.Width(bgLine)
	if bgW < screenW {
		bgLine += strings.Repeat(" ", screenW-bgW)
	}

	// startX より左の背景部分を取得（文字単位でスライス）
	left := truncateANSI(bgLine, startX)
	leftW := lipgloss.Width(left)
	// 幅が足りない場合はスペースで補完
	if leftW < startX {
		left += strings.Repeat(" ", startX-leftW)
	}

	// モーダル行の幅
	modalW := lipgloss.Width(modalLine)

	// モーダルより右の背景部分を取得
	rightStart := startX + modalW
	right := ""
	if rightStart < screenW {
		right = skipANSI(bgLine, rightStart)
	}

	return left + modalLine + right
}

// truncateANSI は ANSI コードを含む文字列を表示幅 n でカットして返す。
func truncateANSI(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return lipgloss.NewStyle().MaxWidth(n).Render(s)
}

// skipANSI は表示幅 skip 文字分スキップして残りを返す（簡易実装）。
func skipANSI(s string, skip int) string {
	// lipgloss でスタイルなしの文字列に変換してからスキップ
	plain := stripANSI(s)
	runes := []rune(plain)
	pos := 0
	cur := 0
	for pos < len(runes) && cur < skip {
		cur++
		pos++
	}
	if pos >= len(runes) {
		return ""
	}
	return string(runes[pos:])
}

// stripANSI は ANSI エスケープシーケンスを除去した文字列を返す（簡易版）。
func stripANSI(s string) string {
	return styles.StripANSI(s)
}

// ModalBox はモーダルボックス全体のスタイル
var ModalBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(styles.ColorAccent).
	Padding(1, 2)

// ModalTitle はモーダルのタイトルスタイル
var ModalTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(styles.ColorAccent).
	MarginBottom(1)
