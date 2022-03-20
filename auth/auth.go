package auth

import (
	"net/http"
)

const UserTokenHeader = "X-User-Token"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		userToken := request.Header.Get(UserTokenHeader)
		if userToken == "" {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}
