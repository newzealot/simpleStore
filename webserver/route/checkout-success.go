package route

import (
	"github.com/gorilla/csrf"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"html/template"
	"net/http"
	"os"
	. "simpleStore/webserver/middleware"
	"time"
)

func CheckoutSuccessGET(w http.ResponseWriter, r *http.Request) {
	// clear cart
	c := http.Cookie{
		Name:     "Cart",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, &c)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	s, _ := session.Get(r.URL.Query()["session_id"][0], nil)
	orderID := s.Metadata["SimpleStoreOrderID"]
	t, _ := template.ParseFiles("template/layout.gohtml", "template/checkout-success.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag:     csrf.TemplateField(r),
		"User":               u,
		"SimpleStoreOrderID": orderID,
		"error":              r.URL.Query().Get("error"),
		"success":            r.URL.Query().Get("success"),
	})
}
