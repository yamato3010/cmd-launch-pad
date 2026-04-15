// Package i18n は多言語対応のためのシンプルな翻訳機能を提供します。
// サポート言語: 英語 (en), 日本語 (ja)
package i18n

import (
	"os"
	"strings"
	"sync"
)

// Lang は言語コード
type Lang string

const (
	LangJa Lang = "ja"
	LangEn Lang = "en"
)

var (
	mu      sync.RWMutex
	current = LangEn
)

// SetLang はアクティブな言語を設定する
func SetLang(l Lang) {
	mu.Lock()
	defer mu.Unlock()
	switch l {
	case LangJa, LangEn:
		current = l
	default:
		current = LangEn
	}
}

// GetLang は現在の言語を返す
func GetLang() Lang {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

// T はキーに対応する翻訳文字列を返す
func T(key string) string {
	mu.RLock()
	lang := current
	mu.RUnlock()

	switch lang {
	case LangJa:
		if s, ok := messagesJa[key]; ok {
			return s
		}
	case LangEn:
		if s, ok := messagesEn[key]; ok {
			return s
		}
	}
	// フォールバック: キーをそのまま返す
	return key
}

// DetectLang は設定値→環境変数の優先順位で言語を自動検出する。
// cfgLang が空でなければその値を優先する。
func DetectLang(cfgLang string) Lang {
	if cfgLang != "" {
		return parseLang(cfgLang)
	}
	// 環境変数から検出
	for _, envKey := range []string{"LANGUAGE", "LANG", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(envKey); val != "" {
			lang := parseLang(val)
			if lang != LangEn {
				return lang
			}
			// en が明示されていればそれを使う
			if strings.HasPrefix(strings.ToLower(val), "en") {
				return LangEn
			}
		}
	}
	return LangEn
}

func parseLang(s string) Lang {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "ja") {
		return LangJa
	}
	return LangEn
}
