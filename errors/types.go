package errors

type I18nMessageError struct {
	error
	I18nId string
	Values map[string]interface{}
}

func NewI18nJustText(i18nId string) I18nMessageError {
	return NewI18n(i18nId, map[string]interface{}{})
}

func NewI18n(i18nId string, values map[string]interface{}) I18nMessageError {
	return I18nMessageError{
		I18nId: i18nId,
		Values: values,
	}
}
