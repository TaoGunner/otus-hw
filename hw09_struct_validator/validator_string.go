package hw09structvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexpStringValidatorLen    = regexp.MustCompile(`^len:\d+$`)
	regexpStringValidatorRegexp = regexp.MustCompile(`^regexp:.+$`)
	regexpStringValidatorIn     = regexp.MustCompile(`^in:\w+(,\w+)*$`)

	errValidationStringLen    = errors.New("value length not equal with 'len' rule")
	errValidationStringRegexp = errors.New("value not match with 'regexp' rule")
	errValidationStringIn     = errors.New("value not found in 'in' rule list")
)

// validateString валидирует строку на соответствие списку правил.
func validateString(value string, rulesList []string) error { //nolint:gocognit
	errs := []error{}

	for _, rule := range rulesList {
		switch strings.SplitN(rule, ":", 2)[0] {
		case "len":
			if !regexpStringValidatorLen.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateStringLen(value, rule); err != nil {
				// ошибки валидации - накапливаем, остальные сразу же возвращаем
				if !errors.Is(err, errValidationStringLen) {
					return err
				}
				errs = append(errs, err)
			}
		case "regexp":
			if !regexpStringValidatorRegexp.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateStringRegexp(value, rule); err != nil {
				if !errors.Is(err, errValidationStringRegexp) {
					return err
				}
				errs = append(errs, err)
			}
		case "in":
			if !regexpStringValidatorIn.MatchString(rule) {
				return errValidationInvalidRule
			}
			if err := validateStringIn(value, rule); err != nil {
				if !errors.Is(err, errValidationStringIn) {
					return err
				}
				errs = append(errs, err)
			}
		}
	}

	// Если errs не пустой - возвращаем его как ValidationError с ошибками через Join
	if len(errs) > 0 {
		return ValidationError{Err: errors.Join(errs...)}
	}

	return nil
}

// validateStringLen валидирует строку на соответствие правилу 'len'.
func validateStringLen(value string, rule string) error {
	ruleLen, err := strconv.Atoi(strings.SplitN(rule, ":", 2)[1])
	if err != nil {
		return err
	}

	if len(value) != ruleLen {
		return errValidationStringLen
	}

	return nil
}

// validateStringLen валидирует строку на соответствие регулярному выражению в правиле 'regexp'.
func validateStringRegexp(value string, rule string) error {
	regexpRule, err := regexp.Compile(strings.SplitN(rule, ":", 2)[1])
	if err != nil {
		return err
	}

	if !regexpRule.MatchString(value) {
		return errValidationStringRegexp
	}

	return nil
}

// validateStringIn валидирует строку на соответствие правилу 'in'.
func validateStringIn(value string, rule string) error {
	ruleInList := strings.SplitN(rule, ":", 2)[1]

	for _, ruleInValue := range strings.Split(ruleInList, ",") {
		if value == ruleInValue {
			return nil
		}
	}

	return errValidationStringIn
}
