package route

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

func CartGET(w http.ResponseWriter, r *http.Request) {
	orderList := []Order{}
	cart, err := r.Cookie("Cart")
	if err != nil {
		log.Println("No cookie")
	} else {
		j, err := base64.StdEncoding.DecodeString(cart.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, r.URL.Path+"?error=Cannot decode base64 string", http.StatusSeeOther)
			return
		}
		if err := json.Unmarshal(j, &orderList); err != nil {
			log.Println(err)
			http.Redirect(w, r, r.URL.Path+"?error=Cannot unmarshall json", http.StatusSeeOther)
			return
		}
	}
	total := 0.0
	for _, v := range orderList {
		total += v.SubTotal
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/cart.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Cart":      orderList,
		"error":     r.URL.Query().Get("error"),
		"success":   r.URL.Query().Get("success"),
		"Total":     total,
		"StripeKey": os.Getenv("STRIPE_PUBLIC_KEY"),
	})
}
