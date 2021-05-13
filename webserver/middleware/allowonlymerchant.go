package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type UserInfo struct {
	Type  string
	ID    string
	Email string
}

func VerifyToken(token_use string, at string) (UserInfo, error) {
	u := UserInfo{}
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
	keyset, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		return u, err
	}
	src := []byte(at)
	iss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
	p, err := jwt.Parse(src,
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true),
		jwt.WithIssuer(iss),
		jwt.WithClaimValue("token_use", token_use),
	)
	if err != nil {
		return u, err
	}
	// no futher processing if access token
	if token_use == "access" {
		return u, nil
	}
	result, _ := p.Get("custom:type")
	u.Type = fmt.Sprint(result)
	result, _ = p.Get("cognito:username")
	u.ID = fmt.Sprint(result)
	result, _ = p.Get("email")
	u.Email = fmt.Sprint(result)
	return u, nil
}

// AllowOnlyMerchant makes sure that protected route contains a valid Access and ID token.
// Customer types are denied access.
func AllowOnlyMerchant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rt, err := r.Cookie("RefreshToken")
		if err != nil {
			log.Printf("RefreshToken error - %s\n", err)
			http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
			return
		}
		at, err1 := r.Cookie("AccessToken")
		it, err2 := r.Cookie("IdToken")
		if err1 != nil || err2 != nil {
			// get new Tokens using RefreshToken
			client := &http.Client{}
			req, err := http.NewRequest("POST", os.Getenv("APISERVER")+"/api/v1/refresh", nil)
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
		u, err := VerifyToken("access", at.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
			return
		}
		u, err = VerifyToken("id", it.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login?status="+"Please login first", http.StatusSeeOther)
			return
		}
		if u.Type == "customer" {
			log.Println("Trying to access admin as customer")
			http.Redirect(w, r, "/?error=Trying to access admin as customer", http.StatusSeeOther)
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
