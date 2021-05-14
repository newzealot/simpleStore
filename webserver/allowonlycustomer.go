package main

import "net/http"

// AllowOnlyCustomer makes sure that protected route contains a valid Access and ID token.
// Merchant types are denied access.
func AllowOnlyCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PerformChecks(w, r, "customer")
		next.ServeHTTP(w, r)
	})
}
