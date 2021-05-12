package route

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	. "simpleStore/webserver/data"
	"strings"
)

func DeleteMediaGET(w http.ResponseWriter, r *http.Request) {
	at := r.Header.Get("Authorization")
	at = strings.Replace(at, "Bearer ", "", 1)
	vars := mux.Vars(r)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", os.Getenv("APISERVER")+"/api/v1/media/"+vars["id"]+"/"+vars["filename"], nil)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/edit_product/"+vars["id"]+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	req.Header.Add("Authorization", "Bearer "+at)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/edit_product/"+vars["id"]+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Redirect(w, r, "/admin/edit_product/"+vars["id"]+"?error=Unable to delete media", http.StatusSeeOther)
		return
	}
	resp.Body.Close()
	D.GetData()
	http.Redirect(w, r, "/admin/edit_product/"+vars["id"]+"?success=Media deleted successfully", http.StatusSeeOther)
}
