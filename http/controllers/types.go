package controllers

type BodyId struct {
	Id []uint `json:"id" binding:"required,min=1" form:"id[]"`
}

type WhereCondition struct {
	Field     string `json:"name"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Page struct {
	PerPage int
	Page    int
}
