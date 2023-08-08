package validations

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strings"
)

type DatabaseValidation struct {
	db *gorm.DB
}

func AddDatabase(db *gorm.DB) error {
	dv := DatabaseValidation{db}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err error
		err = v.RegisterValidation("exists", dv.Exists)
		if err != nil {
			return err
		}
		err = v.RegisterValidation("unique", dv.NotExists)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dv *DatabaseValidation) Exists(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	param := fl.Param()
	field := "id"
	if seps := strings.Split(param, "."); len(seps) > 1 {
		param = seps[0]
		field = seps[1]
	}

	var count int64 = 0
	dv.db.Table(param).Where(field+" = ?", value).Limit(1).Count(&count)

	return count != 0
}

func (dv *DatabaseValidation) NotExists(fl validator.FieldLevel) bool {
	return !dv.Exists(fl)
}
