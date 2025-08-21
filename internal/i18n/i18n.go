package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 支持的语言
const (
	LangZhCN    = "zh-cn"
	LangEnUS    = "en-us"
	DefaultLang = LangZhCN
)

// 全局翻译实例
var globalTranslator *Translator

// Translator 翻译器
type Translator struct {
	translations map[string]map[string]string
	defaultLang  string
}

// NewTranslator 创建新的翻译器
func NewTranslator() *Translator {
	return &Translator{
		translations: make(map[string]map[string]string),
		defaultLang:  DefaultLang,
	}
}

// LoadTranslations 加载翻译文件
func (t *Translator) LoadTranslations(dir string) error {
	// 加载中文翻译
	if err := t.loadLanguageFile(filepath.Join(dir, "zh-cn.json"), LangZhCN); err != nil {
		return fmt.Errorf("加载中文翻译失败: %v", err)
	}

	// 加载英文翻译
	if err := t.loadLanguageFile(filepath.Join(dir, "en-us.json"), LangEnUS); err != nil {
		return fmt.Errorf("加载英文翻译失败: %v", err)
	}

	return nil
}

// loadLanguageFile 加载单个语言文件
func (t *Translator) loadLanguageFile(filename, lang string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	t.translations[lang] = translations
	return nil
}

// Translate 翻译文本
func (t *Translator) Translate(lang, key string, args ...interface{}) string {
	// 如果语言不存在，使用默认语言
	if _, exists := t.translations[lang]; !exists {
		lang = t.defaultLang
	}

	// 获取翻译文本
	if text, exists := t.translations[lang][key]; exists {
		if len(args) > 0 {
			return fmt.Sprintf(text, args...)
		}
		return text
	}

	// 如果当前语言没有找到，尝试默认语言
	if lang != t.defaultLang {
		if text, exists := t.translations[t.defaultLang][key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(text, args...)
			}
			return text
		}
	}

	// 都没找到，返回key本身
	return key
}

// GetLanguageFromContext 从gin上下文获取语言
func GetLanguageFromContext(c *gin.Context) string {
	// 1. 从查询参数获取
	if lang := c.Query("lang"); lang != "" {
		return normalizeLang(lang)
	}

	// 2. 从Header获取
	if lang := c.GetHeader("Accept-Language"); lang != "" {
		return parseAcceptLanguage(lang)
	}

	// 3. 返回默认语言
	return DefaultLang
}

// normalizeLang 标准化语言代码
func normalizeLang(lang string) string {
	lang = strings.ToLower(lang)
	switch lang {
	case "zh", "zh-cn", "chinese":
		return LangZhCN
	case "en", "en-us", "english":
		return LangEnUS
	default:
		return DefaultLang
	}
}

// parseAcceptLanguage 解析Accept-Language头
func parseAcceptLanguage(acceptLang string) string {
	// 简单解析，取第一个语言
	langs := strings.Split(acceptLang, ",")
	if len(langs) > 0 {
		lang := strings.TrimSpace(strings.Split(langs[0], ";")[0])
		return normalizeLang(lang)
	}
	return DefaultLang
}

// Init 初始化国际化
func Init(translationsDir string) error {
	globalTranslator = NewTranslator()
	return globalTranslator.LoadTranslations(translationsDir)
}

// T 全局翻译函数
func T(c *gin.Context, key string, args ...interface{}) string {
	if globalTranslator == nil {
		return key
	}
	lang := GetLanguageFromContext(c)
	return globalTranslator.Translate(lang, key, args...)
}

// TWithLang 指定语言翻译
func TWithLang(lang, key string, args ...interface{}) string {
	if globalTranslator == nil {
		return key
	}
	return globalTranslator.Translate(lang, key, args...)
}
