package route

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	. "simpleStore/webserver/data"
)

func AdminGET(w http.ResponseWriter, r *http.Request) {
	merchantID := r.Header.Get("merchantID")
	t, _ := template.ParseFiles("template/layout.gohtml", "template/admin.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"AdminList":      D.GetAdmin(fmt.Sprint(merchantID)),
	})
}
