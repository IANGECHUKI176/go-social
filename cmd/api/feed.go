package main

import (
	"gopher_social/internal/store"
	"net/http"
)

// @Summary		Fetch user feed
// @Description	Fetch user feed
// @Tags			feed
// @Accept			json
// @Produce		json
//
// @Param			limit	query		int			false	"Limit"
// @Param			offset	query		int			false	"Offset"
// @Param			sort	query		string		false	"Sort"
//
// @Param			tags	query		[]string	false	"Tags (repeatable for multiple values)"
// @Param			search	query		string		false	"Search"
// @Param			since	query		string		false	"Since"
// @Param			until	query		string		false	"Until"
//
// @Success		200		{object}	[]store.PostWithMetadata
// @Failure		500		{object}	error
// @Failure		400		{object}	error
// @Security		ApiKeyAuth
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// pagination ,filters,sort
	//feed?limit=10&offset=0
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Search: "",
		Tags:   []string{},
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(13), fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
