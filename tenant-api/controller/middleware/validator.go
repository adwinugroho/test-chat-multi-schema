package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	var errorArr []string
	err := cv.Validator.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorArr = append(
				errorArr,
				fmt.Sprintf("%v field doesn't satisfy the %v constraint", err.Field(), err.Tag()),
			)
		}
	}

	if len(errorArr) > 0 {
		return errors.New(strings.Join(errorArr, "\n"))
	}
	return nil
}
