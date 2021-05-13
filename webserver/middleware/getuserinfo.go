package middleware

import (
	"net/http"
)

// GetUserInfo is similar to AccessTokenCheck except that it does not block access
func GetUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rt, err := r.Cookie("RefreshToken")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		at, err1 := r.Cookie("AccessToken")
		it, err2 := r.Cookie("IdToken")
		if err1 != nil || err2 != nil {
			// get new Tokens using RefreshToken
			GetNewTokens(w, r, rt.Value)
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
		// verify AccessToken
		u, err := VerifyToken("access", at.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		u, err = VerifyToken("id", it.Value)
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
