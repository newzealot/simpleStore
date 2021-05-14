package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func CollectionGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "All" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/collection.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"User":         u,
		"Collections":  D.GetMenu(),
		"ProductStore": D.GetCollection(vars["id"]),
		"CollectionID": vars["id"],
	})
}
