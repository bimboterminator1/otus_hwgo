package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SystemError struct {
	Err error
}

func (e SystemError) Error() string {
	return fmt.Sprintf("system error: %v", e.Err)
}

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", ve.Field, ve.Err)
}

func (ve ValidationError) Is(target error) bool {
	t, ok := target.(ValidationError)
	if !ok {
		return false
	}
	return ve.Field == t.Field && errors.Is(ve.Err, t.Err)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, e := range ve {
		builder.WriteString(e.Field + ": " + e.Err.Error() + "\n")
	}
	return builder.String()
}

func (ve ValidationErrors) Is(target error) bool {
	if target == nil && len(ve) == 0 {
		return true
	}

	targetVE, ok := target.(ValidationErrors)
	if !ok {
		return false
	}

	if len(ve) != len(targetVE) {
		return false
	}

	for i := range ve {
		if ve[i].Field != targetVE[i].Field || ve[i].Err != targetVE[i].Err {
			return false
		}
	}

	return true
}

// Predefined system errors.
var (
	ErrSysNotAStruct          = SystemError{errors.New("not a struct")}
	ErrSysUnsupportedType     = SystemError{errors.New("unsupported type")}
	ErrSysUnsupportedSlice    = SystemError{errors.New("slice of unsupported type")}
	ErrSysInvalidRule         = SystemError{errors.New("invalid rule")}
	ErrSysCantConvertLenValue = SystemError{errors.New("can't convert len value")}
	ErrSysCantConvertMaxValue = SystemError{errors.New("can't convert max value")}
	ErrSysCantConvertMinValue = SystemError{errors.New("can't convert min value")}
	ErrSysRegexpCompile       = SystemError{errors.New("regexp compile failed")}
)

// Predefined validation errors.
var (
	ErrValueIsLessThanMinValue = errors.New("value is less than min value")
	ErrValueIsMoreThanMaxValue = errors.New("value is more than max value")
	ErrValueNotInList          = errors.New("value is not in the list")
	ErrStringLengthMismatch    = errors.New("string length mismatch")
	ErrRegexpMatchFailed       = errors.New("regexp match failed")
)

// Validate validates the struct fields based on the `validate` tag.
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrSysNotAStruct
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct && validateTag == "nested" {
			if err := Validate(fieldValue.Interface()); err != nil {
				var verr ValidationErrors
				if errors.As(err, &verr) {
					validationErrors = append(validationErrors, verr...)
				} else {
					return err
				}
			}
			continue
		}
		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.String:
			if err := validateField(field.Name, fieldValue.String(), validateTag, &validationErrors); err != nil {
				return err
			}
		case reflect.Int:
			if err := validateField(field.Name, int(fieldValue.Int()), validateTag, &validationErrors); err != nil {
				return err
			}
		case reflect.Slice:
			if err := validateSliceField(field.Name, fieldValue, validateTag, &validationErrors); err != nil {
				return err
			}
		default:
			return ErrSysUnsupportedType
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

//nolint:gocognit
func validateField[T comparable](fieldName string,
	value T,
	validateTag string,
	errs *ValidationErrors,
) error {
	rules := strings.Split(validateTag, "|")
	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return ErrSysInvalidRule
		}
		key, val := parts[0], parts[1]

		switch key {
		case "len":
			if str, ok := any(value).(string); ok {
				length, err := strconv.Atoi(val)
				if err != nil {
					return ErrSysCantConvertLenValue
				}
				if len(str) != length {
					*errs = append(*errs, ValidationError{Field: fieldName, Err: ErrStringLengthMismatch})
				}
			} else {
				return ErrSysUnsupportedType
			}
		case "regexp":
			if str, ok := any(value).(string); ok {
				re, err := regexp.Compile(val)
				if err != nil {
					return ErrSysRegexpCompile
				}
				if !re.MatchString(str) {
					*errs = append(*errs, ValidationError{Field: fieldName, Err: ErrRegexpMatchFailed})
				}
			} else {
				return ErrSysUnsupportedType
			}
		case "in":
			options := strings.Split(val, ",")
			found := false
			switch v := any(value).(type) {
			case int:
				for _, opt := range options {
					optInt, err := strconv.Atoi(opt)
					if err != nil {
						return ErrSysUnsupportedType
					}
					if v == optInt {
						found = true
						break
					}
				}
			case string:
				for _, opt := range options {
					// Trim spaces for string comparison
					opt = strings.TrimSpace(opt)
					if v == opt {
						found = true
						break
					}
				}
			default:
				return ErrSysUnsupportedType
			}
			if !found {
				*errs = append(*errs, ValidationError{Field: fieldName, Err: ErrValueNotInList})
			}
		case "min":
			if num, ok := any(value).(int); ok {
				minim, err := strconv.Atoi(val)
				if err != nil {
					return ErrSysCantConvertMinValue
				}
				if num < minim {
					*errs = append(*errs, ValidationError{Field: fieldName, Err: ErrValueIsLessThanMinValue})
				}
			} else {
				return ErrSysUnsupportedType
			}
		case "max":
			if num, ok := any(value).(int); ok {
				maxim, err := strconv.Atoi(val)
				if err != nil {
					return ErrSysCantConvertMaxValue
				}
				if num > maxim {
					*errs = append(*errs, ValidationError{Field: fieldName, Err: ErrValueIsMoreThanMaxValue})
				}
			} else {
				return ErrSysUnsupportedType
			}
		default:
			return ErrSysInvalidRule
		}
	}
	return nil
}

func validateSliceField(fieldName string, value reflect.Value, validateTag string, errs *ValidationErrors) error {
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		//nolint:exhaustive
		switch elem.Kind() {
		case reflect.String:
			if err := validateField(fieldName, elem.String(), validateTag, errs); err != nil {
				return err
			}
		case reflect.Int:
			if err := validateField(fieldName, int(elem.Int()), validateTag, errs); err != nil {
				return err
			}
		default:
			return ErrSysUnsupportedSlice
		}
	}
	return nil
}
