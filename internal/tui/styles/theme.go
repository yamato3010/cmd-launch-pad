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

// カラーパレット (Init でテーマに応じて設定される)
var (
	ColorBg         lipgloss.TerminalColor
	ColorBgAlt      lipgloss.TerminalColor
	ColorBorder     lipgloss.TerminalColor
	ColorBorderFoc  lipgloss.TerminalColor
	ColorText       lipgloss.TerminalColor
	ColorTextAlt    lipgloss.TerminalColor
	ColorTextDim    lipgloss.TerminalColor
	ColorTextDimmer lipgloss.TerminalColor
	ColorAccent     lipgloss.TerminalColor
	ColorGreen      lipgloss.TerminalColor
	ColorYellow     lipgloss.TerminalColor
	ColorRed        lipgloss.TerminalColor
	ColorCyan       lipgloss.TerminalColor
)

// init はデフォルトテーマで初期化する。
// 設定ファイルを読み込んだ後に Init が呼ばれて上書きされる。
func init() {
	Init("")
}

// Init はテーマ名に応じてカラーパレットとスタイルを初期化する。
// "ansi" を指定するとターミナル側の ANSI 配色に従う。それ以外は固定配色 (Tokyo Night 風)。
func Init(theme string) {
	if theme == "ansi" {
		setANSIPalette()
	} else {
		setFixedPalette()
	}
	buildStyles()
}

// setFixedPalette は固定配色 (Tokyo Night 風) を設定する
func setFixedPalette() {
	ColorBg = lipgloss.Color("#1a1b26")
	ColorBgAlt = lipgloss.Color("#16161e")
	ColorBorder = lipgloss.Color("#3b4261")
	ColorBorderFoc = lipgloss.Color("#7aa2f7")
	ColorText = lipgloss.Color("#c0caf5")
	ColorTextAlt = lipgloss.Color("#a9b1d6")
	ColorTextDim = lipgloss.Color("#565f89")
	ColorTextDimmer = lipgloss.Color("#414868")
	ColorAccent = lipgloss.Color("#7aa2f7")
	ColorGreen = lipgloss.Color("#9ece6a")
	ColorYellow = lipgloss.Color("#e0af68")
	ColorRed = lipgloss.Color("#f7768e")
	ColorCyan = lipgloss.Color("#2ac3de")
}

// setANSIPalette は ANSI パレットのインデックスを設定する。
// 実際の色はターミナルのカラースキームが決める。
// 本文の文字色と背景色は NoColor にして、ターミナルの既定色を透過させる。
func setANSIPalette() {
	ColorBg = lipgloss.NoColor{}
	ColorBgAlt = lipgloss.NoColor{}
	ColorBorder = lipgloss.Color("8")     // bright black
	ColorBorderFoc = lipgloss.Color("12") // bright blue
	ColorText = lipgloss.NoColor{}
	ColorTextAlt = lipgloss.NoColor{}
	ColorTextDim = lipgloss.Color("8")
	ColorTextDimmer = lipgloss.Color("8")
	ColorAccent = lipgloss.Color("4") // blue
	ColorGreen = lipgloss.Color("2")
	ColorYellow = lipgloss.Color("3")
	ColorRed = lipgloss.Color("1")
	ColorCyan = lipgloss.Color("6")
}

// スタイル (buildStyles で構築される)
var (
	// AppTitle はアプリタイトルのスタイル
	AppTitle lipgloss.Style
	// TabActive はアクティブなタブのスタイル
	TabActive lipgloss.Style
	// TabInactive は非アクティブなタブのスタイル
	TabInactive lipgloss.Style
	// CardNormal は通常状態のカードスタイル
	CardNormal lipgloss.Style
	// CardFocused はフォーカス状態のカードスタイル
	CardFocused lipgloss.Style
	// CardTitle はカードタイトルのスタイル
	CardTitle lipgloss.Style
	// CardDesc はカード説明文のスタイル
	CardDesc lipgloss.Style
	// StatusBar はステータスバーのスタイル
	StatusBar lipgloss.Style
	// StatusBarKey はステータスバーのキーのスタイル
	StatusBarKey lipgloss.Style
	// Header はヘッダーのスタイル
	Header lipgloss.Style
	// ErrorStyle はエラーメッセージのスタイル
	ErrorStyle lipgloss.Style
	// SuccessStyle は成功メッセージのスタイル
	SuccessStyle lipgloss.Style
	// HelpKey はヘルプ画面のキー表示スタイル
	HelpKey lipgloss.Style
	// HelpDesc はヘルプ画面の説明文スタイル
	HelpDesc lipgloss.Style
	// InputLabel は入力フォームのラベルスタイル
	InputLabel lipgloss.Style
	// DialogBox はダイアログボックスのスタイル
	DialogBox lipgloss.Style
	// DescPanel はコマンド説明パネルのスタイル
	DescPanel lipgloss.Style
	// DescPanelTitle はコマンド説明パネルのタイトルスタイル
	DescPanelTitle lipgloss.Style
	// DescPanelText はコマンド説明パネルの説明文スタイル
	DescPanelText lipgloss.Style
)

// buildStyles は現在のカラーパレットからスタイルを構築する
func buildStyles() {
	AppTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent).
		Padding(0, 1)

	TabActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent).
		Underline(true).
		Padding(0, 1)

	TabInactive = lipgloss.NewStyle().
		Foreground(ColorTextDim).
		Padding(0, 1)

	CardNormal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1).
		Width(14).
		Height(5)

	CardFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderFoc).
		Padding(0, 1).
		Width(14).
		Height(5)

	CardTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorText).
		MaxWidth(12)

	CardDesc = lipgloss.NewStyle().
		Foreground(ColorTextDim).
		MaxWidth(12)

	StatusBar = lipgloss.NewStyle().
		Foreground(ColorTextDim).
		Background(ColorBgAlt).
		Padding(0, 1)

	StatusBarKey = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Background(ColorBgAlt).
		Bold(true)

	Header = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorBgAlt).
		Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)

	HelpKey = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Width(12)

	HelpDesc = lipgloss.NewStyle().
		Foreground(ColorText)

	InputLabel = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Width(14)

	DialogBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Padding(1, 2)

	DescPanel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	DescPanelTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	DescPanelText = lipgloss.NewStyle().
		Foreground(ColorText)
}
