package repositories

type Operator string

var (
	AndOperator Operator = "and"
)

type ListQueryOption struct {
	operator   Operator
	conditions []interface{}
	page       int
	perPage    int
	order      string
}

type ListQueryExecuteFn func(option *ListQueryOption)

func DefaultListQuery() *ListQueryOption {
	return &ListQueryOption{
		operator:   AndOperator,
		page:       1,
		perPage:    100,
		conditions: make([]interface{}, 0),
		order:      "id desc",
	}
}

func WithPage(page int) ListQueryExecuteFn {
	return func(option *ListQueryOption) {
		option.page = page
	}
}

func WithAppendCondition(condition interface{}) ListQueryExecuteFn {
	return func(option *ListQueryOption) {
		option.conditions = append(option.conditions, condition)
	}
}

func WithPerPage(perPage int) ListQueryExecuteFn {
	return func(option *ListQueryOption) {
		option.perPage = perPage
	}
}

func WithOrder(order string) ListQueryExecuteFn {
	return func(option *ListQueryOption) {
		option.order = order
	}
}
