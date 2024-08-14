package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

const (
	tagValidate  = "validate"
	tagDelimeter = "|"
)

var (
	// Ошибка если на вход подали не структуру.
	errValidationNotStruct = errors.New("value is not a struct")
	// Ошибка в структуре правила валидации (прим: 'max:!23' и 'len:char').
	errValidationInvalidRule = errors.New("invalid validation rule")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result strings.Builder

	result.WriteString("validation errors: ")
	for idx, err := range v {
		result.WriteString(err.Field)
		result.WriteString(": ")
		result.WriteString(err.Err.Error())
		if idx < len(v)-1 {
			result.WriteString(";")
		}
	}

	return result.String()
}

func Validate(v interface{}) error {
	errs := ValidationErrors{}
	var validErr ValidationError

	// Проверка на то, что входной `interface{}` - структура
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		return errValidationNotStruct
	}

	// Вынимаем поля структуры
	for idx := 0; idx < value.Type().NumField(); idx++ {
		fieldValue := value.Field(idx)
		fieldType := value.Type().Field(idx)

		// Если у поля есть тег 'validate'
		if fieldTag, ok := fieldType.Tag.Lookup(tagValidate); ok {
			if err := validate(fieldType, fieldValue, fieldTag); err != nil {
				// Если ошибка не ValidationError, а обычная - то возвращаем её сразу
				if !errors.As(err, &validErr) {
					return err
				}
				// Ошибки типа ValidationError накапливаем
				errs = append(errs, validErr)
			}
		}
	}

	// Если errs не пустой - возвращаем его
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validate(fieldType reflect.StructField, fieldValue reflect.Value, fieldTag string) error {
	// Делим строку с валидациями по "|"
	rules := strings.Split(fieldTag, tagDelimeter)

	var err error
	// Валидируем int, string
	switch fieldValue.Kind() { //nolint:exhaustive
	case reflect.Int:
		err = validateInt(fieldValue.Int(), rules)
	case reflect.String:
		err = validateString(fieldValue.String(), rules)
	case reflect.Slice, reflect.Array:
		// Слайсы и массивы проверяем поэлементно
		for i := 0; i < fieldValue.Len(); i++ {
			switch fieldValue.Index(i).Kind() { //nolint:exhaustive
			case reflect.Int:
				err = validateInt(fieldValue.Index(i).Int(), rules)
			case reflect.String:
				err = validateString(fieldValue.Index(i).String(), rules)
			}
		}
	}

	var validErr ValidationError
	if !errors.As(err, &validErr) {
		return err
	}

	return ValidationError{Field: fieldType.Name, Err: validErr.Err}
}
