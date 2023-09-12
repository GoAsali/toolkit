package multilingual

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type MessageConfig struct {
	Params    map[string]string
	MessageId string
	Default   string
}

func MessageByLanguage(lang string, config MessageConfig) string {
	loc := i18n.NewLocalizer(Bundle(), lang)
	text := loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: config.MessageId,
		},
		TemplateData: config.Params,
	})
	if text == "" && config.Default != "" {
		return config.Default
	}
	return text
}

func MessageByLanguageTextId(lang string, messageId string) string {
	return MessageByLanguage(lang, MessageConfig{MessageId: messageId})
}

func MessageByRequest(c *gin.Context, config MessageConfig) string {
	accept := c.GetHeader("Accept-Language")
	return MessageByLanguage(accept, config)
}

func MessageByRequestTextOnly(c *gin.Context, messageId string) string {
	return MessageByRequest(c, MessageConfig{MessageId: messageId})
}
