package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/errors"
	"github.com/goasali/toolkit/multilingual"
	"github.com/goasali/toolkit/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

type ResourceController interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Create(c *gin.Context)
	Store(c *gin.Context)
	Edit(c *gin.Context)
	Update(c *gin.Context)
	Destroy(c *gin.Context)
	Delete(c *gin.Context)
}

type Controllers struct {
	Response Response
}

func New() *Controllers {
	return &Controllers{
		Response: Response{},
	}
}

func (ctrl Controllers) Create(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) List(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Index(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Edit(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Show(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Store(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Update(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Delete(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) Destroy(c *gin.Context) {
	ctrl.Response.Send404(c)
}

func (ctrl Controllers) AcceptLanguage(c *gin.Context) string {
	return c.GetHeader("Accept-Language")
}

func (ctrl Controllers) LoadLocalize(c *gin.Context) *i18n.Localizer {
	accept := ctrl.AcceptLanguage(c)
	bundle := multilingual.Bundle()
	return i18n.NewLocalizer(bundle, accept)
}

func (ctrl Controllers) AccessDenied(c *gin.Context) {
	ctrl.Response.Send401(c)
}

func (ctrl Controllers) GetIdFromParam(c *gin.Context) int {
	id := c.Param("id")
	re, err := strconv.Atoi(id)
	if err != nil {
		ctrl.Response.Send(
			c,
			WithMessageTextOnly("validation.params_id_not_valid"),
			WithStatus(false),
			WithHttpCode(400),
		)
		return 0
	}
	return re
}

func (ctrl Controllers) GetFieldFromParam(c *gin.Context, field string) int {
	value := c.Param(field)
	re, err := strconv.Atoi(value)
	if err != nil {
		panic(errors.NewI18n("params_not_valid_for_field", map[string]interface{}{
			"Field": field,
		}))
	}
	return re
}

func (ctrl Controllers) GetIdFromBody(c *gin.Context) []uint {
	body := BodyId{}

	if err := c.ShouldBind(&body); err != nil {
		panic(err)
	}

	return body.Id
}

func (ctrl Controllers) GetFilter(c *gin.Context) []WhereCondition {
	conditionJson := c.Query("filter")

	var condition []WhereCondition
	_ = json.Unmarshal([]byte(conditionJson), &condition)

	return condition
}

func (ctrl Controllers) Pagination(c *gin.Context) Page {
	pagination := Page{
		Page:    1,
		PerPage: 100,
	}

	if perPage, _ := strconv.Atoi(c.Query("per_page")); perPage > 0 {
		pagination.PerPage = perPage
	}
	if page, _ := strconv.Atoi(c.Query("page")); page > 0 {
		pagination.Page = page
	}

	return pagination
}

func (ctrl Controllers) Order(c *gin.Context) []services.Order {
	// Order example: ?order=id&order=name,asc&order=wtf,desc
	orders := c.QueryArray("order")

	return lo.Map(orders, func(order string, _ int) services.Order {
		orderSplit := strings.Split(order, ",")
		var direction string
		if len(orderSplit) > 1 {
			direction = orderSplit[1]
		} else {
			direction = "desc"
		}
		return services.Order{
			Column:    orderSplit[0],
			Direction: services.OrderDirection(direction),
		}
	})
}
