package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddArtist(t *testing.T) {
	name := "John Doe"
	expected := http.StatusOK
	rec := httptest.NewRecorder()
	albumHandler.AddAlbum(ctx)
	assert.Equal(t, expected, rec.Code)
}
