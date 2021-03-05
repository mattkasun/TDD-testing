package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func login(router *gin.Engine) (cookie http.Cookie) {
	formData := url.Values{
		"user": {"user"},
		"pass": {"user"},
	}
	payload := formData.Encode()
	request, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	resp := response.Result()
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "todo" {
			return *cookie
		}
		panic("did not get cookie")
	}
	return
}

func TestAdd(t *testing.T) {
	t.Run("shows add form when plus button pressed", func(t *testing.T) {
		router := SetupRouter()

		request, _ := http.NewRequest(http.MethodGet, "/add", nil)
		cookie := login(router)
		request.AddCookie(&cookie)

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Contains(t, response.Body.String(), "<title>Add", "wrong body")
	})
	t.Run("shows main page after new task form completed", func(t *testing.T) {
		router := SetupRouter()

		formData := url.Values{
			"Name": {"New Task"},
		}

		payload := formData.Encode()
		request, _ := http.NewRequest(http.MethodPost, "/add", strings.NewReader(payload))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Type", strconv.Itoa(len(payload)))
		cookie := login(router)
		request.AddCookie(&cookie)

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Contains(t, response.Body.String(), "<title>Application Name", "wrong body")

		assert.Equal(t, http.StatusCreated, response.Code, "Wrong Status Code got %v, wanted %v", response.Code, http.StatusCreated)
	})
}

func TestMainPageWithCookie(t *testing.T) {
	t.Run("shows main page if user is logged in", func(t *testing.T) {
		router := SetupRouter()

		request, _ := http.NewRequest(http.MethodGet, "/main", nil)
		cookie := login(router)
		request.AddCookie(&cookie)

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Contains(t, response.Body.String(), "<title>Application", "wrong body")
		assert.Equal(t, response.Code, http.StatusOK, "wrong status code, got %v, wanted %v", response.Code, http.StatusOK)
	})
}
