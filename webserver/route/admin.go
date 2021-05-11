package route

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	. "simpleStore/webserver/data"
)

func AdminGET(w http.ResponseWriter, r *http.Request) {
	userType := r.Header.Get("SimpleStoreUserType")
	if userType == "customer" {
		log.Println("Trying to access admin as customer")
		http.Redirect(w, r, "/admin?error=Trying to access admin as customer", http.StatusSeeOther)
	}
	merchantID := r.Header.Get("SimpleStoreUserID")
	t, _ := template.ParseFiles("template/layout.gohtml", "template/admin.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"AdminList": D.GetAdmin(fmt.Sprint(merchantID)),
		"error":     r.URL.Query().Get("error"),
		"success":   r.URL.Query().Get("success"),
	})
}
