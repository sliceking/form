package form

import "reflect"

func fields(strct interface{}) []field {
	// using reflect to inspect the structs at runtime
	rv := reflect.ValueOf(strct)
	t := rv.Type()
	var ret []field
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		rvf := rv.Field(i)
		// go to the next field if we cant interact with it
		if !rvf.CanInterface() {
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
