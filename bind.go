package model

import (
	"errors"
	"reflect"
)

const (
	k_MODEL_TAG        = "model"
	k_MODEL_NO_TAG     = "-"
	k_MODEL_CLEAN_DATA = "CleanData"
)

func Bind(source map[string]interface{}, result interface{}) (err error) {
	var objType = reflect.TypeOf(result)
	var objValue = reflect.ValueOf(result)
	var objValueKind = objValue.Kind()

	if objValueKind == reflect.Struct {
		return errors.New("obj is struct")
	}

	if objValue.IsNil() {
		return errors.New("obj is nil")
	}

	for {
		if objValueKind == reflect.Ptr && objValue.IsNil() {
			objValue.Set(reflect.New(objType.Elem()))
		}

		if objValueKind == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			objValueKind = objValue.Kind()
			continue
		}
		break
	}

	var cleanDataValue = objValue.FieldByName(k_MODEL_CLEAN_DATA)
	if cleanDataValue.IsValid() && cleanDataValue.IsNil() {
		cleanDataValue.Set(reflect.MakeMap(cleanDataValue.Type()))
	}
	bindWithMap(objType, objValue, cleanDataValue, source)
	return nil
}

func bindWithMap(objType reflect.Type, objValue, cleanDataValue reflect.Value, source map[string]interface{}) {
	var numField = objType.NumField()
	for i := 0; i < numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		var tag = fieldStruct.Tag.Get(k_MODEL_TAG)

		if tag == "" {
			tag = fieldStruct.Name

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct {
				bindWithMap(fieldValue.Addr().Type().Elem(), fieldValue, cleanDataValue, source)
				continue
			}

		} else if tag == k_MODEL_NO_TAG {
			continue
		}

		var value, exists = source[tag]
		if !exists {
			continue
		}

		fieldValue.Set(reflect.ValueOf(value))

		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
}
