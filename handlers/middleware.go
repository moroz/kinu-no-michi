package handlers

import (
	"context"
	"net/http"

	"github.com/moroz/kinu-no-michi/config"
	"github.com/moroz/kinu-no-michi/db/queries"
	"github.com/moroz/kinu-no-michi/lib/cookies"
)

func decodeSessionCookie(cs cookies.SessionStore, r *http.Request) (*appSession, error) {
	cookie, err := r.Cookie(config.SESSION_COOKIE_NAME)
	if err != nil {
		return nil, err
	}

	binary, err := cs.Decode(cookie.Value)
	if err != nil {
		return nil, err
	}

	return decodeSession(binary)
}

func FetchSession(cs cookies.SessionStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session appSession
			if decoded, err := decodeSessionCookie(cs, r); err == nil {
				session = *decoded
			}

			ctx := context.WithValue(r.Context(), "session", &session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func findCartBySession(db queries.DBTX, r *http.Request) *queries.GetCartByIDRow {
	session, ok := r.Context().Value("session").(*appSession)
	if !ok {
		return nil
	}

	if session.CartID == nil {
		return nil
	}

	cart, _ := queries.New(db).GetCartByID(r.Context(), *session.CartID)
	return cart
}

// LoadCartFromSession is a middleware that fetches the cart pointed to by the CartID field of the appSession struct and stores it at the `"cart"` key in the context as a `*queries.GetCartByIDRow` value.
func LoadCartFromSession(db queries.DBTX) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cart := &queries.GetCartByIDRow{}

			if found := findCartBySession(db, r); found != nil {
				cart = found
			}

			ctx := context.WithValue(r.Context(), "cart", cart)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
