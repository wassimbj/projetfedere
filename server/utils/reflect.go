package utils

//! reflect package utils

import (
	"reflect"
)

type Rflkt struct {
	Elem reflect.Value
	Data interface{}
}

func Reflekt(data interface{}) *Rflkt {
	return &Rflkt{
		Data: data,
		Elem: reflect.ValueOf(data).Elem(),
	}
}

type FieldDetails struct {
	Name  string
	Type  reflect.Kind
	Value reflect.Value
	Tags  map[string]string // [tag => value, ...]
}

func (r *Rflkt) GetFieldDetails(i int, tags ...string) FieldDetails {

	fieldName := r.Elem.Type().Field(i).Name
	fieldType := r.Elem.Type().Field(i).Type.Kind()
	fieldValue := r.Elem.Field(i)

	fieldDetails := FieldDetails{
		Name:  fieldName,
		Type:  fieldType,
		Value: fieldValue,
	}

	if len(tags) > 0 {
		fieldDetails.Tags = make(map[string]string)
		for _, tag := range tags {
			tagVal := r.GetFieldTag(i, tag)
			fieldDetails.Tags[tag] = tagVal
		}
	}

	return fieldDetails
}

func (r *Rflkt) IsFieldEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (r *Rflkt) GetFieldTag(i int, tag string) string {
	fieldTag := r.Elem.Type().Field(i).Tag.Get(tag)
	return fieldTag
}
