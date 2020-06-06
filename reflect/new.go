package reflectUtils

import (
	"fmt"
	"reflect"
	"unsafe"
)

func SetStructPtrUnExportedStrField(source interface{}, fieldName string, fieldVal interface{}) (err error) {
	v := GetStructPtrUnExportedField(source, fieldName)
	rv := reflect.ValueOf(fieldVal)
	if v.Kind() != rv.Kind() {
		return fmt.Errorf("invalid kind: expected kind %v, got kind: %v", v.Kind(), rv.Kind())
	}

	v.Set(rv)
	return nil
}

func SetStructUnExportedStrField(source interface{}, fieldName string, fieldVal interface{}) (addressableSourceCopy reflect.Value, err error) {
	var accessableField reflect.Value
	accessableField, addressableSourceCopy = GetStructUnExportedField(source, fieldName)
	rv := reflect.ValueOf(fieldVal)
	if accessableField.Kind() != rv.Kind() {
		return addressableSourceCopy, fmt.Errorf("invalid kind: expected kind %v, got kind: %v", addressableSourceCopy.Kind(), rv.Kind())
	}
	accessableField.Set(rv)
	return
}

func GetStructPtrUnExportedField(source interface{}, fieldName string) reflect.Value {
	v := reflect.ValueOf(source).Elem().FieldByName(fieldName)
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

func GetStructUnExportedField(source interface{}, fieldName string) (accessableField, addressableSourceCopy reflect.Value) {
	v := reflect.ValueOf(source)
	// since source is not a ptr, get an addressable copy of source to modify it later
	addressableSourceCopy = reflect.New(v.Type()).Elem()
	addressableSourceCopy.Set(v)
	accessableField = addressableSourceCopy.FieldByName(fieldName)
	accessableField = reflect.NewAt(accessableField.Type(), unsafe.Pointer(accessableField.UnsafeAddr())).Elem()
	return
}
