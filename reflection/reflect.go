package reflection

import "reflect"

func IsLiteralType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.String, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func IsCustomType(t reflect.Type) bool {
	if t.PkgPath() != "" {
		return true
	}

	if k := t.Kind(); k == reflect.Array || k == reflect.Chan || k == reflect.Map ||
		k == reflect.Ptr || k == reflect.Slice {
		return IsCustomType(t.Elem()) || k == reflect.Map && IsCustomType(t.Key())
	} else if k == reflect.Struct {
		for i := t.NumField() - 1; i >= 0; i-- {
			if IsCustomType(t.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

func HasUnexportedField(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if IsUnexportedField(field) {
			return true
		}
		if field.Type.Kind() == reflect.Struct && HasUnexportedField(field.Type) {
			return true
		}
	}
	return false
}

func IsUnexportedField(field reflect.StructField) bool {
	if len(field.PkgPath) > 0 {
		return true
	}
	return false
}
