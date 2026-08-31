package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gobackend/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

// createCommentHandler godoc
//	@Summary		Create a comment
//	@Description	Add a comment to a specific post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int						true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Router			/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Hardcoded user ID for now, similar to posts
	userID := 1

	comment := &store.Comment{
		PostID:  postID,
		UserID:  int64(userID),
		Content: payload.Content,
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
