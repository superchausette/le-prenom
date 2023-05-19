package leprenom

import "html/template"

func mod(value, modulo int) int {
	return value % modulo
}

func seq(count uint) []uint {
	ret := make([]uint, count)
	for i := range ret {
		ret[i] = uint(i)
	}
	return ret
}

func percent(value, total uint) float64 {
	return float64(value) * 100. / float64(total)
}

func NewTemplates() *template.Template {
	funcMap := template.FuncMap{
		"mod":              mod,
		"percent":          percent,
		"seq":              seq,
		"sessionTypeToStr": SessionTypeToString}

	files := []string{"template/index.html",
		"template/404.html",
		"template/list.html",
		"template/partial/firstname_list.html",
		"template/partial/footer.html",
		"template/partial/header.html",
		"template/partial/session_first_name_kept.html",
		"template/partial/session_first_name_table_entry.html",
		"template/partial/session_list.html",
		"template/partial/stats.html",
		"template/session.html",
	}
	return template.Must(template.New("templates").
		Funcs(funcMap).
		ParseFiles(files...))
}
