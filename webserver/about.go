package main

import (
	"html/template"
	"net/http"
)

func AboutGET(w http.ResponseWriter, r *http.Request) {
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/about.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"User":    u,
		"error":   r.URL.Query().Get("error"),
		"success": r.URL.Query().Get("success"),
	})
}
