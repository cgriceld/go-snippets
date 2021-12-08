// Non-auth

// POST /signup : register new user with login and password, login must be unique
//                request  : JSON {"login":"ivan", "password":"123"}
//                response : JSON of new created user (on success) {"login":"ivan", "password":"123"}

// POST /signin : signin user and return token
//                request  : JSON {"login":"ivan", "password":"123"}
//                response : JSON of authorized user (on success) {"login":"ivan", "password":"123", "token":"/token/"}

// ==================================================

// Auth (needs token)

// GET  /home                 : "start page"
//                              response : plain text message ("Hi, ivan!")

// GET  /chat                 : read all messages from the chat
//                              response : JSON [{"login":"ivan","text":"hello"}]

// POST /chat                 : send new message to the chat
//                              request  : JSON {"text":"hello"}
//                              response : new message with login of sender {"login":"ivan","text":"hello"}

// GET  /private/me           : read all your private messages
//                              response : JSON [{"login":"ivan","text":"hello"}]

// POST /private/send/{login} : send private message to user with {login}, BadRequest if user which such login doesn't exist
//                              request  : JSON {"text":"hello"}
//                              response : new message with login of reciever {"login":"ivan","text":"hello"}

// GET /chat and GET /private/me supports pagination (e.g. GET /chat?page=1). Size of page is 3.
// If page number is larger than len of chat or page = 0, you get whole chat in response.

package main

import (
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

const (
	port     = ":5000"
	pageSize = 3
)

type Message struct {
	Login string `json:"login"`
	Text  string `json:"text"`
}

type Handler struct {
	users     UserStorage
	chat      ChatStorage
	private   PrivateStorage
	tokenAuth *jwtauth.JWTAuth
}

func newHandler(dbUsers UserStorage, dbChat ChatStorage, dbPrivate PrivateStorage) *Handler {
	h := &Handler{
		users:   dbUsers,
		chat:    dbChat,
		private: dbPrivate,
	}
	h.tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	return h
}

func main() {
	dbUsers := newUserStorage()
	dbPrivate := newPrivateStorage()
	dbChat := newChatStorage()

	handler := newHandler(dbUsers, dbChat, dbPrivate)

	log.Fatal(http.ListenAndServe(port, handler.Routes()))
}
