package templates

type menuItem struct {
	label string
	href  string
}

var menuItems = []menuItem{
	{"Home", "/"},
	{"Fruit", "/fruit"},
	{"Vegetables", "/vegetables"},
	{"Books", "/books"},
	{"Firearms", "/firearms"},
	{"Cosmetics", "/cosmetics"},
}
