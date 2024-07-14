package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrorUnknownRule      = errors.New("unknown validation rule")
	ErrorUnknownFieldType = errors.New("unknown field type")
	ErrorIsNotStructType  = errors.New("not a struct type")
	ErrorStrLengthInvalid = errors.New("value length is invalid")

	ErrorMaxValue     = errors.New("value is greater than max")
	ErrorMinValue     = errors.New("value is less than min")
	ErrorIsNotInSlice = errors.New("value is not in slice")
	ErrorIsNotMatch   = errors.New("value is not match regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultString string
	for _, err := range v {
		resultString += fmt.Sprintf("error in field: %v, %v\n", err.Field, err.Err)
	}
	return resultString
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return ErrorIsNotStructType
	}

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		validators := strings.Split(tag, "|")

		for _, validator := range validators {
			parts := strings.SplitN(validator, ":", 2)
			rule := parts[0]
			param := parts[1]

			switch fieldValue.Kind() {
			case reflect.Int:
				err := validateInt(int(fieldValue.Int()), rule, param)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			case reflect.Slice, reflect.Array:
				if fieldValue.Len() > 0 {
					switch fieldValue.Index(0).Kind() {
					case reflect.Int:
						for j := 0; j < fieldValue.Len(); j++ {
							err := validateInt(int(fieldValue.Index(j).Int()), rule, param)
							if err != nil {
								validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
							}
						}
					case reflect.String:
						for j := 0; j < fieldValue.Len(); j++ {
							err := validateString(fieldValue.Index(j).String(), rule, param)
							if err != nil {
								validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
							}
						}
					default:
						return ErrorUnknownFieldType
					}
				}
			case reflect.String:
				err := validateString(fieldValue.String(), rule, param)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			default:
				return ErrorUnknownFieldType
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func errorFunc(value, rule string) error {
	return fmt.Errorf("%s is not a valid value for validator %s", value, rule)
}

func validateInt(value int, rule string, param string) error {
	switch rule {
	case "min":
		min, err := strconv.Atoi(param)
		if err != nil {
			return errorFunc(param, rule)
		}
		if value < min {
			return ErrorMinValue
		}
	case "max":
		max, err := strconv.Atoi(param)
		if err != nil {
			return errorFunc(param, rule)
		}
		if value > max {
			return ErrorMaxValue
		}
	case "in":
		split := strings.Split(param, ",")

		in := make([]int, len(split))
		for i, s := range split {
			v, err := strconv.Atoi(s)
			if err != nil {
				return errorFunc(split[i], rule)
			}
			in[i] = v
		}

		if !slices.Contains(in, value) {
			return ErrorIsNotInSlice
		}
	default:
		return ErrorUnknownRule
	}
	return nil
}

func validateString(value string, rule string, param string) error {
	switch rule {
	case "len":
		length, err := strconv.Atoi(param)
		if err != nil {
			return errorFunc(param, rule)
		}
		if len(value) != length {
			return ErrorStrLengthInvalid
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return errorFunc(param, rule)
		}
		if !re.MatchString(value) {
			return ErrorIsNotMatch
		}
	case "in":
		split := strings.Split(param, ",")
		if !slices.Contains(split, value) {
			return ErrorIsNotInSlice
		}
	default:
		return ErrorUnknownRule
	}

	return nil
}
