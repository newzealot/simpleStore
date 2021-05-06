package route

import (
	"net/http"
	"time"
)

func LogoutGET(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "AccessToken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	c2 := http.Cookie{
		Name:     "IdToken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	c3 := http.Cookie{
		Name:     "RefreshToken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
	http.SetCookie(w, &c3)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
