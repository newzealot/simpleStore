package route

import (
	"fmt"
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
	. "simpleStore/webserver/middleware"
)

func AdminGET(w http.ResponseWriter, r *http.Request) {
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/admin.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"User":      u,
		"AdminList": D.GetAdmin(fmt.Sprint(u.ID)),
		"error":     r.URL.Query().Get("error"),
		"success":   r.URL.Query().Get("success"),
	})
}
