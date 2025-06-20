package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type usersKey string

const userCtx usersKey = "user"

// GetUser godoc
//
//	@summary		Fetches an user profile
//	@Description	Fetches an user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (api *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := api.getUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			api.notFoundError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	if err := api.jsonResponse(w, http.StatusOK, user); err != nil {
		api.internalServerError(w, r, err)
	}
}

// FollowUser godoc
//
//	@summary		Follows an user
//	@Description	Follows an user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{object}	string  "User followed"
//	@Failure		400	{object}	error "Bad userID"
//	@Failure		404	{object}	error "User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [put]
func (api *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	followedId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if followerUser.ID == followedId {
		api.conflictError(w, r, nil)
		return
	}

	if err := api.store.Followers.Follow(r.Context(), followerUser.ID, followedId); err != nil {

		switch err {
		case store.ErrConflict:
			api.conflictError(w, r, errors.New("user is already being followed"))
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	if err := api.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		api.internalServerError(w, r, err)
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (api *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	unfollowedId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if followerUser.ID == unfollowedId {
		api.conflictError(w, r, nil)
		return
	}

	if err := api.store.Followers.Unfollow(r.Context(), followerUser.ID, unfollowedId); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	if err := api.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		api.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromCtx(r *http.Request) *store.User {
	user := r.Context().Value(userCtx).(*store.User)
	return user
}
