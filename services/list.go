package services

import (
	"fmt"
	"github.com/goasali/toolkit/repositories"
	"github.com/samber/lo"
	"strings"
)

type Pagination struct {
	Page    int
	PerPage int
}

type OrderDirection string

var (
	DESC OrderDirection = "desc"
	ASC  OrderDirection = "asc"
)

type Order struct {
	Column    string
	Direction OrderDirection
}

type orders []Order

func (order Order) string() string {
	return fmt.Sprintf("%s %s", order.Column, order.Direction)
}

func (o orders) string() string {
	ordersStr := lo.Map(o, func(item Order, index int) string {
		return item.string()
	})
	return strings.Join(ordersStr, " ")
}

type ListOptions struct {
	pagination Pagination
	condition  []interface{}
	order      orders
}

type ListOptionsFunc func(options *ListOptions)

func DefaultListOptions() *ListOptions {
	return &ListOptions{
		pagination: Pagination{
			Page:    1,
			PerPage: 100,
		},
		condition: make([]interface{}, 0),
		order:     []Order{{Column: "id", Direction: DESC}},
	}
}

func WithPerPage(perPage int) ListOptionsFunc {
	return func(options *ListOptions) {
		options.pagination.PerPage = perPage
	}
}

func WithPage(page int) ListOptionsFunc {
	return func(options *ListOptions) {
		options.pagination.Page = page
	}
}

func WithPagination(page int, perPage int) ListOptionsFunc {
	return func(options *ListOptions) {
		WithPage(page)(options)
		WithPerPage(perPage)(options)
	}
}

func WithCondition(condition interface{}) ListOptionsFunc {
	return func(options *ListOptions) {
		options.condition = append(options.condition, condition)
	}
}

func WithOrder(order ...Order) ListOptionsFunc {
	return func(options *ListOptions) {
		options.order = order
	}
}

func (as AbstractService[M]) List(models *[]M, optionCallbacks ...ListOptionsFunc) (count int64, err error) {
	option := DefaultListOptions()
	for _, callback := range optionCallbacks {
		callback(option)
	}
	conditionCount := len(option.condition)
	conditions := make([]repositories.ListQueryExecuteFn, conditionCount+3)

	for i, condition := range option.condition {
		conditions[i] = repositories.WithAppendCondition(condition)
	}
	conditions[conditionCount] = repositories.WithPage(option.pagination.Page)
	conditions[conditionCount+1] = repositories.WithPerPage(option.pagination.PerPage)
	conditions[conditionCount+2] = repositories.WithOrder(option.order.string())

	if count, err = as.repository.List(
		models,
		conditions...,
	); err != nil {
		return
	}
	return
}
