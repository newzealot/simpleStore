package route

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
)

func CollectionGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "All" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/collection.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Collections":    D.GetMenu(),
		"ProductStore":   D.GetCollection(vars["id"]),
		"CollectionID":   vars["id"],
	})
}
