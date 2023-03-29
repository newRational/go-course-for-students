package homework

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var strs []string
	for _, e := range v {
		strs = append(strs, e.Err.Error())
	}
	return strings.Join(strs, ", ")
}

func Validate(val any) error {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)
	var err ValidationErrors

	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)

		if !tf.IsExported() && tf.Tag != "" {
			return append(err, ValidationError{Err: ErrValidateForUnexportedFields})
		}

		if tf.IsExported() && tf.Tag != "" {
			if e := validate(tf, v.Field(i)); e != nil {
				err = append(err, ValidationError{Err: e})
			}
		}
	}

	if len(err) == 0 {
		return nil
	}

	return err
}

func validate(f reflect.StructField, v reflect.Value) error {
	switch v.Interface().(type) {
	case int:
		return validateInt(f, int(v.Int()))
	case string:
		return validateStr(f, v.String())
	case []int:
		for _, n := range v.Interface().([]int) {
			if e := validateInt(f, n); e != nil {
				return e
			}
		}
	case []string:
		for _, s := range v.Interface().([]string) {
			if e := validateStr(f, s); e != nil {
				return e
			}
		}
	}
	return nil
}

func validateInt(f reflect.StructField, num int) error {
	tk := strings.Split(f.Tag.Get("validate"), ":")[0]
	tv := strings.Split(f.Tag.Get("validate"), ":")[1]

	switch tk {
	case "min":
		min, err := strconv.Atoi(tv)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if num < min {
			return errors.Errorf("%v: value is less than %v", f.Name, min)
		}
	case "max":
		max, err := strconv.Atoi(tv)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if num > max {
			return errors.Errorf("%v: value is greater than %v", f.Name, max)
		}
	case "in":
		var in bool
		nums := strings.Split(tv, ",")

		for _, n := range nums {
			i, err := strconv.Atoi(n)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			if i == num {
				in = true
			}
		}

		if !in {
			return errors.Errorf("%v: value is not in %v", f.Name, nums)
		}
	}
	return nil
}

func validateStr(f reflect.StructField, str string) error {
	tk := strings.Split(f.Tag.Get("validate"), ":")[0]
	tv := strings.Split(f.Tag.Get("validate"), ":")[1]

	switch tk {
	case "len":
		l, err := strconv.Atoi(tv)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if l < 0 || len(str) != l {
			return errors.Errorf("%v: len is not equal to %v", f.Name, l)
		}
	case "min":
		min, err := strconv.Atoi(tv)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if len(str) < min {
			return errors.Errorf("%v: len is less than %v", f.Name, min)
		}
	case "max":
		max, err := strconv.Atoi(tv)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if max < 0 || len(str) > max {
			return errors.Errorf("%v: len is greater than %v", f.Name, max)
		}
	case "in":
		var in bool
		strs := strings.Split(tv, ",")

		if len(tv) == 0 {
			return errors.Errorf("%v: value is not in empty %v", f.Name, strs)
		}

		for _, s := range strs {
			if s == str {
				in = true
			}
		}

		if !in {
			return errors.Errorf("%v: value is not in %v", f.Name, strs)
		}
	}
	return nil
}
