package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func getUserFromBody(body io.Reader) (*User, int) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	var user User
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, http.StatusBadRequest
	}

	return &user, 0
}

func getMessageFromBody(body io.Reader) (*Message, int) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	var mess Message
	if err = json.Unmarshal(b, &mess); err != nil {
		return nil, http.StatusBadRequest
	}

	return &mess, 0
}

func getMessageByPage(m []Message, page int) []Message {
	chatSize := len(m)
	if page > 0 && chatSize > 0 {
		if start := (page - 1) * pageSize; start < chatSize {
			end := start + pageSize
			if end > chatSize {
				end = chatSize
			}
			return m[start:end]
		}
	}
	return m
}
