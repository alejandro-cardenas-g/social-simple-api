package main

import (
	"net/http"

	"github.com/alejandro-cardenas-g/social/internal/store"
)

type CreateCommentToPostPayload struct {
	Content string `json:"content"`
}

func (api *application) createCommentToPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	payload := &CreateCommentToPostPayload{}

	if err := readJSON(w, r, payload); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		UserID:  1,
		PostID:  post.ID,
		Content: payload.Content,
	}

	if err := api.store.Comments.Create(r.Context(), comment); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	if err := api.jsonResponse(w, http.StatusOK, comment); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}
