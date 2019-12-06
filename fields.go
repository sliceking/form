package form

import "reflect"

func valueOf(v interface{}) reflect.Value {
	var rv reflect.Value
	switch value := v.(type) {
	case reflect.Value:
		rv = value
	default:
		rv = reflect.ValueOf(v)
	}

	// We really just want to work with the underlying type if we get a pointer
	if rv.Kind() == reflect.Ptr {
		// if its nil its kind of useless so we create a new copy
		if rv.IsNil() {
			rv = reflect.New(rv.Type().Elem())
		}
		rv = rv.Elem()
	}

	return rv
}

func fields(strct interface{}) []field {
	// using reflect to inspect the structs at runtime
	rv := valueOf(strct)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		panic("form: invalid value; only structs are supported")
	}
	t := rv.Type()
	var ret []field
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		rvf := valueOf(rv.Field(i))
		// go to the next field if we cant interact with it
		if !rvf.CanInterface() {
			continue
		}
		if rvf.Kind() == reflect.Struct {
			nestedFields := fields(rvf.Interface())
			for i, nf := range nestedFields {
				nestedFields[i].Name = tf.Name + "." + nf.Name
			}
			ret = append(ret, nestedFields...)
			continue
		}

		f := field{
			Label:       tf.Name,
			Name:        tf.Name,
			Type:        "text",
			Placeholder: tf.Name,
			Value:       rvf.Interface(),
		}
		ret = append(ret, f)
	}
	return ret
}

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
}
