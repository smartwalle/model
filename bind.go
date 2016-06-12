package model

import (
	"errors"
	"reflect"
)

const (
	k_MODEL_TAG         = "model"
	k_MODEL_CONSTRUCTOR = "Constructor"
	k_MODEL_NO_TAG      = "-"
	k_MODEL_CLEAN_DATA  = "CleanData"
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
	return bindWithMap(objType, objValue, cleanDataValue, source)
}

func bindWithMap(objType reflect.Type, objValue, cleanDataValue reflect.Value, source map[string]interface{}) (error) {
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
				if err := bindWithMap(fieldValue.Addr().Type().Elem(), fieldValue, cleanDataValue, source); err != nil {
					return err
				}
				continue
			}

		} else if tag == k_MODEL_NO_TAG {
			continue
		}

		var value, exists = source[tag]
		if !exists {
			continue
		}

		//fieldValue.Set(reflect.ValueOf(value))
		if err := setValue(objValue, fieldValue, fieldStruct, value); err != nil {
			return err
		}

		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
	return nil
}

func setValue(objValue, fieldValue reflect.Value, fieldStruct reflect.StructField, value interface{}) (error) {
	var mName = fieldStruct.Name + k_MODEL_CONSTRUCTOR
	var mValue = objValue.MethodByName(mName)
	if mValue.IsValid() == false {
		if objValue.CanAddr() {
			mValue = objValue.Addr().MethodByName(mName)
		}
	}

	if mValue.IsValid() {
		var rList = mValue.Call([]reflect.Value{reflect.ValueOf(value)})
		if len(rList) > 1 {
			var rValue1 = rList[1]
			if rValue1.IsNil() == false {
				return rValue1.Interface().(error)
			}
		}
		fieldValue.Set(rList[0])
	} else {
		fieldValue.Set(reflect.ValueOf(value))
	}
	return nil
}