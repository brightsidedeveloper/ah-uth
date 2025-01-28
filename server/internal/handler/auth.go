package handler

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/buf"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) GetAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(err)
		h.bin.WriteError(w, http.StatusUnauthorized, "error completing user auth")
		return
	}

	session, err := gothic.Store.Get(r, "session")
	fmt.Println("Session:", session)
	fmt.Println("User:", user)
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		h.bin.WriteError(w, http.StatusInternalServerError, "error retrieving session")
		return
	}

	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Error saving session:", err)
		h.bin.WriteError(w, http.StatusInternalServerError, "error saving session")
		return
	}

	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)

}

func (h *Handler) GetAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		session, err := gothic.Store.Get(r, "session")
		if err != nil {
			fmt.Println("Error retrieving session:", err)
			h.bin.WriteError(w, http.StatusInternalServerError, "error retrieving session")
			return
		}

		session.Values["user"] = gothUser
		err = session.Save(r, w)
		if err != nil {
			fmt.Println("Error saving session:", err)
			h.bin.WriteError(w, http.StatusInternalServerError, "error saving session")
			return
		}
		h.bin.ProtoWrite(w, http.StatusOK, &buf.User{
			Name: gothUser.Name,
		})
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) GetLogout(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	session, _ := gothic.Store.Get(r, "session")
	delete(session.Values, "user")
	session.Save(r, w)

	gothic.Logout(w, r)
	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserStore(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "session")
	if err != nil {
		fmt.Println(err)
		h.bin.WriteError(w, http.StatusInternalServerError, "error getting session")
		return
	}
	user, ok := session.Values["user"]
	if !ok {
		h.bin.WriteError(w, http.StatusUnauthorized, "error getting user session")
		return
	}
	gothUser, ok := user.(goth.User)
	if !ok {
		h.bin.WriteError(w, http.StatusInternalServerError, "error getting user")
		return
	}
	fmt.Println("Goth User:", gothUser)
	h.bin.ProtoWrite(w, http.StatusOK, &buf.User{
		Name: gothUser.Email,
	})
}
