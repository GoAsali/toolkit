package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (r Response) handleValidationError(err validator.ValidationErrors, c *gin.Context) {
	message := i18nMessage(c, "validation.error")
	code := http.StatusBadRequest

	formErrors := map[string]string{}
	for _, field := range err {
		id := "validation." + field.Tag()
		fieldKey := field.Field()
		keyId := fmt.Sprintf("validation.fields.%s", fieldKey)
		fieldName := i18nMessageOrDefault(c, keyId, fieldKey)
		params := map[string]string{
			"Param": field.Param(),
			"Field": fieldName,
		}
		formErrors[fieldKey] = i18nMessageOrDefault(c, id, fieldName, params)
	}

	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
		"status":  false,
		"fields":  formErrors,
	})
}

func (r Response) HandleGinError(err error, c *gin.Context) {
	message := err.Error()
	code := http.StatusInternalServerError

	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		r.handleValidationError(verr, c)
		return
	}

	if err.Error() == "EOF" {
		r.Send(c, WithI18n("validation.error"), WithHttpCode(http.StatusBadRequest))
		return
	}

	r.Send(c, WithMessage(message), WithHttpCode(code))
}
