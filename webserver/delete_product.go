package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func DeleteProductGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/delete_product.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           u,
		"Product":        D.GetProduct(vars["id"]),
		"error":          r.URL.Query().Get("error"),
		"info":           r.URL.Query().Get("info"),
	})
}

func DeleteProductPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostFormValue("yesno") == "no" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	at := r.Header.Get("Authorization")
	at = strings.Replace(at, "Bearer ", "", 1)
	vars := mux.Vars(r)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", os.Getenv("APISERVER")+"/api/v1/product/"+vars["id"], nil)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin?error=Something went wrong", http.StatusSeeOther)
		return
	}
	req.Header.Add("Authorization", "Bearer "+at)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin?error=Something went wrong", http.StatusSeeOther)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Redirect(w, r, "/admin?error=Something went wrong", http.StatusSeeOther)
		return
	}
	resp.Body.Close()
	D.GetData()
	http.Redirect(w, r, "/admin?success=Product deleted successfully", http.StatusSeeOther)
}
