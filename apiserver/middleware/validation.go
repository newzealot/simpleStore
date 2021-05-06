package middleware

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Validation(next http.Handler) http.Handler {
	log.Println("ValidationMiddleWare")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header.Get("Authorization")
		at = strings.Replace(at, "Bearer ", "", 1)
		// verify current AccessToken
		url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
		keyset, err := jwk.Fetch(context.Background(), url)
		src := []byte(at)
		iss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
		p, err := jwt.Parse(src,
			jwt.WithKeySet(keyset),
			jwt.WithValidate(true),
			jwt.WithClaimValue("client_id", os.Getenv("AWS_COGNITO_APP_CLIENT_ID")), // replacing aud
			jwt.WithIssuer(iss),
			jwt.WithClaimValue("token_use", "access"),
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, ok := p.Get("username")
		if ok == false || username == "" {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//at, err := r.Cookie("AccessToken")
//if err != nil {
//// get new AccessToken using RefreshToken
//sess, err := session.NewSession()
//if err != nil {
//log.Println(err)
//}
//authTry := &cognito.InitiateAuthInput{
//AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
//AuthParameters: map[string]*string{
//"REFRESH_TOKEN": aws.String(rt.Value),
//},
//ClientId: aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
//}
//res, err := cognito.New(sess).InitiateAuth(authTry)
//if err != nil {
//log.Printf("RefreshToken error - %s\n", err)
//http.Redirect(w, r, "/login", http.StatusSeeOther)
//return
//}
//c1 := http.Cookie{
//Name:     "AccessToken",
//Value:    *res.AuthenticationResult.AccessToken,
//HttpOnly: true,
//Path:     "/",
//Expires:  time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
//}
//c2 := http.Cookie{
//Name:     "IdToken",
//Value:    *res.AuthenticationResult.IdToken,
//HttpOnly: true,
//Path:     "/",
//Expires:  time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
//}
//http.SetCookie(w, &c1)
//http.SetCookie(w, &c2)
//log.Println("New AccessToken and IdToken issued")
//}