package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/multilingual"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strconv"
)

type responseConfig struct {
	c        *gin.Context
	message  string
	other    map[string]interface{}
	status   bool
	httpCode int
}

type OptionFunc func(*responseConfig)

func WithMessage(msg string) OptionFunc {
	return func(config *responseConfig) {
		config.message = msg
	}
}

func i18nMessageOrDefault(c *gin.Context, msgId string, defaultValue string, params ...map[string]string) string {
	accept := c.GetHeader("Accept-Language")
	loc := i18n.NewLocalizer(multilingual.Bundle(), accept)
	value := loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: msgId,
		},
		TemplateData: params,
	})
	if value == "" {
		return defaultValue
	}
	return value
}

func i18nMessage(c *gin.Context, msgId string, params ...map[string]string) string {
	return i18nMessageOrDefault(c, msgId, msgId, params...)
}

func WithI18n(msgId string, params ...map[string]string) OptionFunc {
	return func(config *responseConfig) {
		config.message = i18nMessage(config.c, msgId, params...)
	}
}

func WithStatus(status bool) OptionFunc {
	return func(config *responseConfig) {
		config.status = status
	}
}

func WithFieldKey(key string, value interface{}) OptionFunc {
	return func(config *responseConfig) {
		config.other[key] = value
	}
}

func WithHttpCode(code int) OptionFunc {
	return func(config *responseConfig) {
		config.httpCode = code
	}
}

type Response struct {
}

func (r Response) Send(c *gin.Context, options ...OptionFunc) {
	conf := responseConfig{
		c:        c,
		message:  "",
		other:    make(map[string]interface{}),
		status:   false,
		httpCode: 0,
	}
	for _, option := range options {
		option(&conf)
	}
	result := conf.other
	result["message"] = conf.message
	result["status"] = strconv.FormatBool(conf.status)
	c.JSON(conf.httpCode, result)
}
