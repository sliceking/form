package form

func HTML(strct interface{}, tpl *template.Template) template.HTML{
	return template.HTML
}

func fields(strct interface{}) field {
	return field{}
}

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
}
