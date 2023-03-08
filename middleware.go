package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PathParam is a middleware that parses a path parameter and stores it in the request context.
func PathParam[T any](key string, parser func(r *http.Request, key string) (T, error)) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value, err := parser(r, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), pathParamKey, value))
			h.ServeHTTP(w, r)
		})
	}
}

// PathParamFromRequest retrieves a path parameter from the request context.
func PathParamFromRequest[T any](r *http.Request) (T, error) {
	return PathParamFromContext[T](r.Context())
}

// PathParamFromContext retrieves a path parameter from the request context.
func PathParamFromContext[T any](ctx context.Context) (T, error) {
	v, ok := ctx.Value(pathParamKey).(T)
	if !ok {
		return v, fmt.Errorf("path param not found")
	}

	return v, nil
}

func UUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uid == uuid.Nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), uuidKey, uid))
		next.ServeHTTP(w, r)
	})
}

func UUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	return UUIDFromContext(r.Context())
}

func UUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	uid, ok := ctx.Value(uuidKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("uuid not found in context")
	}
	return uid, nil
}

type contextKey string

const (
	uuidKey      = contextKey("uuid")
	pathParamKey = contextKey("path-param")
)

func (k contextKey) String() string {
	return fmt.Sprintf("registry context key %q", string(k))
}
