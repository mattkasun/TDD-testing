package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMainPage(t *testing.T) {

	t.Run("shows Landing Page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		Router := SetupRouter()
		gin.SetMode(gin.TestMode)

		Router.ServeHTTP(response, request)
		want := "<Title>Landing"
		searchString := regexp.MustCompile(want)

		found := searchString.MatchString(response.Body.String())

		if !found {
			t.Errorf("did not find %q", want)
		}
	})
}
