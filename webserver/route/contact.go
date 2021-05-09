package route

import (
	"html/template"
	"net/http"
)

func ContactGET(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/layout.gohtml", "template/contact.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"error":   r.URL.Query().Get("error"),
		"success": r.URL.Query().Get("success"),
	})
}
