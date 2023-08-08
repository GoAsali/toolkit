package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/multilingual"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"strconv"
)

type Controllers struct {
	Response Response
}

func New() *Controllers {
	return &Controllers{
		Response: Response{},
	}
}

func (ctrl Controllers) LoadLocalize(c *gin.Context) *i18n.Localizer {
	accept := c.GetHeader("Accept-Language")
	bundle := multilingual.Bundle()
	return i18n.NewLocalizer(bundle, accept)
}

func (ctrl Controllers) AccessDenied(c *gin.Context) {
	ctrl.Response.Send(c, WithHttpCode(http.StatusUnauthorized), WithI18n("authorization.access_denied"))
}

func (ctrl Controllers) GetIdFromParam(c *gin.Context) int {
	id := c.Param("id")
	re, err := strconv.Atoi(id)
	if err != nil {
		ctrl.Response.Send(c, WithHttpCode(http.StatusUnauthorized), WithI18n("validation.params_id_not_valid"))
		return 0
	}
	return re
}

func (ctrl Controllers) GetIdFromBody(c *gin.Context) []uint {
	body := BodyId{}

	if err := c.ShouldBind(&body); err != nil {
		ctrl.Response.HandleGinError(err, c)
		return make([]uint, 0)
	}

	return body.Id
}
