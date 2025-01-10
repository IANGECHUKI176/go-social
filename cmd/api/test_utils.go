package main

import (
	"gopher_social/internal/auth"
	"gopher_social/internal/store"
	"gopher_social/internal/store/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewTestApplication(t *testing.T, cfg config) *application {
	t.Helper()
	logger := zap.Must(zap.NewProduction()).Sugar()
	// logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	cacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  cacheStore,
		authenticator: testAuth,
		config:        cfg,
	}
}
func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d, but got %d", expected, actual)
	}
}
