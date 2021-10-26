package schema

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrUnsupportedDataType = errors.New("unsupported data type")

func Parse(dest interface{}) error {
	if dest == nil {
		return fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
	}

	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(dest)).Elem().Type()
	}

	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			return fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
		}
		return fmt.Errorf("%w: %s.%s", ErrUnsupportedDataType, modelType.PkgPath(), modelType.Name())
	}

	return nil
}
