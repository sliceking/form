package form

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
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
