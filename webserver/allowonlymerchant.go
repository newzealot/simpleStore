package main

import (
	"net/http"
)

type UserInfo struct {
	Type  string
	ID    string
	Email string
}

// AllowOnlyMerchant makes sure that protected route contains a valid Access and ID token.
// Customer types are denied access.
func AllowOnlyMerchant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PerformChecks(w, r, "merchant")
		next.ServeHTTP(w, r)
	})
}
