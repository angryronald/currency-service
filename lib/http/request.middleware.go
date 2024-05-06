package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/lib/appcontext"
)

func AuthenticateClient(allowedClients map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// // You can access the parameter passed to the middleware here
			// fmt.Println("Parameter passed to middleware:", param)
			authTokenInString := r.Header.Get("authorization")
			authTokens := strings.Split(authTokenInString, " ")
			if strings.ToLower(authTokens[0]) != strings.ToLower("bearer") {
				log.Printf("forbidden access (not allowed authentication method for client) %v: (from %v) on [%s] %s", time.Now().UTC(), r.RemoteAddr, r.Method, r.URL.String())
				ResponseError(w, http.StatusForbidden, "", logrus.New())
				return
			}
			if allowedClients[authTokens[1]] == "" {
				log.Printf("forbidden access (not registered client) %v: (from %v) on [%s] %s", time.Now().UTC(), r.RemoteAddr, r.Method, r.URL.String())
				ResponseError(w, http.StatusForbidden, "", logrus.New())
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), appcontext.KeyClientID, allowedClients[authTokens[1]]))

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
