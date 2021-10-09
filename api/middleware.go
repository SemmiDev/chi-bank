package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			Error(w, http.StatusUnauthorized, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			Error(w, http.StatusUnauthorized, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			Error(w, http.StatusUnauthorized, err)
			return
		}

		accessToken := fields[1]
		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			Error(w, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w,r.WithContext(context.WithValue(r.Context(), authorizationPayloadKey, payload)))
	})
}