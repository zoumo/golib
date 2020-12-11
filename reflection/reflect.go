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

// IsCustomType returns true if the type is not predeclared
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

// HasUnexportedField returns true if the struct has unexported fields.
// It panics if the type's Kind is not Struct.
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

// IsUnexportedField returns true if the structField is not exported.
func IsUnexportedField(field reflect.StructField) bool {
	return len(field.PkgPath) > 0
}

// Hashable returns true if the type can used as map key
// see https://blog.golang.org/maps
func Hashable(in reflect.Type) bool {
	switch in.Kind() {
	case reflect.Invalid, reflect.Map, reflect.Func, reflect.Slice:
		return false
	case reflect.Struct:
		for i := 0; i < in.NumField(); i++ {
			if !Hashable(in.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Array:
		return Hashable(in.Elem())
	}
	return true
}

// https://stackoverflow.com/questions/36310538/identify-non-builtin-types-using-reflect?answertab=votes#tab-top
func IsAnonymousStruct(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	if len(t.PkgPath()) == 0 && len(t.Name()) == 0 {
		// not custom and unnamed
		return true
	}
	return false
}
