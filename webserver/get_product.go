package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func ProductGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/get_product.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           u,
		"error":          r.URL.Query().Get("error"),
		"success":        r.URL.Query().Get("success"),
		"Product":        D.GetProduct(vars["id"]),
	})
}
