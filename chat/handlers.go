package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/home", h.Home)

		r.With(Pagination).Get("/chat", h.readChat)
		r.Post("/chat", h.postChat)

		r.With(Pagination).Get("/private/me", h.readPrivate)
		r.Post("/private/send/{login}", h.postPrivate)
	})

	r.Group(func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		r.Post("/signin", h.SignIn)
	})

	return r
}

func Pagination(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pageQ := r.URL.Query().Get("page")
		page := 0
		if val, err := strconv.Atoi(pageQ); err == nil {
			page = val
		}

		ctx := context.WithValue(r.Context(), "page_id", page)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromBody(r.Body)
	if err != 0 {
		w.WriteHeader(err)
		return
	}
	defer r.Body.Close()

	created, cerr := h.users.Create(user)
	if cerr != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, cerr.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, created)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromBody(r.Body)
	if err != 0 {
		w.WriteHeader(err)
		return
	}
	defer r.Body.Close()

	curr, serr := h.users.checkUser(user, h.tokenAuth)
	if serr != nil {
		render.Status(r, http.StatusUnauthorized)
		render.PlainText(w, r, serr.Error())
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, curr)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	render.PlainText(w, r, fmt.Sprintf("Hi, %v!\n", claims["login"]))
}

func (h *Handler) postChat(w http.ResponseWriter, r *http.Request) {
	mess, err := getMessageFromBody(r.Body)
	if err != 0 {
		w.WriteHeader(err)
		return
	}
	defer r.Body.Close()

	if mess.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	sessionLogin, ok := claims["login"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	send := h.chat.SendToChat(&Message{Login: sessionLogin, Text: mess.Text})

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, *send)
}

func (h *Handler) readChat(w http.ResponseWriter, r *http.Request) {
	page, ok := r.Context().Value("page_id").(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, h.chat.GetChat(page))
}

func (h *Handler) postPrivate(w http.ResponseWriter, r *http.Request) {
	mess, err := getMessageFromBody(r.Body)
	if err != 0 {
		w.WriteHeader(err)
		return
	}
	defer r.Body.Close()

	if mess.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reciever := chi.URLParam(r, "login")
	if _, ok := h.users.getByLogin(reciever); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	sessionLogin, ok := claims["login"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	send := h.private.SendToPrivate(reciever, &Message{Login: sessionLogin, Text: mess.Text})

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, *send)
}

func (h *Handler) readPrivate(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sessionLogin, ok := claims["login"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	page, ok := r.Context().Value("page_id").(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, h.private.GetPrivate(sessionLogin, page))
}
