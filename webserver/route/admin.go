package route

import (
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

func AdminGET(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/layout.gohtml", "template/admin.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
