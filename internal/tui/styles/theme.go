package styles

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// ansiEscapeRe は ANSI エスケープシーケンスにマッチする正規表現
var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// StripANSI は ANSI エスケープシーケンスを除去した文字列を返す
func StripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

// カラーパレット (Tokyo Night 風)
var (
	ColorBg        = lipgloss.Color("#1a1b26")
	ColorBgAlt     = lipgloss.Color("#16161e")
	ColorBorder    = lipgloss.Color("#3b4261")
	ColorBorderFoc = lipgloss.Color("#7aa2f7")
	ColorText      = lipgloss.Color("#c0caf5")
	ColorTextDim   = lipgloss.Color("#565f89")
	ColorAccent    = lipgloss.Color("#7aa2f7")
	ColorGreen     = lipgloss.Color("#9ece6a")
	ColorYellow    = lipgloss.Color("#e0af68")
	ColorRed       = lipgloss.Color("#f7768e")
	ColorCyan      = lipgloss.Color("#2ac3de")
)

// AppTitle はアプリタイトルのスタイル
var AppTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorAccent).
	Padding(0, 1)

// TabActive はアクティブなタブのスタイル
var TabActive = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorAccent).
	Underline(true).
	Padding(0, 1)

// TabInactive は非アクティブなタブのスタイル
var TabInactive = lipgloss.NewStyle().
	Foreground(ColorTextDim).
	Padding(0, 1)

// CardNormal は通常状態のカードスタイル
var CardNormal = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorBorder).
	Padding(0, 1).
	Width(14).
	Height(5)

// CardFocused はフォーカス状態のカードスタイル
var CardFocused = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorBorderFoc).
	Padding(0, 1).
	Width(14).
	Height(5)

// CardTitle はカードタイトルのスタイル
var CardTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorText).
	MaxWidth(12)

// CardDesc はカード説明文のスタイル
var CardDesc = lipgloss.NewStyle().
	Foreground(ColorTextDim).
	MaxWidth(12)

// StatusBar はステータスバーのスタイル
var StatusBar = lipgloss.NewStyle().
	Foreground(ColorTextDim).
	Background(ColorBgAlt).
	Padding(0, 1)

// StatusBarKey はステータスバーのキーのスタイル
var StatusBarKey = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Background(ColorBgAlt).
	Bold(true)

// Header はヘッダーのスタイル
var Header = lipgloss.NewStyle().
	Foreground(ColorText).
	Background(ColorBgAlt).
	Padding(0, 1)

// ErrorStyle はエラーメッセージのスタイル
var ErrorStyle = lipgloss.NewStyle().
	Foreground(ColorRed).
	Bold(true)

// SuccessStyle は成功メッセージのスタイル
var SuccessStyle = lipgloss.NewStyle().
	Foreground(ColorGreen).
	Bold(true)

// HelpKey はヘルプ画面のキー表示スタイル
var HelpKey = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Bold(true).
	Width(12)

// HelpDesc はヘルプ画面の説明文スタイル
var HelpDesc = lipgloss.NewStyle().
	Foreground(ColorText)

// InputLabel は入力フォームのラベルスタイル
var InputLabel = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Bold(true).
	Width(14)

// DialogBox はダイアログボックスのスタイル
var DialogBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorAccent).
	Padding(1, 2)

// DescPanel はコマンド説明パネルのスタイル
var DescPanel = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorBorder).
	Padding(0, 1)

// DescPanelTitle はコマンド説明パネルのタイトルスタイル
var DescPanelTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorAccent)

// DescPanelText はコマンド説明パネルの説明文スタイル
var DescPanelText = lipgloss.NewStyle().
	Foreground(ColorText)
