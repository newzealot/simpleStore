package middleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// GetUserInfo is similar to AccessTokenCheck except that it does not block access
func GetUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AccessTokenCheck")
		rt, err := r.Cookie("RefreshToken")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		at, err1 := r.Cookie("AccessToken")
		it, err2 := r.Cookie("IdToken")
		if err1 != nil || err2 != nil {
			// get new Tokens using RefreshToken
			client := &http.Client{}
			req, err := http.NewRequest("POST", os.Getenv("APISERVER")+"/api/v1/refresh", nil)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			req.Header.Add("Authorization", "Bearer "+rt.Value)
			resp, err := client.Do(req)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if resp.StatusCode != http.StatusOK {
				next.ServeHTTP(w, r)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				next.ServeHTTP(w, r)
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
			resp.Body.Close()
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
			at = &c1
			log.Println("Cookies set")
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
		// verify AccessToken
		u, err := VerifyAccessToken("access", at.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		u, err = VerifyAccessToken("id", it.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Add("Authorization", "Bearer "+at.Value)
		// setting UserInfo in header
		r.Header.Add("SimpleStoreUserType", u.Type)
		r.Header.Add("SimpleStoreUserID", u.ID)
		r.Header.Add("SimpleStoreUserEmail", u.Email)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//
//log.Println("GetUserInfo")
//it, err := r.Cookie("IdToken")
//if err != nil {
//next.ServeHTTP(w, r)
//return
//}
//// verify IdToken
//u, err := VerifyAccessToken("id", it.Value)
//if err != nil {
//next.ServeHTTP(w, r)
//return
//}
//// setting UserInfo in header
//r.Header.Add("SimpleStoreUserType", u.Type)
//r.Header.Add("SimpleStoreUserID", u.ID)
//r.Header.Add("SimpleStoreUserEmail", u.Email)
//// Call the next handler, which can be another middleware in the chain, or the final handler.
//next.ServeHTTP(w, r)
