package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/goasali/toolkit/errors"
	"github.com/goasali/toolkit/multilingual"
	"github.com/goasali/toolkit/transforms"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	validator2 "gopkg.in/validator.v2"
	"math"
	"net/http"
	"strings"
)

type responseConfig struct {
	c        *gin.Context
	message  string
	other    map[string]interface{}
	status   bool
	httpCode int
}

type OptionFunc func(*responseConfig)

type MessageCount string

//goland:noinspection GoUnusedGlobalVariable
var (
	MessageOne  MessageCount = "one"
	MessageMany MessageCount = "many"
	MessageZero MessageCount = "zero"
	MessageNot  MessageCount = "not"
)

type MessageOption struct {
	Message   string
	Params    map[string]string
	Category  string
	CountName MessageCount
}

func WithMessage(messageOption MessageOption) OptionFunc {
	return func(config *responseConfig) {
		messages := make([]string, 0)
		if messageOption.Category != "" {
			messages = append(messages, messageOption.Category)
		}
		if messageOption.Message != "" {
			messages = append(messages, messageOption.Message)
		}
		if messageOption.CountName != "" {
			messages = append(messages, string(messageOption.CountName))
		}
		id := strings.Join(messages, ".")
		text := multilingual.MessageByRequest(config.c, multilingual.MessageConfig{
			Params:    messageOption.Params,
			MessageId: id,
		})
		if text != "" {
			config.message = text
		} else {
			config.message = id
		}
	}
}

func WithMessageTextOnly(msgId string) OptionFunc {
	return WithMessage(MessageOption{
		Message: msgId,
	})
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

func WithFields(fields map[string]interface{}) OptionFunc {
	return func(config *responseConfig) {
		for key, value := range fields {
			WithFieldKey(key, value)(config)
		}
	}
}

func WithModelList[T any](field string, models []T, count int64, transformer transforms.ITransform[T], pagination Page) OptionFunc {
	total := math.Ceil(float64(count) / float64(pagination.PerPage))
	return func(config *responseConfig) {
		WithFieldKey("count", count)(config)
		WithFieldKey(field, transformer.TransformModels(models))(config)
		WithFieldKey("pagination", map[string]any{
			"page":     pagination.Page,
			"per_page": pagination.PerPage,
			"total":    total,
		})(config)
	}
}

func WithHttpCode(code int) OptionFunc {
	return func(config *responseConfig) {
		config.httpCode = code
	}
}

type Response struct {
}

func (r Response) Send404(c *gin.Context) {
	r.Send(
		c,
		WithStatus(false),
		WithMessageTextOnly("errors.not_found"),
		WithHttpCode(http.StatusNotFound),
	)
}

func (r Response) Send401(c *gin.Context) {
	r.Send(
		c,
		WithHttpCode(http.StatusUnauthorized),
		WithMessageTextOnly("authorization.access_denied"),
	)
}

func (r Response) Send500(c *gin.Context) {
	r.Send(
		c,
		WithStatus(false),
		WithMessageTextOnly("errors.internal_server"),
		WithHttpCode(http.StatusInternalServerError),
	)
}

func (r Response) handleValidationError(c *gin.Context, err validator.ValidationErrors) {
	formErrors := make(map[string]interface{})

	for _, field := range err {
		msgId := "validation." + field.Tag()
		fieldKey := field.Field()
		keyId := fmt.Sprintf("validation.fields.%s", fieldKey)
		fieldName := multilingual.MessageByRequest(c, multilingual.MessageConfig{
			MessageId: keyId,
			Default:   fieldKey,
		})
		formErrors[fieldKey] = multilingual.MessageByRequest(c, multilingual.MessageConfig{
			Params:    map[string]string{"Field": fieldName, "Param": field.Param()},
			MessageId: msgId,
		})
	}

	r.Send(
		c,
		WithMessageTextOnly("validation.error"),
		WithFieldKey("fields", formErrors),
		WithHttpCode(http.StatusBadRequest),
	)
}

// HandleError Handle errors in program with custom response for each error's type
func (r Response) HandleError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := "errors.internal_server"

	if err.Error() == "record not found" {
		r.Send404(c)
		return
	}

	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		r.handleValidationError(c, verr)
		return
	}

	var i18nError pkgErrors.I18nMessageError
	if errors.As(err, &i18nError) {
		msg = i18nError.I18nId
		code = http.StatusBadRequest
	}

	var serviceError pkgErrors.ServiceError
	if errors.As(err, &serviceError) {
		msg = serviceError.Error()
	}

	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		log.Errorf(unmarshalErr.Error())
		msg = "validation.error"
	}

	var validationErr validator2.ErrorMap
	if errors.As(err, &validationErr) {
		msg = validationErr.Error()
		code = http.StatusBadRequest
	}

	if err.Error() == "EOF" {
		msg = "validation.error"
		code = http.StatusBadRequest
	}

	log.Error(err)
	r.Send(c, WithMessageTextOnly(msg), WithHttpCode(code))
}

// SendAccessDenied sends an "access denied" message and aborts the current HTTP context.
func (r Response) SendAccessDenied(c *gin.Context) {
	r.SendMessage(c, http.StatusUnauthorized, "authorization.access_denied")
	c.Abort()
}

func (r Response) SendMessage(c *gin.Context, statusCode int, msg string, fields ...map[string]interface{}) {
	status := false
	if statusCode > 100 && statusCode < 400 {
		status = true
	}

	mergedFields := lo.Assign(fields...)

	r.Send(
		c,
		WithMessageTextOnly(msg),
		WithStatus(status),
		WithHttpCode(statusCode),
		WithFields(mergedFields),
	)
}

func (r Response) SendCreated(c *gin.Context, model any, key string) {
	msg := "created_successfully"
	if key != "" {
		msg = fmt.Sprintf("%s.%s", key, msg)
	}
	msg = fmt.Sprintf("messages.%s", msg)

	r.Send(
		c,
		WithMessageTextOnly(msg),
		WithHttpCode(http.StatusCreated),
		WithStatus(true),
		WithFieldKey("model", model),
	)
}

func (r Response) SendModelUpdated(c *gin.Context, key string) {
	msg := "updated_successfully"
	if key != "" {
		msg = fmt.Sprintf("%s.%s", key, msg)
	}
	msg = fmt.Sprintf("messages.%s", msg)

	r.SendMessage(c, http.StatusOK, msg)
}

func (r Response) SendModelDeleted(c *gin.Context, key string, count int) {
	var msg string
	if count == 1 {
		msg = "deleted_successfully"
	} else {
		msg = "deleted_many_successfully"
	}

	if key != "" {
		msg = fmt.Sprintf("%s.%s", key, msg)
	}
	msg = fmt.Sprintf("messages.%s", msg)

	r.Send(
		c,
		WithMessage(MessageOption{
			Message: msg,
			Params:  map[string]string{"count": fmt.Sprintf("%d", count)},
		}),
		WithHttpCode(http.StatusOK),
		WithStatus(true),
	)
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
	result["status"] = conf.status
	c.JSON(conf.httpCode, result)
}
