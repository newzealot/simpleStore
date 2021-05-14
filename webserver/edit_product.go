package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func EditProductGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := UserInfo{
		Type:  r.Header.Get("SimpleStoreUserType"),
		ID:    r.Header.Get("SimpleStoreUserID"),
		Email: r.Header.Get("SimpleStoreUserEmail"),
	}
	t, _ := template.ParseFiles("template/layout.gohtml", "template/edit_product.gohtml")
	t.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"User":           u,
		"error":          r.URL.Query().Get("error"),
		"success":        r.URL.Query().Get("success"),
		"Product":        D.GetProduct(vars["id"]),
	})
}

func EditProductPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	at := r.Header.Get("Authorization")
	at = strings.Replace(at, "Bearer ", "", 1)
	// Continue to process form
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	formdata := r.MultipartForm
	// Copying from body to new writer
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, v := range formdata.Value {
		br, _ := bodyWriter.CreateFormField(k)
		j, _ := json.Marshal(v[0])
		io.Copy(br, bytes.NewBuffer(j))
	}
	for _, v := range formdata.File["media[]"] {
		fileWriter, _ := bodyWriter.CreateFormFile("media", v.Filename)
		fh, _ := v.Open()
		io.Copy(fileWriter, fh)
		fh.Close()
	}
	if err := bodyWriter.Close(); err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", os.Getenv("APISERVER")+"/api/v1/product/"+vars["id"], bytes.NewBuffer(bodyBuf.Bytes()))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+at)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusInsufficientStorage:
			http.Redirect(w, r, r.URL.Path+"?error=Issue with uploading media files", http.StatusSeeOther)
			return
		default:
			http.Redirect(w, r, r.URL.Path+"?error=Something went wrong", http.StatusSeeOther)
			return
		}
	}
	D.GetData()
	http.Redirect(w, r, r.URL.Path+"?success=Successfully edited product", http.StatusSeeOther)
	return
}
