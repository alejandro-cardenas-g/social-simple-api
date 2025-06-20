package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-cardenas-g/social/internal/auth"
	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/alejandro-cardenas-g/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	cacheStore := cache.NewMockStorage()
	testAuth := auth.NewTestAuthenticator()

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  cacheStore,
		authenticator: testAuth,
		config:        cfg,
	}

}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected the response code to be %d and we got %d", expected, actual)
	}
}
