package main

import (
	"html/template"
	"net/http"
)

func IndexGET(w http.ResponseWriter, r *http.Request) {
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/index.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"User":         u,
		"error":        r.URL.Query().Get("error"),
		"Collections":  D.GetMenu(),
		"ProductStore": D.GetCollection("All"),
	})
}
