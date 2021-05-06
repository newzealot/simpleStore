package route

import (
	"github.com/gorilla/csrf"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func RegisterGET(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/layout.gohtml", "template/register.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"error":          r.URL.Query().Get("error"),
	})
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("APISERVER")+"/api/v1/register", strings.NewReader(r.PostForm.Encode()))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Unable to connect to other server", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			http.Redirect(w, r, r.URL.Path+"?error=Invalid password. Check if both passwords are the same.", http.StatusSeeOther)
			return
		case http.StatusConflict:
			http.Redirect(w, r, r.URL.Path+"?error=Email already exists", http.StatusSeeOther)
			return
		default:
			http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
			return
		}
	}
	log.Println("Registration successful")
	http.Redirect(w, r, "/login?info=Please check your email, click the confirmation link, then login here", http.StatusSeeOther)
}
