package errors

type ServiceError struct {
	i18nId string
}

func NewServiceError(i18nId string) ServiceError {
	return ServiceError{i18nId: i18nId}
}

func (s ServiceError) Error() string {
	return s.i18nId
}
