package hw09structvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexpIntValidatorMin = regexp.MustCompile(`^min:\d+$`)
	regexpIntValidatorMax = regexp.MustCompile(`^max:\d+$`)
	regexpIntValidatorIn  = regexp.MustCompile(`^in:\d+(,\d+)*$`)

	errValidationIntMin = errors.New("value is lower then 'min' rule")
	errValidationIntMax = errors.New("value is bigger than 'max' rule")
	errValidationIntIn  = errors.New("value not found in 'in' rule list")
)

// validateInt валидирует число на соответствие списку правил.
func validateInt(value int64, rulesList []string) error { //nolint:gocognit
	errs := []error{}

	for _, rule := range rulesList {
		switch {
		case strings.HasPrefix(rule, "min:"):
			if !regexpIntValidatorMin.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateIntMin(value, rule); err != nil {
				if !errors.Is(err, errValidationIntMin) {
					return err
				}
				errs = append(errs, errValidationIntMin)
			}
		case strings.HasPrefix(rule, "max:"):
			if !regexpIntValidatorMax.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateIntMax(value, rule); err != nil {
				if !errors.Is(err, errValidationIntMax) {
					return err
				}
				errs = append(errs, errValidationIntMax)
			}
		case strings.HasPrefix(rule, "in:"):
			if !regexpIntValidatorIn.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateIntIn(value, rule); err != nil {
				if !errors.Is(err, errValidationIntIn) {
					return err
				}
				errs = append(errs, errValidationIntIn)
			}
		}
	}

	// Если errs не пустой - возвращаем его как ValidationError с ошибками через Join.
	if len(errs) > 0 {
		return ValidationError{Err: errors.Join(errs...)}
	}

	return nil
}

// validateIntMin проверяет число на соответствие правилу 'min'.
func validateIntMin(value int64, rule string) error {
	ruleMin, err := strconv.ParseInt(strings.SplitN(rule, ":", 2)[1], 10, 64)
	if err != nil {
		return err
	}
	if value < ruleMin {
		return errValidationIntMin
	}

	return nil
}

// validateIntMax проверяет число на соответствие правилу 'max'.
func validateIntMax(value int64, rule string) error {
	ruleMax, err := strconv.ParseInt(strings.SplitN(rule, ":", 2)[1], 10, 64)
	if err != nil {
		return err
	}
	if value >= ruleMax {
		return errValidationIntMax
	}

	return nil
}

// validateIntIn проверяет число на соответствие правилу 'in'.
func validateIntIn(value int64, rule string) error {
	ruleInList := strings.SplitN(rule, ":", 2)[1]
	// Сравниваем строки, а не числа
	valueStr := strconv.FormatInt(value, 10)

	for _, ruleInValue := range strings.Split(ruleInList, ",") {
		if valueStr == ruleInValue {
			return nil
		}
	}

	return errValidationIntIn
}
