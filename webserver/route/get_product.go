package route

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
)

func ProductGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, _ := template.ParseFiles("template/layout.gohtml", "template/get_product.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"error":          r.URL.Query().Get("error"),
		"success":        r.URL.Query().Get("success"),
		"Product":        D.GetProduct(vars["id"]),
	})
}
