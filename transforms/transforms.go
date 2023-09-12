package transforms

import "encoding/json"

type ITransform[T any] interface {
	Transform(model T) map[string]interface{}
	TransformModels(models []T) []map[string]interface{}
	Language() string
	SetLanguage(string)
}

type TransformAbstract[T any] struct {
	language string
}

func GetMany[N any, M any](base ITransform[N], field ITransform[M], models []M) []map[string]any {
	field.SetLanguage(base.Language())
	return field.TransformModels(models)
}

func GetModels[M any](base ITransform[M], models []M) []map[string]any {
	list := make([]map[string]any, len(models))
	for i, model := range models {
		list[i] = base.Transform(model)
	}
	return list
}

func NewTransformAbstract[T any]() *TransformAbstract[T] {
	return &TransformAbstract[T]{}
}

func (t *TransformAbstract[T]) Transform(model T) map[string]any {
	j, _ := json.Marshal(model)
	var result map[string]any
	_ = json.Unmarshal(j, &result)
	return result
}

func (t *TransformAbstract[T]) TransformModels(models []T) []map[string]any {
	list := make([]map[string]any, len(models))
	for i, model := range models {
		list[i] = t.Transform(model)
	}
	return list
}

func (t *TransformAbstract[T]) Language() string {
	return t.language
}

func (t *TransformAbstract[T]) SetLanguage(s string) {
	t.language = s
}
