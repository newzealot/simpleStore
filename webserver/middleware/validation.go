package middleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func Validation(next http.Handler) http.Handler {
	log.Println("ValidationMiddleWare")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rt, err := r.Cookie("RefreshToken")
		if err != nil {
			log.Printf("RefreshToken error - %s\n", err)
			http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
			return
		}
		_, err = r.Cookie("AccessToken")
		if err != nil {
			client := &http.Client{}
			req, err := http.NewRequest("POST", os.Getenv("SERVER")+"/api/v1/refresh", nil)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
				return
			}
			req.Header.Add("Authorization", "Bearer "+rt.Value)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println(err)
				http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
				return
			}
			var j struct {
				AccessToken string
				ExpiresIn   int64
				IdToken     string
			}
			if err = json.Unmarshal(body, &j); err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
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
			http.SetCookie(w, &c1)
			http.SetCookie(w, &c2)
			log.Println("Cookies set")
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
