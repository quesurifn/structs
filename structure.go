// Package structure contains various utilities functions to work with structs.
package structure

import "reflect"

// Map converts the given s struct to a map[string]interface{}, where the keys
// of the map are the field names and the values of the map the associated
// values of the fields. The default key string is the struct field name but
// can be changed in the struct field's tag value. The "structure" key in the
// struct's field tag value is the key name. Example:
//
//   // Field appears in map as key "myName".
//   Name string `structure:"myName"`
//
// A value with the content of "-" ignores that particular field. Example:
//
//   // Field is ignored by this package.
//   Field bool `structure:"-"`
//
// Note that only exported fields of a struct can be accessed, non exported
// fields will be neglected. It panics if s's kind is not struct.
func Map(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v, fields := strctInfo(s)

	for i, field := range fields {
		name := field.Name

		// override if the user passed a structure tag value
		// ignore if the user passed the "-" value
		if tag := field.Tag.Get("structure"); tag != "" {
			name = tag
		}

		out[name] = v.Field(i).Interface()
	}

	return out
}

// Values converts the given s struct's field values to a []interface{}.  A
// struct tag with the content of "-" ignores the that particular field.
// Example:
//
//   // Field is ignored by this package.
//   Field int `structure:"-"`
//
// Note that only exported fields of a struct can be accessed, non exported
// fields  will be neglected.  It panics if s's kind is not struct.
func Values(s interface{}) []interface{} {
	v, fields := strctInfo(s)

	t := make([]interface{}, len(fields))
	for i := range fields {
		t[i] = v.Field(i).Interface()
	}

	return t

}

// IsValid returns true if all fields in a struct are initialized (non zero
// value). A struct tag with the content of "-" ignores the checking of that
// particular field. Example:
//
//   // Field is ignored by this package.
//   Field bool `structure:"-"`
//
// Note that only exported fields of a struct can be accessed, non exported
// fields  will be neglected. It panics if s's kind is not struct.
func IsValid(s interface{}) bool {
	v, fields := strctInfo(s)

	for i := range fields {
		// zero value of the given field, such as "" for string, 0 for int
		zero := reflect.Zero(v.Field(i).Type()).Interface()

		//  current value of the given field
		current := v.Field(i).Interface()

		if reflect.DeepEqual(current, zero) {
			return false
		}
	}

	return true
}

// Fields returns a slice of field names. A struct tag with the content of "-"
// ignores the checking of that particular field. Example:
//
//   // Field is ignored by this package.
//   Field bool `structure:"-"`
//
// Note that only exported fields of a struct can be accessed, non exported
// fields  will be neglected. It panics if s's kind is not struct.
func Fields(s interface{}) []string {
	_, fields := strctInfo(s)

	keys := make([]string, len(fields))
	for i, field := range fields {
		keys[i] = field.Name
	}

	return keys
}

// IsStruct returns true if the given variable is a struct or a pointer to
// struct.
func IsStruct(s interface{}) bool {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}

//  Name returns the structs's type name within its package. It returns an
//  empty string for unnamed types. It panics if s's kind is not struct.
func Name(s interface{}) string {
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("not struct")
	}

	return t.Name()
}

// strctInfo returns the struct value and the exported struct fields for a
// given s struct. This is a convenient helper method to avoid duplicate code
// in some of the functions.
func strctInfo(s interface{}) (reflect.Value, []reflect.StructField) {
	v := strctVal(s)
	t := v.Type()

	f := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get("structure"); tag == "-" {
			continue
		}

		f = append(f, field)
	}

	return v, f
}

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v
}
