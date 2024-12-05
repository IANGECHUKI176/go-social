package main

import (
	"gopher_social/internal/store"
	"log"
	"net/http"
)

type CreateCommentPayload struct {
	Content string `json:"content"`
	UserID  int64  `json:"user_id"`
	PostID  int64  `json:"post_id"`
}

// func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {

// }
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)
	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	log.Println("payload", payload)
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	comment := &store.Comment{
		Content: payload.Content,
		UserID:  payload.UserID,
		PostID:  post.ID,
	}
	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comments); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
