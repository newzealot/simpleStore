package route

import (
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
)

func IndexGET(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/layout.gohtml", "template/index.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Collections":  D.GetMenu(),
		"ProductStore": D.GetCollection("All"),
	})
}
