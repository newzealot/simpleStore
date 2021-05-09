package route

import (
	"fmt"
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
)

func AdminGET(w http.ResponseWriter, r *http.Request) {
	merchantID := r.Header.Get("merchantID")
	t, _ := template.ParseFiles("template/layout.gohtml", "template/admin.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"AdminList": D.GetAdmin(fmt.Sprint(merchantID)),
		"error":     r.URL.Query().Get("error"),
		"success":   r.URL.Query().Get("success"),
	})
}
