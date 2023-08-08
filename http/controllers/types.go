package controllers

type BodyId struct {
	Id []uint `json:"id" binding:"required,min=1"`
}
