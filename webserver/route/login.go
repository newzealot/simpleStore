package route

import (
	"encoding/json"
	"github.com/gorilla/csrf"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	. "simpleStore/webserver/middleware"
	"strings"
	"time"
)

func LoginGET(w http.ResponseWriter, r *http.Request) {
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/login.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           u,
		"error":          r.URL.Query().Get("error"),
		"info":           r.URL.Query().Get("info"),
	})
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("APISERVER")+"/api/v1/login", strings.NewReader(r.PostForm.Encode()))
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
		case http.StatusUnauthorized:
			http.Redirect(w, r, r.URL.Path+"?error=Incorrect username or password", http.StatusSeeOther)
			return
		case http.StatusForbidden:
			http.Redirect(w, r, r.URL.Path+"?error=Email not confirmed", http.StatusSeeOther)
			return
		default:
			http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
			return
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	var j struct {
		AccessToken  string
		ExpiresIn    int64
		IdToken      string
		RefreshToken string
	}
	if err = json.Unmarshal(body, &j); err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	c1 := http.Cookie{
		Name:     "AccessToken",
		Value:    j.AccessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Second * time.Duration(j.ExpiresIn)),
	}
	c2 := http.Cookie{
		Name:     "IdToken",
		Value:    j.IdToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Second * time.Duration(j.ExpiresIn)),
	}
	c3 := http.Cookie{
		Name:     "RefreshToken",
		Value:    j.RefreshToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
	}
	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
	http.SetCookie(w, &c3)
	log.Println("Cookies set")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
