package leprenom

import "html/template"

func mod(value, modulo int) int {
	return value % modulo
}

func NewTemplates() *template.Template {
	funcMap := template.FuncMap{
		"mod":              mod,
		"sessionTypeToStr": SessionTypeToString}

	files := []string{"template/index.html",
		"template/list.html",
		"template/404.html",
		"template/partial/firstname_list.html",
		"template/partial/footer.html",
		"template/partial/header.html",
		"template/partial/session_list.html",
		"template/partial/stats.html",
	}
	return template.Must(template.New("templates").
		Funcs(funcMap).
		ParseFiles(files...))
}
