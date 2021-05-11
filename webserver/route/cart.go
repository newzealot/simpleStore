package route

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/csrf"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func deleteItemFromCart(toDelete string, orderList []Order) (http.Cookie, error) {
	deleteIndex := 0
	for i, v := range orderList {
		if v.ProductID == toDelete {
			deleteIndex = i
		}
	}
	copy(orderList[deleteIndex:], orderList[deleteIndex+1:])
	orderList = orderList[:len(orderList)-1]
	j, err := json.Marshal(orderList)
	jbase := base64.StdEncoding.EncodeToString(j)
	if err != nil {
		return http.Cookie{}, err
	}
	c := http.Cookie{
		Name:     "Cart",
		Value:    jbase,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
	}
	return c, nil
}

func CartGET(w http.ResponseWriter, r *http.Request) {
	toDelete := r.URL.Query().Get("delete")
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
	if toDelete != "" {
		c, err := deleteItemFromCart(toDelete, orderList)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/cart?error=Cannot delete item", http.StatusSeeOther)
			return
		}
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/cart?success=Item deleted", http.StatusSeeOther)
		return
	}
	total := 0.0
	for _, v := range orderList {
		total += v.SubTotal
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/cart.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Cart":           orderList,
		"error":          r.URL.Query().Get("error"),
		"success":        r.URL.Query().Get("success"),
		"Total":          total,
		"StripeKey":      os.Getenv("STRIPE_PUBLIC_KEY"),
	})
}
