package route

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	. "simpleStore/webserver/data"
	"strings"
)

func DeleteProductGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, _ := template.ParseFiles("template/layout.gohtml", "template/delete_product.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Product":        D.GetProduct(vars["id"]),
		"error":          r.URL.Query().Get("error"),
		"info":           r.URL.Query().Get("info"),
	})
}

func DeleteProductPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.PostFormValue("yesno"))
	fmt.Sprintln()
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
	log.Println("Product deleted successfully")
	D.GetData()
	http.Redirect(w, r, "/admin?success=Product deleted successfully", http.StatusSeeOther)
}
