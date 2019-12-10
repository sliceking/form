package form

import (
	"reflect"
	"strings"
)

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

func parseTags(sf reflect.StructField) map[string]string {
	// label=Full Name;name=full_name
	rawTag := sf.Tag.Get("form")
	if len(rawTag) == 0 {
		return nil
	}
	ret := make(map[string]string)
	tags := strings.Split(rawTag, ";")
	for _, tag := range tags {
		kv := strings.Split(tag, "=")
		if len(kv) != 2 {
			panic("form: invalid struct tag")
		}
		k, v := kv[0], kv[1]
		ret[k] = v
	}
	return ret
}
