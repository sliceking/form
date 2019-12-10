package form

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := map[string]struct {
		arg  reflect.StructField
		want map[string]string
	}{
		"empty tag": {
			arg:  reflect.StructField{},
			want: nil,
		},
		"label tag": {
			arg: reflect.StructField{
				Tag: `form:"label=Full Name"`,
			},
			want: map[string]string{
				"label": "Full Name",
			},
		},
		"multiple tags": {
			arg: reflect.StructField{
				Tag: `form:"label=Full Name;name=full_name"`,
			},
			want: map[string]string{
				"label": "Full Name",
				"name":  "full_name",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseTags(tc.arg)
			if len(got) != len(tc.want) {
				t.Errorf("parseTags() len = %d, want %d", len(got), len(tc.want))
			}
			for k, v := range tc.want {
				gv, ok := got[k]
				if !ok {
					t.Errorf("parseTags() missing key %q", k)
					continue
				}
				if gv != v {
					t.Errorf("parseTags()[%q] = %q; want %q", k, gv, v)
				}
				delete(got, k)
			}
			for gk, gv := range got {
				t.Errorf("parseTags() extra key %q, value = %q", gk, gv)
			}
		})
	}
}

func TestParseTags_invalidTag(t *testing.T) {
	tests := []struct {
		arg reflect.StructField
	}{
		{reflect.StructField{Tag: `form:"invalid-value"`}},
	}
	for _, tc := range tests {
		t.Run(string(tc.arg.Tag), func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("parseTags() did not panic")
				}
			}()
			parseTags(tc.arg)
		})
	}
}

func TestFields(t *testing.T) {
	var nilStructPtr *struct {
		Name string
		Age  int
	}

	tests := map[string]struct {
		strct interface{}
		want  []field
	}{
		"Simplest use case": {
			strct: struct {
				Name string
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
			},
		},
		"Field names should be determined from the struct": {
			strct: struct {
				FullName string
			}{},
			want: []field{
				{
					Label:       "FullName",
					Name:        "FullName",
					Type:        "text",
					Placeholder: "FullName",
					Value:       "",
				},
			},
		},
		"Multiple fields should be supported": {
			strct: struct {
				Name  string
				Age   int
				Email string
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       0,
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "",
				},
			},
		},
		"Values should be parsed": {
			strct: struct {
				Name  string
				Age   int
				Email string
			}{
				Name:  "stanny",
				Age:   34,
				Email: "bob@lawblog.com",
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "stanny",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       34,
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "bob@lawblog.com",
				},
			},
		},
		"Unexported fields should be skipped": {
			strct: struct {
				Name  string
				Age   int
				email string
			}{
				Name:  "stanny",
				Age:   34,
				email: "bob@lawblog.com",
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "stanny",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       34,
				},
			},
		},
		"Pointers to structs should be supported": {
			strct: &struct {
				Name string
				Age  int
			}{
				Name: "stanny",
				Age:  34,
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "stanny",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       34,
				},
			},
		},
		"Nil pointers with a struct type should be supported": {
			strct: nilStructPtr,
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       0,
				},
			},
		},
		"Nested structs should be supported": {
			strct: struct {
				Name    string
				Address struct {
					Street string
					Zip    int
				}
			}{
				Name: "Stanley",
				Address: struct {
					Street string
					Zip    int
				}{
					Street: "235 Arcadia rd",
					Zip:    43253,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "Stanley",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "235 Arcadia rd",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
					Value:       43253,
				},
			},
		},
		"Doubly nested structs should be supported": {
			strct: struct {
				A struct {
					B struct {
						C1 string
						C2 int
					}
				}
			}{
				A: struct {
					B struct {
						C1 string
						C2 int
					}
				}{
					B: struct {
						C1 string
						C2 int
					}{
						C1: "C1-value",
						C2: 123,
					},
				},
			},
			want: []field{
				{
					Label:       "C1",
					Name:        "A.B.C1",
					Type:        "text",
					Placeholder: "C1",
					Value:       "C1-value",
				},
				{
					Label:       "C2",
					Name:        "A.B.C2",
					Type:        "text",
					Placeholder: "C2",
					Value:       123,
				},
			},
		},
		"Nested pointer structs should be supported": {
			strct: struct {
				Name    string
				Address *struct {
					Street string
					Zip    int
				}
				ContactCard *struct {
					Phone string
				}
			}{
				Name: "Jon Calhoun",
				Address: &struct {
					Street string
					Zip    int
				}{
					Street: "123 Fake St",
					Zip:    90210,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "Jon Calhoun",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "123 Fake St",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
					Value:       90210,
				},
				{
					Label:       "Phone",
					Name:        "ContactCard.Phone",
					Type:        "text",
					Placeholder: "Phone",
					Value:       "",
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := fields(tc.strct)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("fields () = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestFields_invalidValues(t *testing.T) {
	tests := []struct {
		notAStruct interface{}
	}{
		{"this is a string"},
		{123},
		{nil},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.notAStruct), func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("fields() did not panic: %v", tc.notAStruct)
				}
			}()
			fields(tc.notAStruct)
		})
	}
}
