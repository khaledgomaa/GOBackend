package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gobackend/social/internal/store"
)

// DTO to limit specific properties by the user.
// The `validate` tags act like Data Annotations in .NET (e.g., [Required]).
type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
	UserID  int64    `json:"user_id" validate:"required"`
}

// createPostHandler godoc
//	@Summary		Create a post
//	@Description	Create a new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate the struct against its tags.
	// Equivalent to checking ModelState.IsValid in a .NET Controller or using a FluentValidation validator.
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	// Validate user exists
	_, err := app.store.Users.GetByID(ctx, payload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, errors.New("user not found"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Created as reference type because app.store.Posts.Create will mutate it
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  payload.UserID,
		Tags:    payload.Tags,
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPostHandler godoc
//	@Summary		Get a post
//	@Description	Get a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	store.Post
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deletePostHandler godoc
//	@Summary		Delete a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		204		"No Content"
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, postID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Tags    []string `json:"tags"`
}

// Hanlding concurrent updates is out of scope for now
// updatePostHandler godoc
//	@Summary		Update a post
//	@Description	Partially update a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post update payload"
//	@Success		200		{object}	store.Post
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.Update(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getUserPostsHandler godoc
//	@Summary		Get user posts
//	@Description	Get a paginated list of a user's posts
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Success		200		{array}		store.Post
//	@Router			/users/{userID}/posts [get]
func (app *application) getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err = fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	posts, err := app.store.Posts.GetUserPosts(ctx, userID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Handle case where posts is nil to return empty array instead of null
	if posts == nil {
		posts = []store.Post{}
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}
}
